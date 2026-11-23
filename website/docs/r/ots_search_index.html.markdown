---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_search_index"
sidebar_current: "docs-alicloud-resource-ots-search-index"
description: |-
  Provides an OTS (Open Table Service) search index resource.
---

# alicloud_ots_search_index

Provides an OTS search index resource.

For information about OTS search index and how to use it, see [Search index overview](https://www.alibabacloud.com/help/en/tablestore/latest/search-index-overview).

-> **NOTE:** Available since v1.187.0.

-> **NOTE:** The examples use short fixed OTS instance names (`tf-example-ots2`, `tf-example-ots3`) to fit the 16-byte instance name limit. Destroy the instance before running an example again with the same account, or change the name.

-> **NOTE:** The search index API does not return every configured value, so the imported `schema` is a best-effort reconstruction. `store` is returned only when the service keeps it, the default `PrimaryKeySort` sorter is added when the index has no custom sort, and the `mode` of a `FieldSort` sorter cannot be read back. Review the plan after importing and align the configuration with the imported state; the importer only fills `schema` when the state does not have one, and normal refreshes never rewrite the configured schema.

-> **NOTE:** `schema.field_schema` remains a list in configuration, but the Provider ignores only the top-level field order when comparing the enclosing `schema` block. Changing only the field order does not change the index or propose replacement; changing field membership or values still forces replacement. Order normalization changes the hash used to represent the existing `schema` state, so review the first plan after upgrading the Provider before applying it.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ots_search_index&exampleId=8f64e5b6-4cbb-a79f-5e21-6717a0a2d9ae80e0a156&activeTab=example&spm=docs.r.ots_search_index.0.8f64e5b64c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "tf-example"
}

resource "alicloud_ots_instance" "default" {
  name        = "tf-example-ots3"
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_ots_table" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = "tf_example"
  time_to_live  = -1
  max_version   = 1
  enable_sse    = true
  sse_key_type  = "SSE_KMS_SERVICE"
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  primary_key {
    name = "pk2"
    type = "String"
  }
  primary_key {
    name = "pk3"
    type = "Binary"
  }
}

resource "alicloud_ots_search_index" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = alicloud_ots_table.default.table_name
  index_name    = "example_index"
  time_to_live  = -1
  schema {
    field_schema {
      field_name = "col1"
      field_type = "Text"
      is_array   = false
      index      = true
      analyzer   = "Split"
      store      = true
    }
    field_schema {
      field_name          = "col2"
      field_type          = "Long"
      enable_sort_and_agg = true
    }
    field_schema {
      field_name = "pk1"
      field_type = "Long"
    }
    field_schema {
      field_name = "pk2"
      field_type = "Text"
    }

    index_setting {
      routing_fields = ["pk1", "pk2"]
    }

    index_sort {
      sorter {
        sorter_type = "PrimaryKeySort"
        order       = "Asc"
      }
      sorter {
        sorter_type = "FieldSort"
        order       = "Desc"
        field_name  = "col2"
        mode        = "Max"
      }
    }
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ots_search_index&spm=docs.r.ots_search_index.example&intl_lang=EN_US)

### Search Index with Nested, Vector, Date and Virtual Fields

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

resource "alicloud_ots_instance" "default" {
  name        = "tf-example-ots2"
  description = "tf-example-ots2"
  accessed_by = "Any"
}

resource "alicloud_ots_table" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = "terraform_example_nested"
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  primary_key {
    name = "pk2"
    type = "String"
  }
  defined_column {
    name = "vector_col"
    type = "String"
  }
  defined_column {
    name = "fuzzy_col"
    type = "String"
  }
  defined_column {
    name = "nested_col"
    type = "String"
  }
  defined_column {
    name = "date_col"
    type = "String"
  }
  time_to_live = -1
  max_version  = 1
}

resource "alicloud_ots_search_index" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = alicloud_ots_table.default.table_name
  index_name    = "terraform_example_nested_index"
  time_to_live  = -1

  schema {
    field_schema {
      field_name = "vector_col"
      field_type = "Vector"
      vector_options {
        data_type   = "float32"
        dimension   = 4
        metric_type = "cosine"
      }
    }
    field_schema {
      field_name = "fuzzy_col"
      field_type = "FuzzyKeyword"
    }
    field_schema {
      field_name = "nested_col"
      field_type = "Nested"
      sub_field_schemas {
        field_name          = "nested_long"
        field_type          = "Long"
        enable_sort_and_agg = true
      }
      sub_field_schemas {
        field_name = "nested_text"
        field_type = "Text"
        analyzer   = "SingleWord"
        store      = true
      }
      sub_field_schemas {
        field_name        = "nested_virtual"
        field_type        = "Keyword"
        is_virtual_field  = true
        source_field_name = "pk2"
      }
    }
    field_schema {
      field_name       = "date_col"
      field_type       = "Date"
      date_format_list = ["yyyy-MM-dd HH:mm:ss", "yyyy-MM-dd"]
    }
    field_schema {
      field_name        = "virtual_keyword"
      field_type        = "Keyword"
      is_virtual_field  = true
      source_field_name = "pk2"
    }
    field_schema {
      field_name = "pk2"
      field_type = "Text"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, ForceNew) The name of the OTS instance in which table will located.
* `table_name` - (Required, ForceNew) The name of the OTS table. If changed, a new table would be created.
* `index_name` - (Required, ForceNew) The index name of the OTS Table. If changed, a new index would be created.
* `time_to_live` - (Optional, ForceNew) The index type of the OTS Table. Specifies the retention period of data in the search index. Unit: seconds. Default value: -1.
  If the retention period exceeds the TTL value, OTS automatically deletes expired data.
* `schema` - (Required, ForceNew) The schema of the search index. If changed, a new index would be created. See [`schema`](#schema) below.

### `schema`

The schema supports the following:

* `field_schema` - (Required, ForceNew, List) A list of field schemas. The Provider ignores field order when comparing the enclosing `schema` block. See [`field_schema`](#schema-field_schema) below.
* `index_setting` - (Optional, ForceNew) The settings of the search index, including routingFields. See [`index_setting`](#schema-index_setting) below.
* `index_sort` - (Optional, ForceNew) The presorting settings of the search index, including sorters. If no value is specified for the indexSort parameter, field values are sorted by primary key by default. See [`index_sort`](#schema-index_sort) below.

### `schema-field_schema`

The field_schema supports the following:

* `field_name` - (Required, ForceNew) Specifies the name of the field in the search index. The value is used as a column name. A field in a search index can be a primary key column or an attribute column.
* `field_type` - (Required, ForceNew) Specifies the type of the field. Type: String. Valid values: `Text`, `Long`, `Double`, `Boolean`, `Keyword`, `Date`, `GeoPoint`, `FuzzyKeyword`, `Nested`, `Vector`. The Provider supports `Nested` and `Vector` only as top-level fields.
* `is_array` - (Optional, ForceNew) Specifies whether the value is an array. Type: Boolean. This parameter is not supported for `Nested`, `Vector` or `FuzzyKeyword` fields.
* `index` - (Optional, ForceNew) Specifies whether to enable indexing for the column. Type: Boolean. Defaults to `true`. A `Nested` container may set this to `false`; virtual fields must set it to `true`.
* `analyzer` - (Optional, ForceNew) Specifies the type of the analyzer that you want to use. Type: String. If fieldType is set to Text, you can configure this parameter. Otherwise, the default analyzer type single-word tokenization is used. Valid values: `SingleWord`, `Split`, `MinWord`, `MaxWord`, `Fuzzy`.
* `enable_sort_and_agg` - (Optional, ForceNew) Specifies whether to enable sorting and aggregation. Type: Boolean. Sorting can be enabled only for fields for which enable_sort_and_agg is set to true. `Nested`, `Text` and `FuzzyKeyword` fields do not support sorting or aggregation.
* `store` - (Optional, ForceNew) Specifies whether to store the value of the field in the search index. Type: Boolean. If you set store to true, you can read the value of the field from the search index without querying the data table. This improves query performance. This parameter is not supported for `Nested` fields. When omitted, `Vector`, `FuzzyKeyword` and virtual fields follow the SDK default (`false` for `Text`, `true` for the other types), while the other top-level field types (all except `Nested`, which does not support `store`) keep the previous behavior and are created with `store = false`.
* `date_format_list` - (Optional, ForceNew, Set) Specifies the date formats supported by a `Date` field; other field types do not support this argument. For backward compatibility, the Provider preserves the legacy passthrough behavior for a top-level `Date` field that omits this argument and does not add a local validation error. The OTS service currently requires a format for `Date` fields and may reject an omitted value. A `Date` sub-field of a `Nested` field must configure this argument. Use `epoch_second`, `epoch_millis`, `epoch_micros` or `epoch_nanos` for numeric timestamps, or a custom format matching string date values such as `yyyy-MM-dd HH:mm:ss.SSS`. The Provider does not validate the format strings and the service performs the final check. Available since v2.4.0.
* `is_virtual_field` - (Optional, ForceNew) Specifies whether the field is a virtual column. Type: Boolean. When set to `true`, `source_field_name` is required and `index` must be `true`. This parameter is not supported for `Nested` fields. Available since v2.4.0.
* `source_field_name` - (Optional, ForceNew) Specifies the data table column that provides the values of a virtual column. Type: String. Exactly one source column is supported, it must exist in the data table, and its type must map to the target search index field type. Only valid together with `is_virtual_field = true`. See [Virtual columns](https://help.aliyun.com/zh/tablestore/table-store-service/virtual-columns). Available since v2.4.0.
* `vector_options` - (Optional, ForceNew, List) The vector settings of a field whose `field_type` is `Vector`. Required for `Vector` fields and not supported for other field types. See [`vector_options`](#schema-field_schema-vector_options) below. Available since v2.4.0.
* `sub_field_schemas` - (Optional, ForceNew, Set) The sub-fields of a top-level field whose `field_type` is `Nested`; required in that case. Only one nesting level is supported. See [`sub_field_schemas`](#schema-field_schema-sub_field_schemas) below. Available since v2.4.0.

  -> **NOTE:** The Provider rejects attributes that do not apply to a `Nested` container itself (`is_array = true`, `store = true`, a non-empty `analyzer`, `enable_sort_and_agg = true` and `is_virtual_field = true`). `Vector` fields are not supported as `Nested` sub-fields.

### `schema-field_schema-vector_options`

The vector_options supports the following:

* `data_type` - (Optional, ForceNew) The data type of the vector. Type: String. Valid values: `float32`. Required by the service. Available since v2.4.0.
* `dimension` - (Optional, ForceNew) The number of dimensions of the vector. Type: Int. Valid values: 1 to 4096. Required by the service. Available since v2.4.0.
* `metric_type` - (Optional, ForceNew) The similarity metric. Type: String. Valid values: `euclidean`, `cosine`, `dot_product`. Required by the service. Available since v2.4.0.

### `schema-field_schema-sub_field_schemas`

The sub_field_schemas describes one sub-field of a `Nested` field and supports the following:

* `field_name` - (Required, ForceNew) Specifies the name of the sub-field in the search index. Type: String. Available since v2.4.0.
* `field_type` - (Required, ForceNew) Specifies the type of the sub-field. Type: String. Valid values: `Text`, `Long`, `Double`, `Boolean`, `Keyword`, `Date`, `GeoPoint`, `FuzzyKeyword`. The Provider does not support `Nested` or `Vector` sub-fields. Available since v2.4.0.
* `is_array` - (Optional, ForceNew) Specifies whether the value of the sub-field is an array. Type: Boolean. `FuzzyKeyword` sub-fields do not support arrays. Available since v2.4.0.
* `index` - (Optional, ForceNew) Specifies whether to enable indexing for the sub-field. Type: Boolean. Defaults to `true`. Available since v2.4.0.
* `analyzer` - (Optional, ForceNew) Specifies the type of the analyzer of the sub-field. Type: String. Applies to `Text` sub-fields. Valid values: `SingleWord`, `Split`, `MinWord`, `MaxWord`, `Fuzzy`. Available since v2.4.0.
* `enable_sort_and_agg` - (Optional, ForceNew) Specifies whether to enable sorting and aggregation for the sub-field. Type: Boolean. `Text` and `FuzzyKeyword` sub-fields do not support sorting or aggregation. Available since v2.4.0.
* `store` - (Optional, ForceNew) Specifies whether to store the value of the sub-field in the search index. Type: Boolean. When omitted, the field follows the SDK default: `false` for `Text` and `true` for the other supported sub-field types. Available since v2.4.0.
* `date_format_list` - (Optional, ForceNew, Set) Specifies the date formats supported by a `Date` sub-field; other sub-field types do not support this argument. This parameter is required when the sub-field `field_type` is `Date`. Use `epoch_second`, `epoch_millis`, `epoch_micros` or `epoch_nanos` for numeric timestamps, or a custom format matching string date values such as `yyyy-MM-dd HH:mm:ss.SSS`. The Provider does not validate the format strings and the service performs the final check. Available since v2.4.0.
* `is_virtual_field` - (Optional, ForceNew) Specifies whether the sub-field is a virtual column. Type: Boolean. When set to `true`, `source_field_name` is required and `index` must be `true`. Available since v2.4.0.
* `source_field_name` - (Optional, ForceNew) Specifies the data table column that provides the values of a virtual sub-field. Type: String. Exactly one source column is supported, it must exist in the data table, and its type must map to the target search index sub-field type. Only valid together with `is_virtual_field = true`. See [Virtual columns](https://help.aliyun.com/zh/tablestore/table-store-service/virtual-columns). Available since v2.4.0.

### `schema-index_setting`

The index_setting supports the following:

* `routing_fields` - (Optional, ForceNew) Specifies custom routing fields. You can specify some primary key columns as routing fields. Tablestore distributes data that is written to a search index across different partitions based on the specified routing fields. The data whose routing field values are the same is distributed to the same partition.

### `schema-index_sort`

The index_sort supports the following:

* `sorter` - (Required, ForceNew)  Specifies the presorting method for the search index. PrimaryKeySort and FieldSort are supported. See [`sorter`](#schema-index_sort-sorter) below.

### `schema-index_sort-sorter`

The sorter supports the following:

* `sorter_type` - (Optional, ForceNew) Data is sorted by Which fields or keys. valid values: `PrimaryKeySort`, `FieldSort`.
* `order` - (Optional, ForceNew) The sort order. Data can be sorted in ascending(`Asc`) or descending(`Desc`) order. Default value: `Asc`.
* `field_name` - (Optional, ForceNew) The name of the field that is used to sort data. only required if sorter_type is FieldSort.
* `mode` - (Optional, ForceNew) The sorting method that is used when the field contains multiple values. valid values: `Min`, `Max`, `Avg`. only required if sorter_type is FieldSort.

## Attributes Reference

The following attributes are exported:

* `index_id` - The index id of the search index which could not be changed.
* `create_time` - The search index create time.
* `sync_phase` - The search index sync phase. possible values: `Full`, `Incr`.
* `current_sync_timestamp` - The timestamp for sync phase.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the OTS search index.
* `delete` - (Defaults to 2 mins) Used when delete the OTS search index.

## Import

OTS search index can be imported using id, e.g.

```shell
$ terraform import alicloud_ots_search_index.index1 <instance_name>:<table_name>:<index_name>:<index_type>
```

The importer reads `schema`, `instance_name`, `table_name`, `index_name` and `time_to_live` from the existing search index.
