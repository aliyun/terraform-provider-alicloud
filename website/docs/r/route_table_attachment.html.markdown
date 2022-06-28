---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_route_table_attachment"
sidebar_current: "docs-alicloud-resource-route-table-attachment"
description: |-
  Provides an Alicloud Route Table Attachment resource.
---

# alicloud\_route\_table\_attachment

Provides an Alicloud Route Table Attachment resource for associating Route Table to VSwitch Instance.

-> **NOTE:** Terraform will auto build route table attachment while it uses `alicloud_route_table_attachment` to build a route table attachment resource.

For information about route table and how to use it, see [What is Route Table](https://www.alibabacloud.com/help/doc-detail/87057.htm).

## Example Usage

Basic Usage

```
variable "name" {
  default = "route-table-attachment-example-name"
}

resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name       = var.name
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "foo" {
  vpc_id            = alicloud_vpc.foo.id
  cidr_block        = "172.16.0.0/21"
  zone_id           = data.alicloud_zones.default.zones[0].id
  name              = var.name
}

resource "alicloud_route_table" "foo" {
  vpc_id           = alicloud_vpc.foo.id
  route_table_name = var.name
  description      = "route_table_attachment"
}

resource "alicloud_route_table_attachment" "foo" {
  vswitch_id     = alicloud_vswitch.foo.id
  route_table_id = alicloud_route_table.foo.id
}
```
## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The vswitch_id of the route table attachment, the field can't be changed.
* `route_table_id` - (Required, ForceNew) The route_table_id of the route table attachment, the field can't be changed.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the route table attachment id and formates as `<route_table_id>:<vswitch_id>`.

## Import

The route table attachment can be imported using the id, e.g.

```
$ terraform import alicloud_route_table_attachment.foo vtb-abc123456:vsw-abc123456
```

