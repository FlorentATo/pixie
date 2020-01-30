#include <arrow/array.h>
#include <algorithm>
#include <memory>
#include <string>
#include <vector>

#include <absl/strings/str_format.h>
#include "src/common/base/base.h"
#include "src/shared/types/arrow_adapter.h"
#include "src/shared/types/type_utils.h"
#include "src/table_store/schema/row_batch.h"

namespace pl {
namespace table_store {
namespace schema {

using types::DataType;

std::shared_ptr<arrow::Array> RowBatch::ColumnAt(int64_t i) const { return columns_[i]; }

Status RowBatch::AddColumn(const std::shared_ptr<arrow::Array>& col) {
  if (columns_.size() >= desc_.size()) {
    return error::InvalidArgument("Schema only allows $0 columns", desc_.size());
  }
  if (col->length() != num_rows_) {
    return error::InvalidArgument("Schema only allows $0 rows, got $1", num_rows_, col->length());
  }
  if (col->type_id() != types::ToArrowType(desc_.type(columns_.size()))) {
    return error::InvalidArgument("Column[$0] was given incorrect type", columns_.size());
  }

  columns_.emplace_back(col);
  return Status::OK();
}

bool RowBatch::HasColumn(int64_t i) const { return columns_.size() > static_cast<size_t>(i); }

std::string RowBatch::DebugString() const {
  if (columns_.empty()) {
    return "RowBatch: <empty>";
  }
  std::string debug_string = absl::StrFormat("RowBatch(eow=%d, eos=%d):\n", eow_, eos_);
  for (const auto& col : columns_) {
    debug_string += absl::StrFormat("  %s\n", col->ToString());
  }
  return debug_string;
}

int64_t RowBatch::NumBytes() const {
  if (num_rows() == 0) {
    return 0;
  }

  int64_t total_bytes = 0;
  for (auto col : columns_) {
#define TYPE_CASE(_dt_) total_bytes += types::GetArrowArrayBytes<_dt_>(col.get());
    PL_SWITCH_FOREACH_DATATYPE(types::ArrowToDataType(col->type_id()), TYPE_CASE);
#undef TYPE_CASE
  }
  return total_bytes;
}

// Serialize/deserialize from protobuf.
// TODO(nserrino): Replace these functions when schemapb::RowBatchData is deprecated
// in favor of proper (de)serialization.

// PL_CARNOT_UPDATE_FOR_NEW_TYPES
template <DataType T>
constexpr auto GetMutablePBDataColumn(table_store::schemapb::Column* data_col) {
  if constexpr (T == DataType::BOOLEAN) {
    return data_col->mutable_boolean_data();
  } else if constexpr (T == DataType::INT64) {
    return data_col->mutable_int64_data();
  } else if constexpr (T == DataType::TIME64NS) {
    return data_col->mutable_time64ns_data();
  } else if constexpr (T == DataType::DURATION64NS) {
    return data_col->mutable_duration64ns_data();
  } else if constexpr (T == DataType::UINT128) {
    return data_col->mutable_uint128_data();
  } else if constexpr (T == DataType::FLOAT64) {
    return data_col->mutable_float64_data();
  } else if constexpr (T == DataType::STRING) {
    return data_col->mutable_string_data();
  } else {
    static_assert(sizeof(T) != 0, "Unsupported data type");
  }
}

template <DataType T>
constexpr const auto& GetPBDataColumn(const table_store::schemapb::Column& data_col) {
  if constexpr (T == DataType::BOOLEAN) {
    return data_col.boolean_data();
  } else if constexpr (T == DataType::INT64) {
    return data_col.int64_data();
  } else if constexpr (T == DataType::TIME64NS) {
    return data_col.time64ns_data();
  } else if constexpr (T == DataType::DURATION64NS) {
    return data_col.duration64ns_data();
  } else if constexpr (T == DataType::UINT128) {
    return data_col.uint128_data();
  } else if constexpr (T == DataType::FLOAT64) {
    return data_col.float64_data();
  } else if constexpr (T == DataType::STRING) {
    return data_col.string_data();
  } else {
    static_assert(sizeof(T) != 0, "Unsupported data type");
  }
}

template <DataType T>
void CopyIntoOutputPB(table_store::schemapb::Column* output_column, arrow::Array* input_column) {
  CHECK_NOTNULL(input_column);
  CHECK_NOTNULL(output_column);

  size_t col_length = input_column->length();
  auto casted_output_data = GetMutablePBDataColumn<T>(output_column);
  for (size_t i = 0; i < col_length; ++i) {
    if constexpr (T == DataType::UINT128) {
      auto out_datum = casted_output_data->add_data();
      auto val = types::GetValueFromArrowArray<DataType::UINT128>(input_column, i);
      out_datum->set_high(absl::Uint128High64(val));
      out_datum->set_low(absl::Uint128Low64(val));
    } else {
      casted_output_data->add_data(types::GetValueFromArrowArray<T>(input_column, i));
    }
  }
}

template <DataType T>
Status CopyFromInputPB(std::shared_ptr<arrow::Array>* output_column,
                       const table_store::schemapb::Column& input_column) {
  CHECK_NOTNULL(output_column);

  auto builder = MakeArrowBuilder(T, arrow::default_memory_pool());
  auto input_data = GetPBDataColumn<T>(input_column);
  PL_RETURN_IF_ERROR(builder->Reserve(input_data.data_size()));

  for (const auto& datum : input_data.data()) {
    if constexpr (T == DataType::UINT128) {
      PL_RETURN_IF_ERROR(CopyValue<T>(builder.get(), types::UInt128Value(datum).val));
    } else {
      PL_RETURN_IF_ERROR(CopyValue<T>(builder.get(), datum));
    }
  }
  PL_RETURN_IF_ERROR(builder->Finish(output_column));
  return Status::OK();
}

Status RowBatch::ToProto(table_store::schemapb::RowBatchData* proto) const {
  proto->set_num_rows(num_rows_);
  proto->set_eow(eow_);
  proto->set_eos(eos_);

  for (auto col_idx = 0; col_idx < num_columns(); ++col_idx) {
    auto input_col = ColumnAt(col_idx).get();
    auto output_col_data = proto->add_cols();
    auto dt = desc_.type(col_idx);

#define TYPE_CASE(_dt_) CopyIntoOutputPB<_dt_>(output_col_data, input_col);
    PL_SWITCH_FOREACH_DATATYPE(dt, TYPE_CASE);
#undef TYPE_CASE
  }

  return Status::OK();
}

StatusOr<DataType> ProtoDataType(const table_store::schemapb::Column& proto) {
  switch (proto.col_data_case()) {
    case table_store::schemapb::Column::kBooleanData:
      return DataType::BOOLEAN;
    case table_store::schemapb::Column::kInt64Data:
      return DataType::INT64;
    case table_store::schemapb::Column::kUint128Data:
      return DataType::UINT128;
    case table_store::schemapb::Column::kTime64NsData:
      return DataType::TIME64NS;
    case table_store::schemapb::Column::kFloat64Data:
      return DataType::FLOAT64;
    case table_store::schemapb::Column::kStringData:
      return DataType::STRING;
    default:
      return error::Internal("Received unknown column data type (in ProtoDataType)");
  }
}

StatusOr<std::unique_ptr<RowBatch>> RowBatch::FromProto(
    const table_store::schemapb::RowBatchData& proto) {
  std::vector<DataType> types(proto.cols_size());
  std::vector<std::shared_ptr<arrow::Array>> data_columns(proto.cols_size());

  for (auto i = 0; i < proto.cols_size(); ++i) {
    PL_ASSIGN_OR_RETURN(types[i], ProtoDataType(proto.cols(i)));
    std::shared_ptr<arrow::Array> output_array;

#define TYPE_CASE(_dt_) PL_RETURN_IF_ERROR(CopyFromInputPB<_dt_>(&data_columns[i], proto.cols(i)));
    PL_SWITCH_FOREACH_DATATYPE(types[i], TYPE_CASE);
#undef TYPE_CASE
  }

  RowDescriptor desc(types);
  std::unique_ptr<RowBatch> output_rb = std::make_unique<RowBatch>(desc, proto.num_rows());
  output_rb->set_eow(proto.eow());
  output_rb->set_eos(proto.eos());

  for (auto i = 0; i < proto.cols_size(); ++i) {
    PL_RETURN_IF_ERROR(output_rb->AddColumn(data_columns[i]));
  }

  return output_rb;
}

StatusOr<std::unique_ptr<RowBatch>> RowBatch::FromColumnBuilders(
    const RowDescriptor& desc, bool eow, bool eos,
    std::vector<std::unique_ptr<arrow::ArrayBuilder>>* builders) {
  DCHECK(builders->size());
  int64_t output_rows = builders->at(0)->length();

  auto output_rb = std::make_unique<RowBatch>(desc, output_rows);
  output_rb->set_eow(eow);
  output_rb->set_eos(eos);

  for (auto& column_builder : *builders) {
    std::shared_ptr<arrow::Array> output_array;
    PL_RETURN_IF_ERROR(column_builder->Finish(&output_array));
    PL_RETURN_IF_ERROR(output_rb->AddColumn(output_array));
  }

  return output_rb;
}

StatusOr<std::unique_ptr<RowBatch>> RowBatch::WithZeroRows(const RowDescriptor& desc, bool eow,
                                                           bool eos) {
  std::vector<std::unique_ptr<arrow::ArrayBuilder>> builders(desc.size());
  for (size_t i = 0; i < desc.size(); ++i) {
    builders[i] = types::MakeArrowBuilder(desc.type(i), arrow::default_memory_pool());
    PL_RETURN_IF_ERROR(builders[i]->Reserve(0));
  }
  return RowBatch::FromColumnBuilders(desc, eow, eos, &builders);
}

}  // namespace schema
}  // namespace table_store
}  // namespace pl
