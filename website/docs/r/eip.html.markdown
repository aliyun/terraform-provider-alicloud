---
layout: "alicloud"
page_title: "Alicloud: alicloud_eip"
sidebar_current: "docs-alicloud-resource-eip"
description: |-
  Provides a ECS EIP resource.
---

# alicloud\_eip

Provides an elastic IP resource.

-> **NOTE:** The resource only support to create `PayByTraffic` elastic IP for international account. Otherwise, you will happened error `COMMODITY.INVALID_COMPONENT`.
Your account is international if you can use it to login in [International Web Console](https://account.alibabacloud.com/login/login.htm).

-> **NOTE:** From version 1.10.1, this resource supports creating "PrePaid" EIP.

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
* `internet_charge_type` - (Optional, ForceNew) Internet charge type of the EIP, Valid values are `PayByBandwidth`, `PayByTraffic`. Default to `PayByBandwidth`. From version `1.7.1`, default to `PayByTraffic`.
* `instance_charge_type` - (Optional, ForceNew) Elastic IP instance charge type. Valid values are "PrePaid" and "PostPaid". Default to "PostPaid".
* `period` - (Optional, ForceNew) The duration that you will buy the resource, in month. It is valid when `instance_charge_type` is `PrePaid`.
Default to 1. Valid values: [1-9, 12, 24, 36]. At present, the provider does not support modify "period" and you can do that via web console.

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