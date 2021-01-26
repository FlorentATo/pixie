from pxapi import cloudapi_pb2_grpc, cpb, vizier_pb2_grpc, vpb, test_utils as utils
import pxapi
import unittest
import grpc
import asyncio
import uuid
from concurrent import futures
from typing import List, Any, Coroutine, Dict


ACCESS_TOKEN = "12345678-0000-0000-0000-987654321012"
pxl_script = """
import px
px.display(px.DataFrame('http_events')[
           ['http_resp_body','http_resp_status']].head(10), 'http')
px.display(px.DataFrame('process_stats')[
           ['upid','cpu_ktime_ns', 'rss_bytes']].head(10), 'stats')
"""


async def run_script_and_tasks(
    script: pxapi.Script,
    processors: List[Coroutine[Any, Any, Any]]
) -> None:
    """
    Runs data processors in parallel with a script, returns the result of the coroutine.
    """
    tasks = [asyncio.create_task(p) for p in processors]
    try:
        await script.run_async()
    except Exception as e:
        for t in tasks:
            t.cancel()
        raise e

    await asyncio.gather(*tasks)


class VizierServiceFake:
    def __init__(self) -> None:
        self.cluster_id_to_fake_data: Dict[str,
                                           List[vpb.ExecuteScriptResponse]] = {}
        self.cluster_id_to_error: Dict[str, Exception] = {}

    def add_fake_data(self, cluster_id: str, data: List[vpb.ExecuteScriptResponse]) -> None:
        if cluster_id not in self.cluster_id_to_fake_data:
            self.cluster_id_to_fake_data[cluster_id] = []
        self.cluster_id_to_fake_data[cluster_id].extend(data)

    def trigger_error(self, cluster_id: str, exc: Exception) -> None:
        """ Adds an error that triggers after the data is yielded. """
        self.cluster_id_to_error[cluster_id] = exc

    def ExecuteScript(self, request: vpb.ExecuteScriptRequest, context: Any) -> Any:
        cluster_id = request.cluster_id
        assert cluster_id in self.cluster_id_to_fake_data, f"need data for cluster_id {cluster_id}"
        data = self.cluster_id_to_fake_data[cluster_id]
        for d in data:
            yield d
        # Trigger an error for the cluster ID if the user added one.
        if cluster_id in self.cluster_id_to_error:
            raise self.cluster_id_to_error[cluster_id]

    def HealthCheck(self, request: Any, context: Any) -> Any:
        yield vpb.Status(code=1, message="fail")


def create_cluster_info(
    cluster_id: str,
    cluster_name: str,
    status: cpb.ClusterStatus = cpb.CS_HEALTHY,
    passthrough_enabled: bool = True
) -> cpb.ClusterInfo:
    return cpb.ClusterInfo(
        id=utils.create_uuid_pb(cluster_id),
        status=status,
        config=cpb.VizierConfig(
            passthrough_enabled=passthrough_enabled,
        ),
        pretty_cluster_name=cluster_name,
    )


class CloudServiceFake(cloudapi_pb2_grpc.VizierClusterInfoServicer):
    def __init__(self) -> None:
        self.clusters = [
            create_cluster_info(
                utils.cluster_uuid1,
                "cluster1",
            ),
            create_cluster_info(
                utils.cluster_uuid2,
                "cluster2",
            ),
            # One cluster marked as unhealthy.
            create_cluster_info(
                utils.cluster_uuid3,
                "cluster3",
                status=cpb.CS_UNHEALTHY,
            ),
        ]
        self.direct_conn_info: Dict[str,
                                    cpb.GetClusterConnectionInfoResponse] = {}

    def add_direct_conn_cluster(
        self,
        cluster_id: str,
        name: str,
        url: str,
        token: str,
    ) -> None:
        for c in self.clusters:
            if cluster_id == c.id:
                raise ValueError(
                    f"Cluster id {cluster_id} already exists in cloud.")

        self.clusters.append(create_cluster_info(
            cluster_id,
            name,
            passthrough_enabled=False,
        ))

        self.direct_conn_info[cluster_id] = cpb.GetClusterConnectionInfoResponse(
            ipAddress=url,
            token=token,
        )

    def GetClusterInfo(
        self,
        request: cpb.GetClusterInfoRequest,
        context: Any,
    ) -> cpb.GetClusterInfoResponse:
        if request.id.data:
            for c in self.clusters:
                if c.id == request.id:
                    return cpb.GetClusterInfoResponse(clusters=[c])
            return cpb.GetClusterInfoResponse(clusters=[])

        return cpb.GetClusterInfoResponse(clusters=self.clusters)

    def GetClusterConnectionInfo(
        self,
        request: cpb.GetClusterConnectionInfoRequest,
        context: Any,
    ) -> cpb.GetClusterConnectionInfoResponse:
        cluster_id = request.id.data.decode('utf-8')
        if cluster_id not in self.direct_conn_info:
            raise KeyError(
                f"Cluster ID {cluster_id} not in conn_info. Must add first.")
        return self.direct_conn_info[cluster_id]


class TestClient(unittest.TestCase):
    def setUp(self) -> None:
        # Create a fake server for the VizierService
        self.server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
        self.fake_vizier_service = VizierServiceFake()
        self.fake_cloud_service = CloudServiceFake()

        vizier_pb2_grpc.add_VizierServiceServicer_to_server(
            self.fake_vizier_service, self.server)

        cloudapi_pb2_grpc.add_VizierClusterInfoServicer_to_server(
            self.fake_cloud_service, self.server)
        self.port = self.server.add_insecure_port("[::]:0")
        self.server.start()

        self.http_table_factory = utils.FakeTableFactory("http", vpb.Relation(columns=[
            utils.string_col("http_resp_body"),
            utils.int64_col("http_resp_status"),
        ]))

        self.stats_table_factory = utils.FakeTableFactory("stats", vpb.Relation(columns=[
            utils.uint128_col("upid"),
            utils.int64_col("cpu_ktime_ns"),
            utils.int64_col("rss_bytes"),
        ]))

    def url(self) -> str:
        return f"localhost:{self.port}"

    def tearDown(self) -> None:
        self.server.stop(None)

    def test_list_healthy_clusters(self) -> None:
        # Tests that users can list healthy clusters and then
        # script those clusters.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )

        clusters = px_client.list_healthy_clusters()
        self.assertSetEqual(
            set([c.name() for c in clusters]),
            {"cluster1", "cluster2"}
        )

        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(clusters[0])

        # Create one http table.
        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            # Init "http".
            http_table1.metadata_response(),
            # Send data for "http".
            http_table1.row_batch_response([["foo"], [200]]),
            # End "http".
            http_table1.end(),
        ])

        script = px_client.connect_to_cluster(
            clusters[0]).create_script(pxl_script)

        script.add_callback("http", lambda row: None)

        # Run the script synchronously.
        script.run()

    def test_one_conn_one_table(self) -> None:
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            # Channel functions for testing.
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        # Create table for cluster_uuid1.
        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            # Initialize the table on the stream with the metadata.
            http_table1.metadata_response(),
            # Send over a single-row batch.
            http_table1.row_batch_response([["foo"], [200]]),
            # Send an end-of-stream for the table.
            http_table1.end(),
        ])

        # Create the script object.
        script = conn.create_script(pxl_script)
        # Subscribe to the http table.
        http_tb = script.subscribe("http")

        # Define an async function that processes the TableSub
        # while the API script can run concurrently.
        async def process_table(table_sub: pxapi.TableSub) -> None:
            num_rows = 0
            async for row in table_sub:
                self.assertEqual(row["http_resp_body"], "foo")
                self.assertEqual(row["http_resp_status"], 200)
                num_rows += 1

            self.assertEqual(num_rows, 1)

        # Run the script and process_table concurrently.
        loop = asyncio.get_event_loop()
        loop.run_until_complete(
            run_script_and_tasks(script, [process_table(http_tb)]))

    def test_multiple_rows_and_rowbatches(self) -> None:
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        # Create table for the first cluster.
        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        rb_data: List[List[Any]] = [
            ["foo", "bar", "baz", "bat"], [200, 500, 301, 404]]

        # Here we split the above data into two rowbatches.
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            # Initialize the table on the stream with the metadata.
            http_table1.metadata_response(),
            # Row batch 1 has data 1-3,
            http_table1.row_batch_response([rb_data[0][:3], rb_data[1][:3]]),
            # Row batch 2 has data 4
            http_table1.row_batch_response([rb_data[0][3:], rb_data[1][3:]]),
            # Send an end-of-stream for the table.
            http_table1.end(),
        ])

        # Create the script object.
        script = conn.create_script(pxl_script)
        # Subscribe to the http table.
        http_tb = script.subscribe("http")

        # Verify that the rows returned by the table_sub match the
        # order and values of the input test data.
        async def process_table(table_sub: pxapi.TableSub) -> None:
            row_i = 0
            # table_sub hides the batched rows and delivers them in
            # the same order as the batches sent.
            async for row in table_sub:
                self.assertEqual(row["http_resp_body"], rb_data[0][row_i])
                self.assertEqual(row["http_resp_status"], rb_data[1][row_i])
                row_i += 1

            self.assertEqual(row_i, 4)

        # Run the script and process_table concurrently.
        loop = asyncio.get_event_loop()
        loop.run_until_complete(
            run_script_and_tasks(script, [process_table(http_tb)]))

    def test_one_conn_two_tables(self) -> None:
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        # We will send two tables for this test "http" and "stats".
        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        stats_table1 = self.stats_table_factory.create_table(utils.table_id3)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            # Initialize "http" on the stream.
            http_table1.metadata_response(),
            # Send over a row-batch from "http".
            http_table1.row_batch_response([["foo"], [200]]),
            # Initialize "stats" on the stream.
            stats_table1.metadata_response(),
            # Send over a row-batch from "stats".
            stats_table1.row_batch_response([
                [vpb.UInt128(high=123, low=456)],
                [1000],
                [999],
            ]),
            # Send an end-of-stream for "http".
            http_table1.end(),
            # Send an end-of-stream for "stats".
            stats_table1.end(),
        ])

        script = conn.create_script(pxl_script)
        # Subscribe to both tables.
        http_tb = script.subscribe("http")
        stats_tb = script.subscribe("stats")

        # Async function that makes sure "http" table returns the expected rows.
        async def process_http_tb(table_sub: pxapi.TableSub) -> None:
            num_rows = 0
            async for row in table_sub:
                self.assertEqual(row["http_resp_body"], "foo")
                self.assertEqual(row["http_resp_status"], 200)
                num_rows += 1

            self.assertEqual(num_rows, 1)

        # Async function that makes sure "stats" table returns the expected rows.
        async def process_stats_tb(table_sub: pxapi.TableSub) -> None:
            num_rows = 0
            async for row in table_sub:
                self.assertEqual(row["upid"], uuid.UUID(
                    '00000000-0000-007b-0000-0000000001c8'))
                self.assertEqual(row["cpu_ktime_ns"], 1000)
                self.assertEqual(row["rss_bytes"], 999)
                num_rows += 1

            self.assertEqual(num_rows, 1)
        # Run the script and the processing tasks concurrently.
        loop = asyncio.get_event_loop()
        loop.run_until_complete(
            run_script_and_tasks(script, [process_http_tb(http_tb), process_stats_tb(stats_tb)]))

    def test_run_script_with_invalid_arg_error(self) -> None:
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        # Send over an error in the Status field. This is the exact error you would
        # get if you sent over an empty pxl function in the ExecuteScriptRequest.
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            vpb.ExecuteScriptResponse(status=utils.invalid_argument(
                message="Script should not be empty."
            ))
        ])

        # Prepare the script and run synchronously.
        script = conn.create_script("")
        # Although we add a callback, we don't want this to throw an error.
        # Instead the error should be returned by the run function.
        script.add_callback("http_table", lambda row: print(row))
        with self.assertRaisesRegex(ValueError, "Script should not be empty."):
            script.run()

    def test_run_script_with_line_col_error(self) -> None:
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        # Send over an error a line, column error. These kinds of errors come
        # from the compiler pointing to a specific failure in the pxl script.
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            vpb.ExecuteScriptResponse(status=utils.line_col_error(
                1,
                2,
                message="name 'aa' is not defined"
            ))
        ])

        # Prepare the script and run synchronously.
        script = conn.create_script("aa")
        # Although we add a callback, we don't want this to throw an error.
        # Instead the error should be returned by the run function.
        script.add_callback("http_table", lambda row: print(row))
        with self.assertRaisesRegex(pxapi.PxLError, "PxL, line 1.*name 'aa' is not defined"):
            script.run()

    def test_run_script_with_api_errors(self) -> None:
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        # Only send data for "http".
        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            http_table1.metadata_response(),
            http_table1.row_batch_response([["foo"], [200]]),
            http_table1.end(),
        ])

        script = conn.create_script(pxl_script)

        # Subscribe to a table that doesn't exist shoudl throw an error.
        foobar_tb = script.subscribe("foobar")

        # Try to pull data from the foobar_tb, but error out when the script
        # never produces that data.
        loop = asyncio.get_event_loop()
        with self.assertRaisesRegex(ValueError, "Table 'foobar' not received"):
            loop.run_until_complete(
                run_script_and_tasks(script, [utils.iterate_and_pass(foobar_tb)]))

    def test_run_script_callback(self) -> None:
        # Test the callback API. Callback API is a simpler alternative to the TableSub
        # API that allows you to designate a function that runs on individual rows. Users
        # can process data without worrying about async processing by using this API.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )

        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        # Create two tables: "http" and "stats"
        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        stats_table1 = self.stats_table_factory.create_table(utils.table_id3)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            # Init "http".
            http_table1.metadata_response(),
            # Send data for "http".
            http_table1.row_batch_response([["foo"], [200]]),
            # Init "stats".
            stats_table1.metadata_response(),
            # Send data for "stats".
            stats_table1.row_batch_response([
                [vpb.UInt128(high=123, low=456)],
                [1000],
                [999],
            ]),
            # End "http".
            http_table1.end(),
            # End "stats".
            stats_table1.end(),
        ])

        script = conn.create_script(pxl_script)
        http_counter = 0
        stats_counter = 0

        # Define callback function for "http" table.
        def http_fn(row: pxapi.Row) -> None:
            nonlocal http_counter
            http_counter += 1
            self.assertEqual(row["http_resp_body"], "foo")
            self.assertEqual(row["http_resp_status"], 200)
        script.add_callback("http", http_fn)

        # Define a callback function for the stats_fn.
        def stats_fn(row: pxapi.Row) -> None:
            nonlocal stats_counter
            stats_counter += 1
            self.assertEqual(row["upid"], uuid.UUID(
                '00000000-0000-007b-0000-0000000001c8'))
            self.assertEqual(row["cpu_ktime_ns"], 1000)
            self.assertEqual(row["rss_bytes"], 999)
        script.add_callback("stats", stats_fn)

        # Run the script synchronously.
        script.run()

        # We expect each callback function to only be called once.
        self.assertEqual(stats_counter, 1)
        self.assertEqual(http_counter, 1)

    def test_run_script_callback_with_error(self) -> None:
        # Test to demonstrate how errors raised in callbacks can be handled.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )

        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        # Create HTTP table and add to the stream.
        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            http_table1.metadata_response(),
            http_table1.row_batch_response([["foo"], [200]]),
            http_table1.end(),
        ])

        script = conn.create_script(pxl_script)

        # Add callback function that raises an error.
        def http_fn(row: pxapi.Row) -> None:
            raise ValueError("random internal error")
        script.add_callback("http", http_fn)

        # Run script synchronously, expecting the internal error to propagate up.
        with self.assertRaisesRegex(ValueError, "random internal error"):
            script.run()

    def test_subscribe_all(self) -> None:
        # Tests `subscribe_all_tables()`.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        # Create two tables and simulate them sent over as part of the ExecuteScript call.
        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        stats_table1 = self.stats_table_factory.create_table(utils.table_id3)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            http_table1.metadata_response(),
            http_table1.row_batch_response([["foo"], [200]]),
            stats_table1.metadata_response(),
            stats_table1.row_batch_response([
                [vpb.UInt128(high=123, low=456)],
                [1000],
                [999],
            ]),
            http_table1.end(),
            stats_table1.end(),
        ])
        # Create script.
        script = conn.create_script(pxl_script)
        # Get a subscription to all of the tables that arrive over the stream.
        tables = script.subscribe_all_tables()

        # Async function to run on the "http" table.
        async def process_http_tb(table_sub: pxapi.TableSub) -> None:
            num_rows = 0
            async for row in table_sub:
                self.assertEqual(row["http_resp_body"], "foo")
                self.assertEqual(row["http_resp_status"], 200)
                num_rows += 1

            self.assertEqual(num_rows, 1)

        # Async function that processes the tables subscription and runs the
        # async function above to process the "http" table when that table shows
        # we see thtparticular table.
        async def process_all_tables(tables_gen: pxapi.TableSubGenerator) -> None:
            table_names = set()
            async for table in tables_gen:
                table_names.add(table.table_name)
                if table.table_name == "http":
                    # Once we find the http_tb, process it.
                    await process_http_tb(table)
            # Make sure we see both tables on the generator.
            self.assertEqual(table_names, {"http", "stats"})

        # Run the script and process_all_tables function concurrently.
        # We expect no errors.
        loop = asyncio.get_event_loop()
        loop.run_until_complete(
            run_script_and_tasks(script, [process_all_tables(tables())]))

    def test_subscribe_same_table_twice(self) -> None:
        # Only on subscription allowed per table. Users should handle data from
        # the single alloatted subscription to enable the logical equivalent
        # of multiple subscriptions to one table.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        script = conn.create_script(pxl_script)

        # First subscription is fine.
        script.subscribe("http")

        # Second raises an error.
        with self.assertRaisesRegex(ValueError, "Already subscribed to 'http'"):
            script.subscribe("http")

    def test_fail_on_multi_run(self) -> None:
        # Tests to show that queries may only be run once. After a script has been
        # run, calling data grabbing methods like subscribe, add_callback, etc. will
        # raise an error.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        stats_table1 = self.stats_table_factory.create_table(utils.table_id3)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            stats_table1.metadata_response(),
            stats_table1.row_batch_response([
                [vpb.UInt128(high=123, low=456)],
                [1000],
                [999],
            ]),
            stats_table1.end(),
        ])

        script = conn.create_script(pxl_script)

        # Create a dummy callback.
        def stats_cb(row: pxapi.Row) -> None:
            pass

        script.add_callback("stats", stats_cb)
        # Run the script for the first time. Should not return an error.
        script.run()
        # Each of the following methods should fail if called after script.run()
        script_ran_message = "Script already ran"
        # Adding a callback should fail.
        with self.assertRaisesRegex(ValueError, script_ran_message):
            script.add_callback("stats", stats_cb)
        # Subscribing to a table should fail.
        with self.assertRaisesRegex(ValueError, script_ran_message):
            script.subscribe("stats")
        # Subscribing to all tables should fail.
        with self.assertRaisesRegex(ValueError, script_ran_message):
            script.subscribe_all_tables()
        # Synchronous run should error out.
        with self.assertRaisesRegex(ValueError, script_ran_message):
            script.run()
        # Async run should error out.
        loop = asyncio.get_event_loop()
        with self.assertRaisesRegex(ValueError, script_ran_message):
            loop.run_until_complete(script.run_async())

    def test_send_error_id_table_prop(self) -> None:
        # Sending an error over the stream should cause the table sub to exit.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            # Initialize the table on the stream and send over a rowbatch.
            http_table1.metadata_response(),
            http_table1.row_batch_response([["foo"], [200]]),
            # Send over an error on the stream after we've started sending data.
            # this should happen if something breaks on the Pixie side.
            # Note: the table does not send an end message over the stream.
            vpb.ExecuteScriptResponse(status=utils.invalid_argument(
                message="server error"
            ))
        ])

        # Create the script object.
        script = conn.create_script(pxl_script)
        # Add callback for http table.
        script.add_callback("http", lambda _: None)

        with self.assertRaisesRegex(ValueError, "server error"):
            script.run()

    def test_stop_sending_data_before_eos(self) -> None:
        # If the stream stops before sending over an eos for each table that should be an error.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            # Initialize the table on the stream and send over a rowbatch.
            http_table1.metadata_response(),
            http_table1.row_batch_response([["foo"], [200]]),
            # Note: the table does not send an end message over the stream.
        ])

        # Create the script object.
        script = conn.create_script(pxl_script)
        # Subscribe to the http table.
        http_tb = script.subscribe("http")

        # Run the script and process_table concurrently.
        loop = asyncio.get_event_loop()
        with self.assertRaisesRegex(ValueError, "Closed before receiving end-of-stream."):
            loop.run_until_complete(
                run_script_and_tasks(script, [utils.iterate_and_pass(http_tb)]))

    def test_handle_server_side_errors(self) -> None:
        # Test to make sure server side errors are handled somewhat.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a single fake cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            # Initialize the table on the stream with the metadata.
            http_table1.metadata_response(),
            # Send over a single-row batch.
            http_table1.row_batch_response([["foo"], [200]]),
            # NOTE: don't send over the eos -> simulating error midway through
            # stream.

        ])
        self.fake_vizier_service.trigger_error(
            utils.cluster_uuid1, ValueError('hi'))
        # Create the script object.
        script = conn.create_script(pxl_script)
        # Subscribe to the http table.
        http_tb = script.subscribe("http")

        # Run the script and process_table concurrently.
        loop = asyncio.get_event_loop()
        with self.assertRaisesRegex(grpc.aio.AioRpcError, "hi"):
            loop.run_until_complete(
                run_script_and_tasks(script, [utils.iterate_and_pass(http_tb)]))

    def test_direct_conns(self) -> None:
        # Test the direct connections.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Create the direct conn server.
        dc_server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
        fake_dc_service = VizierServiceFake()
        # Create one http table.
        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        cluster_id = "10000000-0000-0000-0000-000000000004"
        fake_dc_service.add_fake_data(cluster_id, [
            # Init "http".
            http_table1.metadata_response(),
            # Send data for "http".
            http_table1.row_batch_response([["foo"], [200]]),
            # End "http".
            http_table1.end(),
        ])

        vizier_pb2_grpc.add_VizierServiceServicer_to_server(
            fake_dc_service, dc_server)
        port = dc_server.add_insecure_port("[::]:0")
        dc_server.start()

        url = f"http://[::]:{port}"
        token = cluster_id
        self.fake_cloud_service.add_direct_conn_cluster(
            cluster_id, "dc_cluster", url, token)

        clusters = px_client.list_healthy_clusters()
        self.assertSetEqual(
            set([c.name() for c in clusters]),
            {"cluster1", "cluster2", "dc_cluster"}
        )

        conns = [
            px_client.connect_to_cluster(c) for c in clusters if c.name() == 'dc_cluster'
        ]

        self.assertEqual(len(conns), 1)

        script = conns[0].create_script(pxl_script)

        # Define callback function for "http" table.
        def http_fn(row: pxapi.Row) -> None:
            self.assertEqual(row["http_resp_body"], "foo")
            self.assertEqual(row["http_resp_status"], 200)

        script.add_callback("http", http_fn)

        # Run the script synchronously.
        script.run()

    def test_ergo_api(self) -> None:
        # Create a new API where we can run and get results for a table simultaneously.
        px_client = pxapi.Client(
            token=ACCESS_TOKEN,
            server_url=self.url(),
            # Channel functions for testing.
            channel_fn=lambda url: grpc.insecure_channel(url),
            conn_channel_fn=lambda url: grpc.aio.insecure_channel(url),
        )
        # Connect to a cluster.
        conn = px_client.connect_to_cluster(
            px_client.list_healthy_clusters()[0])

        # Create the script.
        script = conn.create_script(pxl_script)

        # Create table for cluster_uuid1.
        http_table1 = self.http_table_factory.create_table(utils.table_id1)
        self.fake_vizier_service.add_fake_data(conn.cluster_id, [
            # Initialize the table on the stream with the metadata.
            http_table1.metadata_response(),
            # Send over a single-row batch.
            http_table1.row_batch_response([["foo"], [200]]),
            # Send an end-of-stream for the table.
            http_table1.end(),
        ])

        # Use the results API to run and get the data from the http table.
        for row in script.results("http"):
            self.assertEqual(row["http_resp_body"], "foo")
            self.assertEqual(row["http_resp_status"], 200)


if __name__ == "__main__":
    unittest.main()
