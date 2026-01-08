---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_instance_v2"
description: |-
  Provides a Alicloud Cloud Firewall Instance V2 resource.
---

# alicloud_cloud_firewall_instance_v2

Provides a Cloud Firewall Instance V2 resource.

Cloud Firewall instance.

For information about Cloud Firewall Instance V2 and how to use it, see [What is Instance V2](https://www.alibabacloud.com/help/en/product/90174.htm).

-> **NOTE:** Available since v1.269.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_firewall_instance_v2" "default" {
  payment_type = "PayAsYouGo"
  product_code = "cfw"
  product_type = "cfw_elasticity_public_cn"
  spec         = "payg_version"
}
```

## Argument Reference

The following arguments are supported:
* `cfw_log` - (Optional, Bool) Whether to use log audit. Valid values:
  - `true`: Enabled.
  - `false`: Disabled.
* `modify_type` - (Optional) The type of modification. Valid values: `Upgrade`, `Downgrade`. **NOTE:** The `modify_type` is required when you execute an update operation.
* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `PayAsYouGo`, `Subscription`.
* `period` - (Optional) The prepaid period. **NOTE:** If `payment_type` is set to `Subscription`, `period` is required.
* `product_code` - (Required, ForceNew) The product code. Valid values: `cfw`.
* `product_type` - (Required, ForceNew) The product type. Valid values: `cfw_elasticity_public_cn`, `cfw_elasticity_public_intl`, `cfw_sub_public_cn`, `cfw_sub_public_intl`.
* `renewal_duration` - (Optional, Int) The auto-renewal duration. **NOTE:** `renewal_duration` takes effect only if `payment_type` is set to `Subscription`, and `renewal_status` is set to `AutoRenewal`.
* `renewal_duration_unit` - (Optional) The unit of the auto-renewal period. Valid values:
  - `M`: Month.
  - `Y`: Year.
* `renewal_status` - (Optional) Whether to renew an instance automatically or not.
  - `AutoRenewal`: Auto renewal.
  - `ManualRenewal`: Manual renewal.
  - `NotRenewal`: No renewal any longer. After you specify this value, Alibaba Cloud stop sending notification of instance expiry, and only gives a brief reminder on the third day before the instance expiry.
* `sdl` - (Optional, Bool) Data leakage protection status. Valid values: `true`, `false`.
* `spec` - (Required, ForceNew) The edition of the Cloud Firewall instance. Valid values: `payg_version`, `premium_version`, `enterprise_version`,`ultimate_version`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time.
* `end_time` - The end time.
* `release_time` - The release time.
* `status` - The status of Cloud Firewall Instance.
* `user_status` - The user status of Cloud Firewall Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Instance V2.
* `delete` - (Defaults to 5 mins) Used when delete the Instance V2.
* `update` - (Defaults to 5 mins) Used when update the Instance V2.

## Import

Cloud Firewall Instance V2 can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_instance_v2.example <id>
```
