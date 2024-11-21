---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_instance"
sidebar_current: "docs-alicloud-resource-cloud-firewall-instance"
description: |-
  Provides an Alicloud Cloud Firewall Instance resource.
---

# alicloud_cloud_firewall_instance

Provides a Cloud Firewall Instance resource.

For information about Cloud Firewall Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/product/90174.htm).

-> **NOTE:** Available since v1.139.0.


## Example Usage

Basic Usage

create a pay-as-you-go instance

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_instance&exampleId=6adc2c91-ffcd-d93f-6b09-9ee9035d294d9c9e777d&activeTab=example&spm=docs.r.cloud_firewall_instance.0.6adc2c91ff&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_cloud_firewall_instance" "PayAsYouGo" {
  payment_type = "PayAsYouGo"
}
```

create a subscription instance
<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_instance&exampleId=8bf5e4ee-d6d4-0b01-fbb3-3164539c417910842f3a&activeTab=example&spm=docs.r.cloud_firewall_instance.1.8bf5e4eed6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_cloud_firewall_instance" "Subscription" {
  payment_type = "Subscription"
  spec         = "premium_version"
  ip_number    = 20
  band_width   = 10
  cfw_log      = false
  period       = 1
}
```

## Argument Reference

The following arguments are supported:

* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `Subscription`, `PayAsYouGo`. **NOTE:** From version 1.220.0, `payment_type` can be set to `PayAsYouGo`.
* `period` - (Optional) The prepaid period. Valid values: `1`, `3`, `6`, `12`, `24`, `36`. **NOTE:** 1 and 3 available since 1.204.1. If `payment_type` is set to `Subscription`, `period` is required. Otherwise, it will be ignored.
* `renew_period` - (Deprecated since v1.209.1) Automatic renewal period. Attribute `renew_period` has been deprecated since 1.209.1. Using `renewal_duration` instead.
* `renewal_duration` - (Optional) Auto-Renewal Duration. It is required under the condition that `renewal_status` is `AutoRenewal`. Valid values: `1`, `2`, `3`, `6`, `12`.
**NOTE:** `renewal_duration` takes effect only if `payment_type` is set to `Subscription`, and `renewal_status` is set to `AutoRenewal`.
* `renewal_duration_unit` - (Optional) Auto-Renewal Cycle Unit Values Include: Month: Month. Year: Years. Valid values: `Month`, `Year`.
* `renewal_status` - (Optional) Whether to renew an instance automatically or not. Default value: `ManualRenewal`.
  - `AutoRenewal`: Auto renewal.
  - `ManualRenewal`: Manual renewal.
  - `NotRenewal`: No renewal any longer. After you specify this value, Alibaba Cloud stop sending notification of instance expiry, and only gives a brief reminder on the third day before the instance expiry.
**NOTE:** `renewal_status` takes effect only if `payment_type` is set to `Subscription`.
* `logistics` - (Optional) The logistics.
* `modify_type` - (Optional) The type of modification. Valid values: `Upgrade`, `Downgrade`. **NOTE:** The `modify_type` is required when you execute an update operation.
* `cfw_service` - (Removed since v1.209.1) Attribute `cfw_service` does not support longer, and it has been removed since v1.209.1.
* `spec` - (Optional) Current version. Valid values: `premium_version`, `enterprise_version`,`ultimate_version`.
* `cfw_log` - (Optional) Whether to use log audit. Valid values: `true`, `false`. **NOTE:** From version 1.232.0, When `payment_type` is set to `PayAsYouGo`, `cfw_log` can only be set to `true`, `cfw_log` cannot be modified to `false`.
* `cfw_log_storage` - (Optional) The log storage capacity. **NOTE:** From version 1.232.0, When `payment_type` is set to `PayAsYouGo`, or `cfw_log` is set to `false`, `cfw_log_storage` will be ignored.
  * `premium_version` - The valid cfw_log_storage is [1000, 500000] with the step size 1000. Default Value: `1000`. Unit: GB.
  * `enterprise_version` - The valid cfw_log_storage is [3000, 500000] with the step size 1000. Default Value: `3000`. Unit: GB.
  * `ultimate_version` - The valid cfw_log_storage is [5000, 500000] with the step size 1000. Default Value: `5000`. Unit: GB.
* `ip_number` - (Optional) The number of public IPs that can be protected. Valid values: 20 to 4000.
  * `premium_version` - The valid cfw_log_storage is [60, 1000] with the step size 1. Default Value: `20`.
  * `enterprise_version` - The valid cfw_log_storage is [60, 1000] with the step size 1. Default Value: `50`.
  * `ultimate_version` - The valid cfw_log_storage is [400, 4000] with the step size 1. Default Value: `400`.
* `band_width` - (Optional) Public network processing capability. Valid values: 10 to 15000. Unit: Mbps.
  * `premium_version` - The valid cfw_log_storage is [10, 2000] with the step size 5. Default Value: `10`. Unit: Mbps.
  * `enterprise_version` - The valid cfw_log_storage is [50, 5000] with the step size 5. Default Value: `50`. Unit: Mbps.
  * `ultimate_version` - The valid cfw_log_storage is [200, 15000] with the step size 5. Default Value: `200`. Unit: Mbps.
* `fw_vpc_number` - (Optional) The number of protected VPCs. It will be ignored when `spec = "premium_version"`. Valid values between 2 and 500.
  * `enterprise_version` - The valid cfw_log_storage is [2, 200] with the step size 1. Default Value: `2`.
  * `ultimate_version` - The valid cfw_log_storage is [5, 500] with the step size 1. Default Value: `5`.
* `instance_count` - (Optional)  The number of assets.
* `cfw_account` - (Optional, Available since v1.209.1, Bool) Whether to use multi-account. Valid values: `true`, `false`.
* `account_number` - (Optional, Available since v1.209.1, Int) The number of multi account. It will be ignored when `cfw_account = false`.
  * `premium_version` - The valid account number is [1, 20].
  * `enterprise_version` - The valid account number is [1, 50].
  * `ultimate_version` - The valid account number is [1, 1000].

## Attributes Reference

The following attributes are exported:

* `user_status` - (Available since v1.232.0) The user status of Cloud Firewall Instance.
* `status` - The status of Cloud Firewall Instance.
* `create_time` - The creation time.
* `end_time` - The end time.
* `release_time` - The release time.

## Import

Cloud Firewall Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_instance.example <id>
```
