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

## Example Usage

```
data "alicloud_ess_scaling_rules" "scalingrules_ds" {
  scaling_group_id = "scaling_group_id"
  ids              = ["scaling_rule_id1", "scaling_rule_id2"]
  name_regex       = "scaling_rule_name"
}

output "first_scaling_rule" {
  value = "${data.alicloud_ess_scaling_rules.scalingrules_ds.rules.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional) Scaling group id the scaling rules belong to.
* `type` - (Optional) Type of scaling rule.
* `name_regex` - (Optional) A regex string to filter resulting scaling rules by name.
* `ids` - (Optional) A list of scaling rule IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

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
