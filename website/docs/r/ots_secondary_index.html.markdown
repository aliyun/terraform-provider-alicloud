---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_secondary_index"
sidebar_current: "docs-alicloud-resource-ots-secondary-index"
description: |-
  Provides an OTS (Open Table Service) secondary index resource.
---

# alicloud_ots_secondary_index

Provides an OTS secondary index resource.

For information about OTS secondary index and how to use it, see [Secondary index overview](https://www.alibabacloud.com/help/en/tablestore/latest/secondary-index-overview).

-> **NOTE:** Available since v1.187.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ots_instance" "default" {
  name        = "${var.name}-${random_integer.default.result}"
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
}

resource "alicloud_ots_secondary_index" "default" {
  instance_name     = alicloud_ots_instance.default.name
  table_name        = alicloud_ots_table.default.table_name
  index_name        = "example_index"
  index_type        = "Global"
  include_base_data = true
  primary_keys      = ["pk1", "pk2", "pk3"]
  defined_columns   = ["col1", "col2", "col3"]
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the OTS secondary index.
* `delete` - (Defaults to 2 mins) Used when delete the OTS secondary index.

## Import

OTS secondary index can be imported using id, e.g.

```shell
$ terraform import alicloud_ots_secondary_index.index1 <instance_name>:<table_name>:<index_name>:<index_type>
```
