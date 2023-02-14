---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_instance"
sidebar_current: "docs-alicloud-resource-threat_detection-instance"
description: |-
  Provides a Alicloud Threat Detection Instance resource.
---

# alicloud_threat_detection_instance

Provides a Threat Detection Instance resource.

For information about Threat Detection Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/security-center/latest/what-is-security-center).

-> **NOTE:** Available in v1.199.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_instance" "default" {
  payment_type           = "Subscription"
  period                 = 12
  renewal_status         = "ManualRenewal"
  sas_sls_storage        = "100"
  sas_anti_ransomware    = "100"
  container_image_scan   = "100"
  sas_webguard_order_num = "100"
  sas_sc                 = "true"
  version_code           = "level2"
  buy_number             = "30"
  honeypot_switch        = "1"
  sas_sdk_switch         = "1"
  sas_sdk                = "1000"
  honeypot               = "32"
  v_core                 = "100"
}
```

## Argument Reference

The following arguments are supported:
* `modify_type` - (Optional) Change configuration type, value
  - Upgrade: Upgrade.
  - Downgrade: Downgrade.
* `buy_number` - (Optional) Number of servers.
* `container_image_scan` - (Optional) Container Image security scan.
* `honeypot` - (Optional) Cloud honeypot authorization number.
* `honeypot_switch` - (Optional) Cloud honeypot. Valid values: `1`, `2`.
* `payment_type` - (Required) The payment type of the resource.
* `period` - (Optional) Prepaid cycle. The unit is Monthly, please enter an integer multiple of 12 for annual paid products. **NOTE:** must be set when creating a prepaid instance.
* `renew_period` - (Optional) Automatic renewal cycle, in months. **NOTE:** The `renew_period` is required under the condition that `renewal_status` is `AutoRenewal`.
* `renewal_status` - (Optional,Computed) Automatic renewal status, Default ManualRenewal. value:
  - `AutoRenewal`: automatic renewal.
  - `ManualRenewal`: manual renewal.
* `renewal_period_unit` - (Optional,Computed) The unit of the auto-renewal period. **NOTE:** The `renewal_period_unit` is required under the condition that `renewal_status` is `AutoRenewal`. Valid values: 
  - `M`: months.
  - `Y`: years.
* `sas_anti_ransomware` - (Optional) Anti-extortion.
* `sas_sc` - (Optional) Large security screen.
* `sas_sdk` - (Optional) Number of malicious file detections.
* `sas_sdk_switch` - (Optional) Malicious file detection SDK. Valid values: `0`, `1`.
* `sas_sls_storage` - (Optional) Log analysis.
* `sas_webguard_boolean` - (Optional) Web page tamper-proof.  Valid values: `0`, `1`.
* `sas_webguard_order_num` - (Optional) Number of tamper-proof authorizations.
* `threat_analysis` - (Optional) The amount of threat analysis log storage.
* `threat_analysis_switch` - (Optional) Threat analysis.  Valid values: `0`, `1`.
* `v_core` - (Optional) Number of cores.
* `version_code` - (Required) Version selection. Valid values: `level10`, `level2`, `level3`, `level7`, `level8`.



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - The creation time of the resource
* `instance_id` - The first ID of the resource
* `status` - The status of the resource

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Threat Detection Instance do not support import.