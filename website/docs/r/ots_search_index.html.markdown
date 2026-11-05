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
  name        = "tf-example-ots1"
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
  defined_column {
    name = "col1"
    type = "String"
  }
  defined_column {
    name = "col2"
    type = "Integer"
  }
  defined_column {
    name = "col3"
    type = "String"
  }
  defined_column {
    name = "col4"
    type = "String"
  }
  defined_column {
    name = "col5"
    type = "String"
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
    field_schema {
      field_name = "col3"
      field_type = "Vector"
      vector_options {
        data_type   = "float32"
        dimension   = 4
        metric_type = "cosine"
      }
    }
    field_schema {
      field_name = "col4"
      field_type = "FuzzyKeyword"
    }
    field_schema {
      field_name       = "col5"
      field_type       = "Date"
      date_format_list = ["yyyy-MM-dd HH:mm:ss", "yyyy-MM-dd"]
    }
    field_schema {
      field_name        = "col6"
      field_type        = "Keyword"
      is_virtual_field  = true
      source_field_name = "pk2"
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

### Example Usage with Nested Fields

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

resource "alicloud_ots_instance" "default" {
  name          = "tf-example-ots2"
  description   = "tf-example-nested"
  accessed_by   = "Any"
  instance_type = "Capacity"
}

resource "alicloud_ots_table" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = "tf_example_nested"
  time_to_live  = -1
  max_version   = 1
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  defined_column {
    name = "nested_col"
    type = "String"
  }
}

resource "alicloud_ots_search_index" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = alicloud_ots_table.default.table_name
  index_name    = "example_nested_index"
  time_to_live  = -1
  schema {
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
      }
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
* `field_schema` - (Required, ForceNew) A list of field schemas. See [`field_schema`](#schema-field_schema) below.
* `index_setting` - (Optional, ForceNew) The settings of the search index, including routingFields. See [`index_setting`](#schema-index_setting) below.
* `index_sort` - (Optional, ForceNew) The presorting settings of the search index, including sorters. If no value is specified for the indexSort parameter, field values are sorted by primary key by default. See [`index_sort`](#schema-index_sort) below.

### `schema-field_schema`

The field_schema supports the following:
* `field_name` - (Required, ForceNew) Specifies the name of the field in the search index. The value is used as a column name. A field in a search index can be a primary key column or an attribute column.
* `field_type` - (Required, ForceNew, String) Specifies the type of the field. Valid values: `Text`, `Long`, `Double`, `Boolean`, `Keyword`, `FuzzyKeyword`, `Date`, `GeoPoint`, `Nested`, `Vector`.
* `is_array` - (Optional, ForceNew, Bool) Specifies whether the value is an array. This parameter is not supported for `Nested` fields. It is not required for `Vector` fields because they are arrays by nature.
* `index` - (Optional, ForceNew) Specifies whether to enable indexing for the column. Type: Boolean. This parameter is not supported for `Nested` fields.
* `analyzer` - (Optional, ForceNew) Specifies the type of the analyzer that you want to use. If fieldType is set to Text, you can configure this parameter. Otherwise, the default analyzer type single-word tokenization is used.
* `enable_sort_and_agg` - (Optional, ForceNew, Bool) Specifies whether to enable sorting and aggregation. Sorting can be enabled only for fields for which enable_sort_and_agg is set to true. `Text`, `FuzzyKeyword` and `Nested` fields do not support this parameter, but the sub-fields of a `Nested` field do.
* `store` - (Optional, ForceNew) Specifies whether to store the value of the field in the search index. Type: Boolean. If you set store to true, you can read the value of the field from the search index without querying the data table. This improves query performance. This parameter is not supported for `Nested` fields.
* `date_format_list` - (Optional, ForceNew, Set) Specifies the date formats supported by the field. This parameter is required for `Date` fields. Available since v2.3.0.
* `is_virtual_field` - (Optional, ForceNew, Bool) Specifies whether the field is a virtual column. Defaults to `false`. When you set this parameter to `true`, `source_field_name` is required. Available since v2.3.0.
* `source_field_name` - (Optional, ForceNew, String) Specifies the source field to which the virtual column maps. This parameter is required for virtual columns, and only one source field is supported. This parameter can be set only when `is_virtual_field` is `true`. Available since v2.3.0.
* `sub_field_schemas` - (Optional, ForceNew, Set) Specifies the sub-fields of a `Nested` field. This parameter is required for `Nested` fields. See [`sub_field_schemas`](#schema-field_schema-sub_field_schemas) below. A sub-field cannot contain a further `sub_field_schemas` block. Available since v2.3.0.
* `vector_options` - (Optional, ForceNew, List) Specifies the vector configuration. This parameter is required for `Vector` fields. See [`vector_options`](#schema-field_schema-vector_options) below. Available since v2.3.0.

### `schema-field_schema-vector_options`

The vector_options supports the following:
* `data_type` - (Required, ForceNew, String) Specifies the vector data type. Valid values: `float32`. Available since v2.3.0.
* `dimension` - (Required, ForceNew, Int) Specifies the vector dimension. Valid values: 1-4096. Available since v2.3.0.
* `metric_type` - (Required, ForceNew, String) Specifies the distance metric that is used between vectors. Valid values: `euclidean`, `cosine`, `dot_product`. Available since v2.3.0.

### `schema-field_schema-sub_field_schemas`

The sub_field_schemas supports the following:
* `field_name` - (Required, ForceNew, String) Specifies the name of the sub-field. Available since v2.3.0.
* `field_type` - (Required, ForceNew, String) Specifies the type of the sub-field. Valid values: `Text`, `Long`, `Double`, `Boolean`, `Keyword`, `FuzzyKeyword`, `Date`, `GeoPoint`, `Vector`. Available since v2.3.0.
* `is_array` - (Optional, ForceNew, Bool) Specifies whether the value is an array. Available since v2.3.0.
* `index` - (Optional, ForceNew, Bool) Specifies whether to enable indexing for the sub-field. Available since v2.3.0.
* `analyzer` - (Optional, ForceNew, String) Specifies the type of the analyzer that you want to use. Available since v2.3.0.
* `enable_sort_and_agg` - (Optional, ForceNew, Bool) Specifies whether to enable sorting and aggregation for the sub-field. Available since v2.3.0.
* `store` - (Optional, ForceNew, Bool) Specifies whether to store the value of the sub-field in the search index. Available since v2.3.0.
* `date_format_list` - (Optional, ForceNew, Set) Specifies the date formats supported by the sub-field. This parameter is required for `Date` sub-fields. Available since v2.3.0.
* `is_virtual_field` - (Optional, ForceNew, Bool) Specifies whether the sub-field is a virtual column. Defaults to `false`. Available since v2.3.0.
* `source_field_name` - (Optional, ForceNew, String) Specifies the source field to which the virtual sub-field maps. Available since v2.3.0.
* `vector_options` - (Optional, ForceNew, List) Specifies the vector configuration. This parameter is required for `Vector` sub-fields. See [`vector_options`](#schema-field_schema-sub_field_schemas-vector_options) below. Available since v2.3.0.

### `schema-field_schema-sub_field_schemas-vector_options`

The vector_options supports the following:
* `data_type` - (Required, ForceNew, String) Specifies the vector data type. Valid values: `float32`. Available since v2.3.0.
* `dimension` - (Required, ForceNew, Int) Specifies the vector dimension. Valid values: 1-4096. Available since v2.3.0.
* `metric_type` - (Required, ForceNew, String) Specifies the distance metric that is used between vectors. Valid values: `euclidean`, `cosine`, `dot_product`. Available since v2.3.0.

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
