---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eip_address"
sidebar_current: "docs-alicloud-resource-eip-address"
description: |-
  Provides a Alicloud EIP Address resource.
---

# alicloud\_eip\_address

Provides a EIP Address resource.

For information about EIP Address and how to use it, see [What is EIP Address](https://www.alibabacloud.com/help/en/doc-detail/36016.htm).

-> **NOTE:** Available in v1.126.0+.

-> **NOTE:** BGP (Multi-ISP) lines are supported in all regions. BGP (Multi-ISP) Pro lines are supported only in the China (Hong Kong) region.

-> **NOTE:** The resource only supports to create `PayAsYouGo PayByTraffic`  or `Subscription PayByBandwidth` elastic IP for international account. Otherwise, you will happened error `COMMODITY.INVALID_COMPONENT`.
Your account is international if you can use it to login in [International Web Console](https://account.alibabacloud.com/login/login.htm).

## Example Usage

Basic Usage

```terraform
resource "alicloud_eip_address" "example" {
  address_name         = "tf-testAcc1234"
  isp                  = "BGP"
  internet_charge_type = "PayByBandwidth"
  payment_type         = "PayAsYouGo"
}

```

## Argument Reference

The following arguments are supported:

* `activity_id` - (Optional) The activity id.
* `address_name` - (Optional) The name of the EIP instance. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://.
* `bandwidth` - (Optional, Computed) The maximum bandwidth of the EIP. Valid values: `1` to `200`. Unit: Mbit/s. Default value: `5`.
* `deletion_protection` - (Optional, Computed) Whether enable the deletion protection or not. Default value: `false`.
* `description` - (Optional) The description of the EIP.
* `internet_charge_type` - (Optional, ForceNew) The metering method of the EIP. 
  Valid values: `PayByDominantTraffic`, `PayByBandwidth` and `PayByTraffic`. Default to `PayByBandwidth`. **NOTE:** It must be set to "PayByBandwidth" when `payment_type` is "Subscription".
* `isp` - (Optional, ForceNew) The line type. You can set this parameter only when you create a `PayAsYouGo` EIP. Valid values: `BGP`: BGP (Multi-ISP) lines.Up to 89 high-quality BGP lines are available worldwide. Direct connections with multiple Internet Service Providers (ISPs), including Telecom, Unicom, Mobile, Railcom, Netcom, CERNET, China Broadcast Network, Dr. Peng, and Founder, can be established in all regions in mainland China. `BGP_PRO`:  BGP (Multi-ISP) Pro lines optimize data transmission to mainland China and improve connection quality for international services. Compared with BGP (Multi-ISP), when BGP (Multi-ISP) Pro provides services to clients in mainland China (excluding data centers), cross-border connections are established without using international ISP services. This reduces network latency.
* `netmode` - (Optional) The type of the network. Valid value is `public` (Internet).
* `payment_type` - (Optional, ForceNew) The billing method of the EIP. Valid values: `Subscription` and `PayAsYouGo`. Default value is `PayAsYouGo`. 
* `period` - (Optional) The duration that you will buy the resource, in month. It is valid when `payment_type` is `Subscription`. Valid values: [1-9, 12, 24, 36]. At present, the provider does not support modify "period" and you can do that via web console.
* `resource_group_id` - (Optional) The ID of the resource group.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `auto_pay`  - (Optional, Available in v1.140.0+) Whether to pay automatically. Valid values: `true` and `false`. Default value: `true`. When `auto_pay` is `true`, The order will be automatically paid. When `auto_pay` is `false`, The order needs to go to the order center to complete the payment. **NOTE:** When `payment_type` is `Subscription`, this parameter is valid.
* `instance_charge_type` - (Optional, ForceNew, Deprecated in v1.126.0+) Field `instance_charge_type` has been deprecated from provider version 1.126.0, and it will be removed in the future version. Please use the new attribute `payment_type` instead.
* `name` - (Optional, Computed, Deprecated in v1.126.0+) Field `name` has been deprecated from provider version 1.126.0, and it will be removed in the future version. Please use the new attribute `address_name` instead.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Address.
* `status` - The status of the EIP. Valid values:  `Associating`: The EIP is being associated. `Unassociating`: The EIP is being disassociated. `InUse`: The EIP is allocated. `Available`:The EIP is available.
* `ip_address` - The address of the EIP.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Address.
* `delete` - (Defaults to 9 mins) Used when delete the Address.

## Import

EIP Address can be imported using the id, e.g.

```
$ terraform import alicloud_eip_address.example <id>
```
