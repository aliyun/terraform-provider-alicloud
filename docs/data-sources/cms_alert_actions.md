---
subcategory: "Cloud Monitor (CMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_actions"
description: |-
  Provides a list of Cms Alert Actions to the user.
---

# alicloud_cms_alert_actions

This data source provides the list of Cms Alert Actions in an Alibaba Cloud account.

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_alert_actions" "example" {
  type = "WEBHOOK"
  output_file = "alert_actions.txt"
}

output "first_action_id" {
  value = data.alicloud_cms_alert_actions.example.alert_actions.0.alert_action_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of alert action IDs to filter by.
* `name_regex` - (Optional) A regex string to filter alert actions by name.
* `type` - (Optional) The type of alert actions to filter by. Valid values: `FC`, `MNS`, `OPEN_API`, `SLS`, `ESS`, `PAGER_DUTY`, `WEBHOOK`, `EB`, `FC3`.
* `output_file` - (Optional) File name where to save the result after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of alert action IDs.
* `alert_actions` - A list of Cms Alert Actions. Each element contains the following attributes:
  * `alert_action_id` - The ID of the alert action.
  * `alert_action_name` - The name of the alert action.
  * `type` - The type of the alert action.
  * `region_id` - The region ID of the alert action.
  * `webhook_param` - The webhook parameters.
  * `mns_param` - The message queue parameters.
  * `sls_param` - The log service parameters.
  * `ess_param` - The elastic scaling parameters.
  * `fc_param` - The function compute parameters.
  * `pager_duty_param` - The PagerDuty parameters.
  * `fc3_param` - The FC3.0 parameters.
  * `eb_param` - The EventBridge parameters.
