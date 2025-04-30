---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_baseline_strategy"
sidebar_current: "docs-alicloud-resource-threat-detection-baseline-strategy"
description: |-
  Provides a Alicloud Threat Detection Baseline Strategy resource.
---

# alicloud_threat_detection_baseline_strategy

Provides a Threat Detection Baseline Strategy resource.

For information about Threat Detection Baseline Strategy and how to use it, see [What is Baseline Strategy](https://www.alibabacloud.com/help/en/security-center/latest/api-sas-2018-12-03-modifystrategy).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_baseline_strategy&exampleId=1975341a-56d8-9ff0-627e-9afc19a8a3cd52b1de21&activeTab=example&spm=docs.r.threat_detection_baseline_strategy.0.1975341a56&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_threat_detection_baseline_strategy" "default" {
  custom_type            = "custom"
  end_time               = "08:00:00"
  baseline_strategy_name = "apispec"
  cycle_days             = 3
  target_type            = "groupId"
  start_time             = "05:00:00"
  risk_sub_type_name     = "hc_exploit_redis"
}
```

## Argument Reference

The following arguments are supported:
* `baseline_strategy_name` - (Required, ForceNew) Policy name.
* `custom_type` - (Required) The type of policy. Value:
  * **common**: standard policy
  * **custom**: custom policy
* `cycle_days` - (Required) The detection period of the policy.
* `cycle_start_time` - (Optional) The detection period of the policy. Value:
  * **0**: 0:00~06:00
  * **6**: 6:00~12:00
  * **12**: 12:00~18:00
  * **18**: 18:00~24:00
* `end_time` - (Required) The baseline check policy execution end time.
* `risk_sub_type_name` - (Required) Detection item subtype.
* `start_time` - (Required) The baseline check policy start time.
* `target_type` - (Required) The method of adding assets that take effect from the policy. Value:
  * **groupId**: Added by asset group.
  * **uuid**: Add by single asset.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the baseline check policy, same with `baseline_strategy_id`.
* `baseline_strategy_id` - The ID of the baseline check policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Baseline Strategy.
* `delete` - (Defaults to 5 mins) Used when delete the Baseline Strategy.
* `update` - (Defaults to 5 mins) Used when update the Baseline Strategy.

## Import

Threat Detection Baseline Strategy can be imported using the id, e.g.

```shell
$terraform import alicloud_threat_detection_baseline_strategy.example <id>
```