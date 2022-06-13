---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_table"
sidebar_current: "docs-alicloud-resource-ots-table"
description: |-
  Provides an OTS (Open Table Service) table resource.
---

# alicloud\_ots\_table

Provides an OTS table resource.

-> **NOTE:** From Provider version 1.10.0, the provider field 'ots_instance_name' has been deprecated and
you should use resource alicloud_ots_table's new field 'instance_name' and 'table_name' to re-import this resource.

## Example Usage

```
variable "name" {
  default = "terraformtest"
}

resource "alicloud_ots_instance" "foo" {
  name        = var.name
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_ots_table" "basic" {
  instance_name = alicloud_ots_instance.foo.name
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

  time_to_live                  = -1
  max_version                   = 1
  deviation_cell_version_in_sec = 1
  enable_sse                    = true
  sse_key_type                  = "SSE_KMS_SERVICE"
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, ForceNew) The name of the OTS instance in which table will located.
* `table_name` - (Required, ForceNew) The table name of the OTS instance. If changed, a new table would be created.
* `primary_key` - (Required, ForceNew) The property of `TableMeta` which indicates the structure information of a table. It describes the attribute value of primary key. The number of `primary_key` should not be less than one and not be more than four.
    * `name` - (Required, ForceNew) Name for primary key.
    * `type` - (Required, ForceNew) Type for primary key. Only `Integer`, `String` or `Binary` is allowed.
* `time_to_live` - (Required) The retention time of data stored in this table (unit: second). The value maximum is 2147483647 and -1 means never expired.
* `max_version` - (Required) The maximum number of versions stored in this table. The valid value is 1-2147483647.
* `deviation_cell_version_in_sec` - (Optional, Available in 1.42.0+) The max version offset of the table. The valid value is 1-9223372036854775807. Defaults to 86400.
* `enable_sse` - (Optional, Available in 1.172.0+) Whether enable OTS server side encryption. Default value is false.
* `sse_key_type` - (Optional, Available in 1.172.0+) The key type of OTS server side encryption. Only `SSE_KMS_SERVICE` is allowed.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is `<instance_name>:<table_name>`.
* `instance_name` - The OTS instance name.
* `table_name` - The table name of the OTS which could not be changed.
* `primary_key` - The property of `TableMeta` which indicates the structure information of a table.
* `time_to_live` - The retention time of data stored in this table.
* `max_version` - The maximum number of versions stored in this table.
* `deviation_cell_version_in_sec` - The max version offset of the table.
* `enable_sse` - Whether enable OTS server side encryption.
* `sse_key_type` - The key type of OTS server side encryption.

## Import

OTS table can be imported using id, e.g.

```
$ terraform import alicloud_ots_table.table "my-ots:ots_table"
```

