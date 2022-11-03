---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_search_index"
sidebar_current: "docs-alicloud-resource-ots-search-index"
description: |-
  Provides an OTS (Open Table Service) search index resource.
---

# alicloud\_ots\_search_index

Provides an OTS search index resource.

For information about OTS search index and how to use it, see [Search index overview](https://www.alibabacloud.com/help/en/tablestore/latest/search-index-overview).

-> **NOTE:** Available in v1.187.0+.

## Example Usage

```terraform
variable "name" {
  default = "terraformtest"
}

resource "alicloud_ots_instance" "instance1" {
  name        = var.name
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_ots_table" "table1" {
  instance_name = alicloud_ots_instance.instance1.name
  table_name    = var.name
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  primary_key {
    name = "pk2"
    type = "String"
  }

  defined_column {
    name = "col1"
    type = "String"
  }

  defined_column {
    name = "col2"
    type = "Integer"
  }

  time_to_live                  = -1
  max_version                   = 1
  deviation_cell_version_in_sec = 1
}

resource "alicloud_ots_search_index" "default" {
  instance_name = alicloud_ots_instance.instance1.name
  table_name = alicloud_ots_table.table1.table_name

  index_name = var.name
  time_to_live = -1
  schema {
      field_schema {
        field_name = "col1"
        field_type = "Text"
        is_array = false
        index = true
        analyzer = "Split"
        store = true
      }
      field_schema {
        field_name =  "col2"
         field_type = "Long"
         enable_sort_and_agg = true
      }

      field_schema {
        field_name =  "pk1"
        field_type = "Long"
        
      }
      field_schema {
        field_name =  "pk2"
        field_type = "Text"
      }

      index_setting {
        routing_fields = [ "pk1", "pk2"]
      }

      index_sort {
        sorter {
          sorter_type = "PrimaryKeySort"
          order = "Asc"
        }
        sorter {
          sorter_type = "FieldSort"
          order = "Desc"
          field_name =  "col2"
          mode = "Max"
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
* `schema` - (Required, ForceNew) The schema of the search index. If changed, a new index would be created.
  * `field_schema` - (Required, ForceNew) A list of field schemas. Each field schema contains the following parameters:
    * `field_name` - (Required, ForceNew) Specifies the name of the field in the search index. The value is used as a column name. Type: String.
      A field in a search index can be a primary key column or an attribute column.
    * `field_type` - (Required, ForceNew) Specifies the type of the field. Use FieldType.XXX to set the type.
    * `is_array` - (Optional, ForceNew) Specifies whether the value is an array. Type: Boolean.
    * `index` - (Optional, ForceNew) Specifies whether to enable indexing for the column. Type: Boolean.
    * `analyzer` - (Optional, ForceNew) Specifies the type of the analyzer that you want to use. If fieldType is set to Text, you can configure this parameter. Otherwise, the default analyzer type single-word tokenization is used.
    * `enable_sort_and_agg` - (Optional, ForceNew) Specifies whether to enable sorting and aggregation. Type: Boolean. Sorting can be enabled only for fields for which enable_sort_and_agg is set to true.
    * `store` - (Optional, ForceNew) Specifies whether to store the value of the field in the search index. Type: Boolean. If you set store to true, you can read the value of the field from the search index without querying the data table. This improves query performance.
  * `index_setting` - (Optional, ForceNew) The settings of the search index, including routingFields.
    * `routing_fields` - (Optional, ForceNew) Specifies custom routing fields. You can specify some primary key columns as routing fields. Tablestore distributes data that is written to a search index across different partitions based on the specified routing fields. The data whose routing field values are the same is distributed to the same partition.
  * `index_sort` - (Optional, ForceNew) The presorting settings of the search index, including sorters. If no value is specified for the indexSort parameter, field values are sorted by primary key by default.
    * `sorter` - (Required, ForceNew)  Specifies the presorting method for the search index. PrimaryKeySort and FieldSort are supported.
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the OTS search index.
* `delete` - (Defaults to 2 mins) Used when delete the OTS search index.

## Import

OTS search index can be imported using id, e.g.

```shell
$ terraform import alicloud_ots_search_index.index1 <instance_name>:<table_name>:<index_name>:<index_type>
```
