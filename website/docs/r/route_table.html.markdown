---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_route_table"
description: |-
  Provides a Alicloud VPC Route Table resource.
---

# alicloud_route_table

Provides a VPC Route Table resource.

Currently, customized route tables are available in most regions apart from China (Beijing), China (Hangzhou), and China (Shenzhen) regions.

For information about VPC Route Table and how to use it, see [What is Route Table](https://www.alibabacloud.com/help/doc-detail/87057.htm).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_vpc" "defaultVpc" {
  vpc_name = var.name
}


resource "alicloud_route_table" "default" {
  description      = "test-description"
  vpc_id           = alicloud_vpc.defaultVpc.id
  route_table_name = var.name
  associate_type   = "VSwitch"
}
```

## Argument Reference

The following arguments are supported:
* `associate_type` - (Optional, ForceNew, Computed) The type of cloud resource that is bound to the routing table. Value:
  - `VSwitch`: switch.
  - `Gateway`:IPv4 Gateway.
* `description` - (Optional) Description of the routing table.
* `route_propagation_enable` - (Optional, Available since v1.245.0) Route Table Receive Propagate Route State
* `route_table_name` - (Optional) The name of the routing table.
* `tags` - (Optional, Map) The tag
* `vpc_id` - (Required, ForceNew) The ID of VPC.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.119.1). Field 'name' has been deprecated from provider version 1.119.1. New field 'route_table_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the routing table
* `resource_group_id` - Resource group ID.
* `status` - Routing table state

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Route Table.
* `delete` - (Defaults to 5 mins) Used when delete the Route Table.
* `update` - (Defaults to 5 mins) Used when update the Route Table.

## Import

VPC Route Table can be imported using the id, e.g.

```shell
$ terraform import alicloud_route_table.example <id>
```