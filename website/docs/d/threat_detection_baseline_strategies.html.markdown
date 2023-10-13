---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_baseline_strategies"
sidebar_current: "docs-alicloud-datasource-threat-detection-baseline-strategies"
description: |-
  Provides a list of Threat Detection Baseline Strategy owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_baseline_strategies

This data source provides Threat Detection Baseline Strategy available to the user.[What is Baseline Strategy](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-describestrategy)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
resource "alicloud_threat_detection_baseline_strategy" "default" {
  custom_type            = "custom"
  end_time               = "08:00:00"
  baseline_strategy_name = "apispec"
  cycle_days             = 3
  target_type            = "groupId"
  start_time             = "05:00:00"
  risk_sub_type_name     = "hc_exploit_redis"
}

data "alicloud_threat_detection_baseline_strategies" "default" {
  ids         = ["${alicloud_threat_detection_baseline_strategy.default.id}"]
  name_regex  = alicloud_threat_detection_baseline_strategy.default.name
  custom_type = "custom"
}

output "alicloud_threat_detection_baseline_strategy_example_id" {
  value = data.alicloud_threat_detection_baseline_strategys.default.strategys.0.id
}
```

## Argument Reference

The following arguments are supported:
* `custom_type` - (ForceNew,Optional) The type of policy. Value:-**common**: standard policy-**custom**: custom policy
* `ids` - (Optional, ForceNew, Computed) A list of Baseline Strategy IDs.
* `baseline_strategy_names` - (Optional, ForceNew) The name of the Baseline Strategy. You can specify at most 10 names.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Baseline Strategy IDs.
* `names` - A list of name of Baseline Strategys.
* `strategys` - A list of Baseline Strategy Entries. Each element contains the following attributes:
    * `id` - The ID of the baseline check policy.
    * `baseline_strategy_id` - The ID of the baseline check policy.
    * `baseline_strategy_name` - Policy name.
    * `custom_type` - The type of policy. Value:
      * **common**: standard policy
      * **custom**: custom policy
    * `cycle_days` - The detection period of the policy.
    * `cycle_start_time` - The detection period of the policy. Value:
      * **0**: 0:00~06:00
      * **6**: 6:00~12:00
      * *12**: 12:00~18:00
      * **18**: 18:00~24:00
    * `end_time` - The baseline check policy execution end time.
    * `start_time` - The baseline check policy start time.
