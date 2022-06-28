---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_customer_gateway"
sidebar_current: "docs-alicloud-resource-vpn-customer-gateway"
description: |-
  Provides a Alicloud VPN customer gateway resource.
---

# alicloud\_vpn_customer_gateway

Provides a VPN customer gateway resource.

-> **NOTE:** Terraform will auto build vpn customer gateway instance  while it uses `alicloud_vpn_customer_gateway` to build a vpn customer gateway resource.

For information about VPN customer gateway and how to use it, see [What is VPN customer gateway](https://www.alibabacloud.com/help/en/doc-detail/120368.html).


## Example Usage

Basic Usage

```
resource "alicloud_vpn_customer_gateway" "foo" {
  name        = "vpnCgwNameExample"
  ip_address  = "43.104.22.228"
  description = "vpnCgwDescriptionExample"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the VPN customer gateway. Defaults to null.
* `ip_address` - (Required, ForceNew) The IP address of the customer gateway.
* `description` - (Optional) The description of the VPN customer gateway instance.
* `asn` - (Optional, ForceNew, Available in v1.160.0+) The autonomous system number of the gateway device in the data center. The `asn` is a 4-byte number. You can enter the number in two segments and separate the first 16 bits from the following 16 bits with a period (.). Enter the number in each segment in the decimal format.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPN customer gateway instance id.


#### Timeouts

-> **NOTE:** Available in 1.160.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the vpn customer gateway.
* `update` - (Defaults to 1 mins) Used when update the vpn customer gateway.
* `delete` - (Defaults to 1 mins) Used when delete the vpn customer gateway.


## Import

VPN customer gateway can be imported using the id, e.g.

```
$ terraform import alicloud_vpn_customer_gateway.example cgw-abc123456
```



