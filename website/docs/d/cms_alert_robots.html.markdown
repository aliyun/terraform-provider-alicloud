---
subcategory: "Cloud Monitor (CMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_robots"
sidebar_current: "docs-alicloud-datasource-cms-alert-robots"
description: |-
  Provides a list of Cloud Monitor (CMS) Alert Robots.
---

# alicloud_cms_alert_robots

This data source provides a list of Cloud Monitor (CMS) Alert Robots in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_alert_robots" "example" {
  type = "DING"
}

output "first_robot_id" {
  value = data.alicloud_cms_alert_robots.example.robots.0.alert_robot_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of alert robot IDs to filter.
* `name_regex` - (Optional) A regex string to filter alert robots by name.
* `type` - (Optional) The type of the webhook. Valid values: `DING`, `WEIXIN`, `FEISHU`, `SLACK`, `TEAMS`, `DING_COOL_APP`.
* `workspace` - (Optional) The workspace ID to filter.
* `output_file` - (Optional) File path where to save the results.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of alert robot IDs.
* `robots` - A list of alert robots. Each element contains the following attributes:
  * `alert_robot_id` - The ID of the alert robot.
  * `alert_robot_name` - The name of the alert robot.
  * `type` - The type of the webhook.
  * `lang` - The language of the alert robot.
  * `url` - The webhook URL of the alert robot.
  * `workspace` - The workspace ID of the alert robot.
  * `digital_employee_name` - The name of the associated digital employee.
  * `robot_sign_key` - The sign key of the alert robot.
