---
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_table"
sidebar_current: "docs-alicloud-resource-ots-table"
description: |-
  Provides an OTS (Open Table Service) table resource.
---

# alicloud\_ots\_table

Provides an OTS table resource.

~> **NOTE:** From Provider version 1.9.7, the provider field 'ots_instance_name' has been deprecated and
you should use resource alicloud_ots_table's new field 'instance_name' and 'table_name' to re-import this resource.

## Example Usage

```
# Create an OTS table

resource "alicloud_ots_instance" "foo" {
  name = "my-ots"
  description = "ots instance"
  accessed_by = "Any"
  tags {
    Created = "TF"
    For = "acceptance test"
  }
}

resource "alicloud_ots_table" "table" {
  instance_name = "${alicloud_ots_instance.foo.name}"
  table_name = "ots-table"
  primary_key = [
    {
      name = "${var.primary_key_1_name}"
      type = "${var.primary_key_integer_type}"
    },
    {
      name = "${var.primary_key_2_name}"
      type = "${var.primary_key_integer_type}"
    },
    {
      name = "${var.primary_key_3_name}"
      type = "${var.primary_key_integer_type}"
    },
    {
      name = "${var.primary_key_4_name}"
      type = "${var.primary_key_string_type}"
    },
  ]
  time_to_live = "${var.time_to_live}"
  max_version = "${var.max_version}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, ForceNew) The name of the OTS instance in which table will located.
* `table_name` - (Required, ForceNew) The table name of the OTS instance. If changed, a new table would be created.
* `primary_key` - (Required, Type: List) The property of `TableMeta` which indicates the structure information of a table. It describes the attribute value of primary key. The number of `primary_key` should not be less than one and not be more than four.
    * `name` - (Required) Name for primary key.
    * `type` - (Required, Type: list) Type for primary key. Only `Integer`, `String` or `Binary` is allowed.
* `time_to_live` - (Required) The retention time of data stored in this table (unit: second). The value maximum is 2147483647 and -1 means never expired.
* `max_version` - (Required) The maximum number of versions stored in this table. The valid value is 1-2147483647.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is "<instance_name>:<table_name>".
* `instance_name` - The OTS instance name.
* `table_name` - The table name of the OTS which could not be changed.
* `primary_key` - The property of `TableMeta` which indicates the structure information of a table.
* `time_to_live` - The retention time of data stored in this table.
* `max_version` - The maximum number of versions stored in this table.

## Import

OTS table can be imported using id, e.g.

```
$ terraform import alicloud_ots_table.table "my-ots:ots_table"
```

