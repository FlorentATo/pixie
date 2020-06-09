#include <google/protobuf/text_format.h>
#include <gtest/gtest.h>
#include <type_traits>

#include <absl/strings/match.h>
#include "src/carnot/planner/compiler_state/registry_info.h"
#include "src/common/base/base.h"
#include "src/common/testing/testing.h"
#include "src/shared/types/types.h"

namespace pl {
namespace carnot {
namespace planner {
using ::testing::UnorderedElementsAre;

constexpr char kExpectedUDFInfo[] = R"(
udas {
  name: "uda1"
  update_arg_types: INT64
  finalize_type: INT64
}
scalar_udfs {
  name: "add"
  exec_arg_types: FLOAT64
  exec_arg_types: FLOAT64
  return_type: FLOAT64
}
scalar_udfs {
  name: "scalar1"
  exec_arg_types: BOOLEAN
  exec_arg_types: INT64
  return_type: INT64
}
udtfs {
  name: "OpenNetworkConnections"
  args {
    name: "upid"
    arg_type: UINT128
    semantic_type: ST_UPID
  }
  executor: UDTF_SUBSET_PEM
  relation {
    columns {
      column_name: "time_"
      column_type: TIME64NS
    }
    columns {
      column_name: "fd"
      column_type: INT64
    }
    columns {
      column_name: "name"
      column_type: STRING
    }
  }
}
semantic_type_rules {
  name: "add"
  udf_exec_type: SCALAR_UDF
  exec_arg_types: ST_BYTES
  exec_arg_types: ST_BYTES
  output_type: ST_BYTES
}
semantic_type_rules {
  name: "uda1"
  udf_exec_type: UDA
  update_arg_types: ST_BYTES
  output_type: ST_BYTES
}
)";

TEST(RegistryInfo, basic) {
  auto info = RegistryInfo();
  udfspb::UDFInfo info_pb;
  google::protobuf::TextFormat::MergeFromString(kExpectedUDFInfo, &info_pb);
  EXPECT_OK(info.Init(info_pb));

  EXPECT_EQ(UDFExecType::kUDA, info.GetUDFExecType("uda1").ConsumeValueOrDie());
  EXPECT_EQ(UDFExecType::kUDF, info.GetUDFExecType("scalar1").ConsumeValueOrDie());
  EXPECT_NOT_OK(info.GetUDFExecType("dne"));

  EXPECT_EQ(types::INT64, info.GetUDADataType("uda1", std::vector<types::DataType>({types::INT64}))
                              .ConsumeValueOrDie());
  EXPECT_NOT_OK(info.GetUDADataType("uda2", std::vector<types::DataType>({types::INT64})));
  EXPECT_EQ(
      types::INT64,
      info.GetUDFDataType("scalar1", std::vector<types::DataType>({types::BOOLEAN, types::INT64}))
          .ConsumeValueOrDie());
  EXPECT_FALSE(
      info.GetUDFDataType("scalar1", std::vector<types::DataType>({types::BOOLEAN, types::FLOAT64}))
          .ok());
  EXPECT_EQ(
      types::FLOAT64,
      info.GetUDFDataType("add", std::vector<types::DataType>({types::FLOAT64, types::FLOAT64}))
          .ConsumeValueOrDie());

  EXPECT_THAT(info.func_names(), UnorderedElementsAre("uda1", "add", "scalar1"));

  ASSERT_EQ(info.udtfs().size(), 1);
  EXPECT_EQ(info.udtfs()[0].name(), "OpenNetworkConnections");
}

TEST(RegistryInfo, semantic_types) {
  auto info = RegistryInfo();
  udfspb::UDFInfo info_pb;
  google::protobuf::TextFormat::MergeFromString(kExpectedUDFInfo, &info_pb);
  EXPECT_OK(info.Init(info_pb));

  EXPECT_OK_AND_PTR_VAL_EQ(
      info.ResolveUDFType("add", {ValueType::Create(types::FLOAT64, types::ST_BYTES),
                                  ValueType::Create(types::FLOAT64, types::ST_BYTES)}),
      ValueType::Create(types::FLOAT64, types::ST_BYTES));
  EXPECT_OK_AND_PTR_VAL_EQ(
      info.ResolveUDFType("add", {ValueType::Create(types::FLOAT64, types::ST_BYTES),
                                  ValueType::Create(types::FLOAT64, types::ST_BYTES)}),
      ValueType::Create(types::FLOAT64, types::ST_BYTES));

  EXPECT_OK_AND_PTR_VAL_EQ(
      info.ResolveUDFType("uda1", {ValueType::Create(types::INT64, types::ST_BYTES)}),
      ValueType::Create(types::INT64, types::ST_BYTES));
  EXPECT_OK_AND_PTR_VAL_EQ(
      info.ResolveUDFType("uda1", {ValueType::Create(types::INT64, types::ST_UPID)}),
      ValueType::Create(types::INT64, types::ST_UNSPECIFIED));
}

TEST(SemanticRuleRegistry, semantic_lookup) {
  std::vector<types::SemanticType> arg_types1(
      {types::ST_UNSPECIFIED, types::ST_UNSPECIFIED, types::ST_BYTES});
  std::vector<types::SemanticType> arg_types2(
      {types::ST_UPID, types::ST_UNSPECIFIED, types::ST_BYTES});

  std::vector<types::SemanticType> arg_types3(
      {types::ST_UPID, types::ST_SERVICE_NAME, types::ST_BYTES});
  std::vector<types::SemanticType> arg_types4(
      {types::ST_UNSPECIFIED, types::ST_SERVICE_NAME, types::ST_BYTES});

  SemanticRuleRegistry map_;
  map_.Insert("test", arg_types1, types::ST_POD_NAME);
  map_.Insert("test", arg_types2, types::ST_BYTES);

  auto out_type_or_s = map_.Lookup("test", arg_types3);
  ASSERT_OK(out_type_or_s);
  EXPECT_EQ(types::ST_BYTES, out_type_or_s.ConsumeValueOrDie());

  out_type_or_s = map_.Lookup("test", arg_types4);
  ASSERT_OK(out_type_or_s);
  EXPECT_EQ(types::ST_POD_NAME, out_type_or_s.ConsumeValueOrDie());
}

}  // namespace planner
}  // namespace carnot
}  // namespace pl
