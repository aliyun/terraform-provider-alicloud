---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv4_gateway"
sidebar_current: "docs-alicloud-resource-vpc-ipv4-gateway"
description: |-
  Provides a Alicloud VPC Ipv4 Gateway resource.
---

# alicloud\_vpc\_ipv4\_gateway

Provides a VPC Ipv4 Gateway resource.

For information about VPC Ipv4 Gateway and how to use it, see [What is Ipv4 Gateway](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/createipv4gateway).

-> **NOTE:** Available in v1.181.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpcs" "default" {
  name_regex = "default-NoDeleting"
}
resource "alicloud_vpc_ipv4_gateway" "example" {
  ipv4_gateway_name = "example_value"
  vpc_id            = data.alicloud_vpcs.default.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `ipv4_gateway_description` - (Optional) The description of the IPv4 gateway. The description must be `2` to `256` characters in length. It must start with a letter but cannot start with `http://` or `https://`.
* `ipv4_gateway_name` - (Optional) The name of the IPv4 gateway. The name must be `2` to `128` characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). It must start with a letter.
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC) where you want to create the IPv4 gateway. You can create only one IPv4 gateway in a VPC.
* `enabled` - (Optional, Available in v1.193.1+) Whether the IPv4 gateway is active or not. Valid values are `true` and `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ipv4 Gateway.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Ipv4 Gateway.
* `update` - (Defaults to 1 mins) Used when updating the Ipv4 Gateway.
* `delete` - (Defaults to 1 mins) Used when deleting the Ipv4 Gateway.


## Import

VPC Ipv4 Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipv4_gateway.example <id>
```