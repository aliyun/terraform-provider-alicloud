---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_alert_robot"
sidebar_current: "docs-alicloud-resource-arms-alert-robot"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Alert Robot resource.
---

# alicloud_arms_alert_robot

Provides a Application Real-Time Monitoring Service (ARMS) Alert Robot resource.

For information about Application Real-Time Monitoring Service (ARMS) Alert Robot and how to use it, see [What is Alert Robot](https://next.api.alibabacloud.com/document/ARMS/2019-08-08/CreateOrUpdateIMRobot).

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_arms_alert_robot&exampleId=bb18d5b9-0088-3372-d9c0-fc554d0d65c5fecbdba1&activeTab=example&spm=docs.r.arms_alert_robot.0.bb18d5b900&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_arms_alert_robot" "wechat" {
  alert_robot_name = "example_wechat"
  robot_type       = "wechat"
  robot_addr       = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=1c704e23"
  daily_noc        = true
  daily_noc_time   = "09:30,17:00"
}

resource "alicloud_arms_alert_robot" "dingding" {
  alert_robot_name = "example_dingding"
  robot_type       = "dingding"
  robot_addr       = "https://oapi.dingtalk.com/robot/send?access_token=1c704e23"
  daily_noc        = true
  daily_noc_time   = "09:30,17:00"
}

resource "alicloud_arms_alert_robot" "feishu" {
  alert_robot_name = "example_feishu"
  robot_type       = "feishu"
  robot_addr       = "https://open.feishu.cn/open-apis/bot/v2/hook/a48efa01"
  daily_noc        = true
  daily_noc_time   = "09:30,17:00"
}
```

## Argument Reference

The following arguments are supported:

* `alert_robot_name` - (Required) The name of the resource.
* `robot_type` - (Required, ForceNew) The type of the robot, Valid values: `wechat`, `dingding`, `feishu`.
* `robot_addr` - (Required) The webhook url of the robot.
* `daily_noc` - (Optional) Specifies whether the alert robot receives daily notifications. Valid values: `true`: receives daily notifications. `false`: does not receive daily notifications, default to `false`.
* `daily_noc_time` - (Optional) The time of the daily notification.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Alert Robot.

## Import

Application Real-Time Monitoring Service (ARMS) Alert Robot can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_alert_robot.example <id>
```
