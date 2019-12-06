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

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPN customer gateway instance id.

## Import

VPN customer gateway can be imported using the id, e.g.

```
$ terraform import alicloud_vpn_customer_gateway.example cgw-abc123456
```



