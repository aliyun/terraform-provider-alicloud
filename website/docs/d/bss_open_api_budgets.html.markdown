---
subcategory: "Bss Open Api"
layout: "alicloud"
page_title: "Alicloud: alicloud_bss_open_api_budgets"
sidebar_current: "docs-alicloud-datasource-bss-open-api-budgets"
description: |-
  Provides a list of Bss Open Api Budgets owned by an Alibaba Cloud account.
---

# alicloud\_bss_open_api_budgets

This data source provides Bss Open Api Budgets available to the user. [What is Budget](https://next.api.alibabacloud.com/document/BssOpenApi/2023-09-30/DescribeBudgets)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
data "alicloud_bss_open_api_budgets" "default" {
  ids        = ["my-budget"]
  name_regex = "terraform-.*"
}

output "first_budget_name" {
  value = data.alicloud_bss_open_api_budgets.default.budgets.0.budget_name
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional) A list of Budget IDs (budget names).
* `name_regex` - (Optional) A regex string to filter results by budget name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Budget IDs (budget names).
* `names` - A list of budget names.
* `budgets` - A list of Budget Entries. Each element contains the following attributes:
  * `budget_name` - The name of the budget.
  * `budget_type` - The type of the budget.
  * `cycle_end_period` - The end period of the budget cycle.
  * `cycle_start_period` - The start period of the budget cycle.
  * `cycle_type` - The cycle type of the budget.
  * `metric` - The metric of the budget.
  * `quota` - The fixed quota value. When the quota type is a quota, the unit is a percentage.
  * `quota_type` - The quota type of the budget.
  * `comment` - The comment of the budget.
  * `cycle_quota` - The specified quota for each cycle.
    * `cycle_period` - The cycle period.
    * `quota` - The quota.
  * `query_filter` - The filter conditions.
    * `code` - The parameter code.
    * `select_type` - The select mode.
    * `values` - The list of filter values.
  * `warn_confs` - The alert configurations.
    * `comment` - The comment.
    * `event_bridge` - Whether to enable EventBridge.
    * `msc_channels` - The list of message center channels.
    * `msc_contacts` - The list of message center contacts.
    * `name` - The alert name.
    * `threshold_type` - The threshold type.
    * `threshold_value` - The threshold value.
    * `warn_target` - The alert target.
