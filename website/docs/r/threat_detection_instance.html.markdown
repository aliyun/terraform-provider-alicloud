---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_instance"
description: |-
  Provides a Alicloud Threat Detection Instance resource.
---

# alicloud_threat_detection_instance

Provides a Threat Detection Instance resource.

Cloud Security Center instance.

For information about Threat Detection Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/security-center/latest/what-is-security-center).

-> **NOTE:** Available since v1.199.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_instance&exampleId=b7939698-097c-9f53-a286-e40654e839a08cfdfa68&activeTab=example&spm=docs.r.threat_detection_instance.0.b793969809&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_threat_detection_instance" "default" {
  threat_analysis          = "10"
  sas_sls_storage          = "10"
  v_core                   = "10"
  sas_sc                   = "false"
  buy_number               = "10"
  honeypot_switch          = "2"
  payment_type             = "Subscription"
  sas_sdk                  = "10"
  sas_anti_ransomware      = "10"
  renewal_status           = "ManualRenewal"
  period                   = "1"
  vul_switch               = "1"
  rasp_count               = "1"
  vul_count                = "20"
  version_code             = "level3"
  sas_cspm                 = "1000"
  renewal_period_unit      = "M"
  container_image_scan_new = "100"
  honeypot                 = "20"
}
```

### Deleting `alicloud_threat_detection_instance` or removing it from your configuration

Terraform cannot destroy resource `alicloud_threat_detection_instance`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `buy_number` - (Optional) Number of servers.
* `container_image_scan` - (Optional, Deprecated since v1.247.0) Container Image security scan. Interval type, value interval:[0,200000].

-> **NOTE:**  The step size is 20, that is, only multiples of 20 can be filled in.

* `container_image_scan_new` - (Optional, Available since v1.212.0) Container Image security scan. Interval type, value interval:[0,200000].

-> **NOTE:**  The step size is 20, that is, only multiples of 20 can be filled in.

* `honeypot` - (Optional) Number of cloud honeypot licenses. Interval type, value interval:[20,500].

-> **NOTE:**  This module can only be purchased when honeypot_switch = 1, starting with 20.

* `honeypot_switch` - (Optional, Computed) Cloud honeypot. Value:
  - 1: Yes.
  - 2: No.
* `modify_type` - (Optional) Change configuration type, value
  - Upgrade: Upgrade.
  - Downgrade: Downgrade.
* `payment_type` - (Required, ForceNew) The payment type of the resource
* `period` - (Optional, Int) Prepaid cycle. The unit is Monthly, please enter an integer multiple of 12 for annual paid products.

-> **NOTE:**  must be set when creating a prepaid instance.

* `post_paid_flag` - (Optional, Int, Available since v1.247.0) Post-paid signage. Value:
  - `1`: Postpaid instance
* `post_pay_module_switch` - (Optional, JsonString, Available since v1.247.0) Pay-as-you-go module switch mapping, in JsonString format. Valid values:
  - Key:
  - `VUL`: vulnerability repair module
  - `CSPM`: Cloud platform configuration check module
  - `AGENTLESS`: AGENTLESS detection module
  - `SERVERLESS`:Serverless asset module
  - `CTDR`: threat analysis and response module
  - Value:0 means off, 1 means on

-> **NOTE:**  The module value of the unpassed value will not change.

* `rasp_count` - (Optional, Available since v1.212.0) Number of application protection licenses. Interval type, value interval:[1,100000000].
* `renew_period` - (Optional, Computed, Int) Automatic renewal cycle, in months.

-> **NOTE:**  When `RenewalStatus` is set to `AutoRenewal`, it must be set.

* `renewal_period_unit` - (Optional, Computed) Automatic renewal period unit, value:
  - M: month.
  - Y: years.

-> **NOTE:**  Must be set when RenewalStatus = AutoRenewal.

* `renewal_status` - (Optional) Automatic renewal status, value:
  - AutoRenewal: automatic renewal.
  - ManualRenewal: manual renewal.

Default ManualRenewal.
* `sas_anti_ransomware` - (Optional) Anti-ransomware capacity. Unit: GB. Interval type, value interval:[0,9999999999].

-> **NOTE:**  The step size is 10, that is, only multiples of 10 can be filled in.

* `sas_cspm` - (Optional, Available since v1.212.0) Cloud platform configuration check scan times, interval type, value range:[15000,9999999999].

-> **NOTE:**  You must have sas_cspm_switch = 1 to purchase this module. The step size is 55000, that is, only multiples of 55000 can be filled in.

* `sas_cspm_switch` - (Optional, Computed, Available since v1.212.0) Cloud platform configuration check switch. Value:
  - 0: No.
  - 1: Yes.
* `sas_sc` - (Optional) Security screen. Value:
  - true: Yes.
  - false: No.
* `sas_sdk` - (Optional) Number of malicious file detections. Unit: 10,000 times. Interval type, value interval:[10,9999999999].

-> **NOTE:**  This module can only be purchased when sas_sdk_switch = 1. The step size is 10, that is, only multiples of 10 can be filled in.

* `sas_sdk_switch` - (Optional, Computed) Malicious file detection SDK.
* `sas_sls_storage` - (Optional) Log analysis storage capacity. Unit: GB. Interval type, value interval:[0,600000].

-> **NOTE:**  The step size is 10, that is, only multiples of 10 can be filled in.

* `sas_webguard_boolean` - (Optional, Computed) Web tamper-proof switch. Value:
  - 0: No.
  - 1: Yes.
* `sas_webguard_order_num` - (Optional) Tamper-proof authorization number. Value:
  - 0: No
1: Yes.
* `subscription_type` - (Optional, Available since v1.247.0) The subscription type. Value:
  - Subscription: Prepaid.
  - PayAsYouGo: Post-paid.
* `threat_analysis` - (Optional) Threat Analysis log storage capacity. Interval type, value interval:[0,9999999999].

-> **NOTE:**  This module can only be purchased when Threat_analysis_switch = 1. The step size is 10, that is, only multiples of 10 can be filled in.

* `threat_analysis_flow` - (Optional, Available since v1.232.0) Threat analysis and response log access traffic. After ThreatAnalysisSwitch1 is selected, it must be selected. Interval type, value interval:[0,9999999999].

-> **NOTE:**  Step size is 1.

* `threat_analysis_sls_storage` - (Optional, Available since v1.232.0) Threat analysis and response log storage capacity. Interval type, value interval:[100,9999999999].

-> **NOTE:**  The step size is 100, that is, only multiples of 100 can be filled in.

* `threat_analysis_switch` - (Optional, Computed) Threat analysis. Value:
  - 0: No.
  - 1: Yes.
* `threat_analysis_switch1` - (Optional, Computed, Available since v1.232.0) Threat analysis and response switch (new), values:
  - 0: No
  - 1: Yes
* `v_core` - (Optional) Number of cores.
* `version_code` - (Optional) Select the security center version. Value:
  - level7: Antivirus Edition.
  - level3: Premium version.
  - level2: Enterprise Edition.
  - level8: Ultimate.
  - level10: Purchase value-added services only.
* `vul_count` - (Optional, Available since v1.212.0) Vulnerability repair times, interval type, value range:[20,100000000].

-> **NOTE:**  This module can only be purchased when vul_switch = 1. Only when the version_code value is level7 or level10. other versions do not need to be purchased separately.

* `vul_switch` - (Optional, Computed) Vulnerability fix switch. Value:
  - 0: No.
  - 1: Yes.

-> **NOTE:**  When the value of version_code is level7 or level10, the purchase is allowed. Other versions do not need to be purchased separately.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `status` - The resource attribute field representing the resource status. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Threat Detection Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_instance.example <id>
```