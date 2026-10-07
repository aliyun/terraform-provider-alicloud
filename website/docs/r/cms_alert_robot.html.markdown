---
subcategory: "Cloud Monitor (CMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_robot"
sidebar_current: "docs-alicloud-resource-cms-alert-robot"
description: |-
  Provides a Cloud Monitor (CMS) Alert Robot resource.
---

# alicloud_cms_alert_robot

Provides a Cloud Monitor (CMS) Alert Robot resource.

For information about CMS Alert Robot and how to use it, see [What is Alert Robot](https://www.alibabacloud.com/help/en/cms/developer-reference/api-cms-2024-03-30-createalertrobot).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_alert_robot" "example" {
  alert_robot_name      = "example-robot"
  type                  = "DING"
  url                   = "https://oapi.dingtalk.com/robot/send?access_token=xxxxx"
  lang                  = "zh_CN"
  digital_employee_name = "example-employee"
  robot_sign_key        = "example-sign-key"
}
```

## Argument Reference

The following arguments are supported:

* `alert_robot_name` - (Optional) The name of the alert robot.
* `type` - (Required, ForceNew) The type of the webhook. Valid values: `DING`, `WEIXIN`, `FEISHU`, `SLACK`, `TEAMS`, `DING_COOL_APP`.
* `url` - (Required) The webhook URL of the alert robot.
* `lang` - (Optional) The language of the alert robot. Valid values: `zh_CN`, `en_US`. Default is `zh_CN`.
* `workspace` - (Optional, Computed, ForceNew) The workspace ID of the alert robot.
* `digital_employee_name` - (Optional) The name of the associated digital employee.
* `robot_sign_key` - (Optional, Sensitive) The sign key of the alert robot.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alert robot.
* `alert_robot_id` - The ID of the alert robot, same as `id`.

## Import

CMS Alert Robot can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_alert_robot.example <robot-id>
```
