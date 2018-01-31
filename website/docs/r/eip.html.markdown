---
layout: "alicloud"
page_title: "Alicloud: alicloud_eip"
sidebar_current: "docs-alicloud-resource-eip"
description: |-
  Provides a ECS EIP resource.
---

# alicloud\_eip

Provides an elastic IP resource.

~> **NOTE:** The resource only support to create `PayByTraffic` elastic IP for international account. Otherwise, you will happened error `COMMODITY.INVALID_COMPONENT`.
Your account is international if you can use it to login in [International Web Console](https://account.alibabacloud.com/login/login.htm).

## Example Usage

```
# Create a new EIP.
resource "alicloud_eip" "example" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
}
```
## Argument Reference

The following arguments are supported:

* `bandwidth` - (Optional) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). If this value is not specified, then automatically sets it to 5 Mbps.
* `internet_charge_type` - (Optional, Forces new resource) Internet charge type of the EIP, Valid values are `PayByBandwidth`, `PayByTraffic`. Default is `PayByBandwidth`. From version `1.7.1`, default to `PayByTraffic`.

## Attributes Reference

The following attributes are exported:

* `id` - The EIP ID.
* `bandwidth` - The elastic public network bandwidth.
* `internet_charge_type` - The EIP internet charge type.
* `status` - The EIP current status.
* `ip_address` - The elastic ip address

## Import

Elastic IP address can be imported using the id, e.g.

```
$ terraform import alicloud_eip.example eip-abc12345678
```