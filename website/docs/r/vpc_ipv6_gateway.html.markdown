---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_gateway"
sidebar_current: "docs-alicloud-resource-vpc-ipv6-gateway"
description: |-
  Provides a Alicloud VPC Ipv6 Gateway resource.
---

# alicloud\_vpc\_ipv6\_gateway

Provides a VPC Ipv6 Gateway resource.

For information about VPC Ipv6 Gateway and how to use it, see [What is Ipv6 Gateway](https://www.alibabacloud.com/help/doc-detail/102214.htm).

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc" "default" {
  vpc_name    = "example_value"
  enable_ipv6 = "true"
}

resource "alicloud_vpc_ipv6_gateway" "example" {
  ipv6_gateway_name = "example_value"
  vpc_id            = alicloud_vpc.default.id
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the IPv6 gateway. The description must be `2` to `256` characters in length. It cannot start with `http://` or `https://`.
* `ipv6_gateway_name` - (Optional) The name of the IPv6 gateway. The name must be `2` to `128` characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`.
* `spec` - (Optional, Computed) The edition of the IPv6 gateway. Valid values: `Large`, `Medium` and `Small`. `Small` (default): Free Edition. `Medium`: Enterprise Edition . `Large`: Enhanced Enterprise Edition. The throughput capacity of an IPv6 gateway varies based on the edition. For more information, see [Editions of IPv6 gateways](https://www.alibabacloud.com/help/doc-detail/98926.htm). 
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC) for which you want to create the IPv6 gateway.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ipv6 Gateway.
* `status` - The status of the resource. Valid values: `Available`, `Pending` and `Deleting`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Ipv6 Gateway.
* `update` - (Defaults to 1 mins) Used when update the Ipv6 Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Ipv6 Gateway.

## Import

VPC Ipv6 Gateway can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_ipv6_gateway.example <id>
```
