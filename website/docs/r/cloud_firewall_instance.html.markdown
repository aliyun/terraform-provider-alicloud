---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_instance"
sidebar_current: "docs-alicloud-resource-cloud-firewall-instance"
description: |-
  Provides a Alicloud Cloud Firewall Instance resource.
---

# alicloud\_cloud\_firewall\_instance

Provides a Cloud Firewall Instance resource.

For information about Cloud Firewall Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/product/90174.htm).

-> **NOTE:** Available in v1.139.0+.


## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_firewall_instance" "example" {
  payment_type    = "Subscription"
  spec            = "premium_version"
  ip_number       = 20
  band_width      = 10
  cfw_log         = false
  cfw_log_storage = 1000
  cfw_service     = false
  period          = 6
}
```

## Argument Reference

The following arguments are supported:

* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `Subscription`.
* `period` - (Required) The prepaid period. Valid values: `6`, `12`, `24`, `36`.
* `renew_period` - (Optional) Automatic renewal period. **NOTE:** The `renew_period` is required under the condition that renewal_status is `AutoRenewal`.
* `renewal_status` - (Optional) Automatic renewal status. Valid values: `AutoRenewal`,`ManualRenewal`. Default Value: `ManualRenewal`.
* `logistics` - (Optional) The logistics.
* `modify_type` - (Optional) The modify type. Valid values: `Upgrade`, `Downgrade`.  **NOTE:** The `modify_type` is required when you execute an update operation.
* `cfw_service` - (Required) Whether to use expert service. Valid values: `true`, `false`.
* `spec` - (Required) Current version. Valid values: `premium_version`, `enterprise_version`,`ultimate_version`.
* `cfw_log_storage` - (Required) The log storage capacity. 
  * `premium_version` - The valid cfw_log_storage is [1000, 500000] with the step size 1000. Default Value: `1000`. Unit: GB.
  * `enterprise_version` - The valid cfw_log_storage is [3000, 500000] with the step size 1000. Default Value: `3000`. Unit: GB.
  * `ultimate_version` - The valid cfw_log_storage is [5000, 500000] with the step size 1000. Default Value: `5000`. Unit: GB.
* `ip_number` - (Required) The number of public IPs that can be protected. Valid values: 20 to 4000.
  * `premium_version` - The valid cfw_log_storage is [60, 1000] with the step size 1. Default Value: `20`. 
  * `enterprise_version` - The valid cfw_log_storage is [60, 1000] with the step size 1. Default Value: `50`. 
  * `ultimate_version` - The valid cfw_log_storage is [400, 4000] with the step size 1. Default Value: `400`. 
* `cfw_log` - (Required) Whether to use log audit. Valid values: `true`, `false`.
* `band_width` - (Required) Public network processing capability. Valid values: 10 to 15000. Unit: Mbps.
  * `premium_version` - The valid cfw_log_storage is [10, 2000] with the step size 5. Default Value: `10`. Unit: Mbps.
  * `enterprise_version` - The valid cfw_log_storage is [50, 5000] with the step size 5. Default Value: `50`. Unit: Mbps.
  * `ultimate_version` - The valid cfw_log_storage is [200, 15000] with the step size 5. Default Value: `200`. Unit: Mbps.
* `fw_vpc_number` - (Optional) The number of protected VPCs. Valid values between 2 and 500.
  * `enterprise_version` - The valid cfw_log_storage is [2, 200] with the step size 1. Default Value: `2`. 
  * `ultimate_version` - The valid cfw_log_storage is [5, 500] with the step size 1. Default Value: `5`. 
* `instance_count` - (Optional)  The number of assets.


## Attributes Reference

The following attributes are exported:

* `create_time` - The creation time.
* `renewal_duration_unit` - Automatic renewal period unit. Valid values: `Month`,`Year`.
* `status` - The status of Instance.
* `end_time` - The end time.
* `release_time` - The release time.

## Import

Cloud Firewall Instance can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_firewall_instance.example <id>
```
