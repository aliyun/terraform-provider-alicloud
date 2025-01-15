---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_rules"
sidebar_current: "docs-alicloud_ess_scaling_rules"
description: |-
    Provides a list of scaling rules available to the user.
---

# alicloud_ess_scaling_rules

This data source provides available scaling rule resources. 

-> **NOTE:** Available since v1.39.0

## Example Usage

```terraform
variable "name" {
  default = "terraform-ex"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

locals {
  name = "${var.name}-${random_integer.default.result}"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = local.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = local.name
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = local.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
}

resource "alicloud_ess_scaling_rule" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  scaling_rule_name = local.name
  adjustment_type   = "PercentChangeInCapacity"
  adjustment_value  = 1
}


data "alicloud_ess_scaling_rules" "scalingrules_ds" {
  scaling_group_id = alicloud_ess_scaling_group.default.id
  ids              = [alicloud_ess_scaling_rule.default.id]
  name_regex       = local.name
}

output "first_scaling_rule" {
  value = data.alicloud_ess_scaling_rules.scalingrules_ds.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional, ForceNew) Scaling group id the scaling rules belong to.
* `type` - (Optional, ForceNew) Type of scaling rule.
* `name_regex` - (Optional, ForceNew) A regex string to filter resulting scaling rules by name.
* `ids` - (Optional, ForceNew) A list of scaling rule IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of scaling rule ids.
* `names` - A list of scaling rule names.
* `rules` - A list of scaling rules. Each element contains the following attributes:
  * `id` - ID of the scaling rule.
  * `scaling_group_id` - ID of the scaling group.
  * `name` - Name of the scaling rule.
  * `type` - Type of the scaling rule.
  * `cooldown` - Cooldown time of the scaling rule.
  * `adjustment_type` - Adjustment type of the scaling rule.
  * `adjustment_value` - Adjustment value of the scaling rule.
  * `min_adjustment_magnitude` - Min adjustment magnitude of scaling rule.
  * `scaling_rule_ari` - Ari of scaling rule.
  * `initial_max_size` - (Available since v1.242.0) The maximum number of ECS instances that can be added to the scaling group.
  * `predictive_value_behavior` - (Available since v1.242.0) The action on the predicted maximum value.
  * `predictive_scaling_mode` - (Available since v1.242.0) The mode of the predictive scaling rule.
  * `predictive_value_buffer` - (Available since v1.242.0) The ratio based on which the predicted value is increased if you set predictive_value_behavior to PredictiveValueOverrideMaxWithBuffer. If the predicted value that is increased by this ratio is greater than the initial maximum capacity, the increased value is used as the maximum value for prediction tasks.
  * `predictive_task_buffer_time` - (Available since v1.242.0) The amount of buffer time before the prediction task is executed. By default, all prediction tasks that are automatically created by a predictive scaling rule are executed on the hour. You can set a buffer time to execute prediction tasks and prepare resources in advance.
  * `target_value` - (Available since v1.242.0) The target value of the metric.
  * `metric_name` - (Available since v1.242.0) The predefined metric of the scaling rule. 
