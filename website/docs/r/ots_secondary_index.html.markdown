---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_secondary_index"
sidebar_current: "docs-alicloud-resource-ots-secondary-index"
description: |-
  Provides an OTS (Open Table Service) secondary index resource.
---

# alicloud\_ots\_secondary_index

Provides an OTS secondary index resource.

For information about OTS secondary index and how to use it, see [Secondary index overview](https://www.alibabacloud.com/help/en/tablestore/latest/secondary-index-overview).

-> **NOTE:** Available in v1.187.0+.

## Example Usage

``` terraform
variable "name" {
  default = "terraformtest"
}

variable "pks" {
    default = ["pk1", "pk2", "pk3"]
    type    = list(string)
}

variable "defined_cols" {
    default = ["col1", "col2", "col3"]
    type    = list(string)
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
  primary_key {
    name = "pk3"
    type = "Binary"
  }

  defined_column {
    name = "col1"
    type = "Integer"
  }

  defined_column {
    name = "col2"
    type = "String"
  }

  defined_column {
    name = "col3"
    type = "Binary"
  }
  

  time_to_live                  = -1
  max_version                   = 1
  deviation_cell_version_in_sec = 1
}

resource "alicloud_ots_secondary_index" "index1" {
  instance_name = alicloud_ots_instance.instance1.name
  table_name = alicloud_ots_table.table1.table_name

  index_name = var.name
  index_type = "Global"
  include_base_data = true
  primary_keys = var.pks
  defined_columns = var.defined_cols
}
```

## Argument Reference

The following arguments are supported:
* `instance_name` - (Required, ForceNew) The name of the OTS instance in which table will located.
* `table_name` - (Required, ForceNew) The name of the OTS table. If changed, a new table would be created.
* `index_name` - (Required, ForceNew) The index name of the OTS Table. If changed, a new index would be created.
* `index_type` - (Required, ForceNew) The index type of the OTS Table. If changed, a new index would be created, only `Global` or `Local` is allowed.
* `include_base_data` - (Required, ForceNew) whether the index contains data that already exists in the data table. When include_base_data is set to true, it means that stock data is included.
* `primary_keys` - (Required, ForceNew) A list of primary keys for index, referenced from Table's primary keys or predefined columns.
* `defined_columns` - (Optional, ForceNew) A list of defined column for index, referenced from Table's primary keys or predefined columns.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is `<instance_name>:<table_name>:<indexName>:<indexType>`.
* `primary_keys` - The primary keys of OTS secondary index. Each element contains the following attributes:
    * `name` - The name of the key.
    * `type` - The type of the key, valid values: `Integer`, `Binary`, `String`.
* `defined_columns` - The defined columns of OTS secondary index. Each element contains the following attributes:
  * `name` - The name of the defined columns.
  * `type` - The type of the defined columns, valid values: `Integer`, `Binary`, `String`, `Double`, `Boolean`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the OTS secondary index.
* `delete` - (Defaults to 2 mins) Used when delete the OTS secondary index.

## Import

OTS secondary index can be imported using id, e.g.

```shell
$ terraform import alicloud_ots_secondary_index.index1 "<instance_name>:<table_name>:<index_name>:<index_type>"
```
