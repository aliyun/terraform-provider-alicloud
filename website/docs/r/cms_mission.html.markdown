---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_mission"
description: |-
  Provides a Alicloud Cms Mission resource.
---

# alicloud_cms_mission

Provides a Cms Mission resource.

Mission is an async session for digital employees that contains multiple Blueprint schedule tasks (cron, calendar, or event triggered).

For information about Cms Mission and how to use it, see [What is Mission](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateMission).

-> **NOTE:** Available since v1.293.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cms_mission" "default" {
  name                  = var.name
  digital_employee_name = "my-digital-employee"
  display_name          = var.name
  description           = "terraform-example"
  enabled               = true

  variables = {
    workspace = var.name
  }

  notification_policy {
    region    = "cn-hangzhou"
    workspace = var.name
  }

  configuration {
    credits = 100
  }

  blueprint {
    cron {
      id                 = "${var.name}-cron"
      display_name       = "${var.name}-cron"
      description        = "cron blueprint"
      prompt             = "hello"
      cron_expression    = "0 0 * * * ?"
      time_zone          = "+0800"
      delay              = 0
      run_immediately    = false
      priority           = 1
      concurrency_policy = "skip"
      timeout_seconds    = 300
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the mission. It is the unique identifier and cannot be modified after creation. The name must match the pattern `^[a-zA-Z0-9_-]+$` with length between 1 and 64.
* `digital_employee_name` - (Required) The name of the associated digital employee.
* `display_name` - (Optional) The display name of the mission. Maximum length 512.
* `description` - (Optional) The description of the mission. Maximum length 512.
* `enabled` - (Optional) Whether the mission is enabled. Default to `true`.
* `variables` - (Optional) The variables of the mission. Must contain `workspace` or `project`.
* `notification` - (Optional) The notification configuration. See [`notification`](#notification) below.
* `notification_policy` - (Optional) The notification policy. See [`notification_policy`](#notification_policy) below.
* `blueprint` - (Optional) The blueprint schedule configuration. See [`blueprint`](#blueprint) below.
* `configuration` - (Optional) The runtime configuration. See [`configuration`](#configuration) below.

### `notification`

The notification block supports:

* `dingtalk` - (Optional) Dingtalk webhook URLs.
* `feishu` - (Optional) Feishu webhook URLs.
* `slack` - (Optional) Slack webhook URLs.
* `wechat` - (Optional) WeChat webhook URLs.
* `call` - (Optional) Call notification phone numbers.
* `sms` - (Optional) SMS notification phone numbers.
* `email` - (Optional) Email notification addresses.
* `webhook` - (Optional) Generic webhook URLs.

### `notification_policy`

The notification_policy block supports:

* `region` - (Optional) The notification policy region.
* `workspace` - (Optional) The notification policy workspace name.

### `configuration`

The configuration block supports:

* `credits` - (Optional) The credits limit for the mission.

### `blueprint`

The blueprint block supports:

* `cron` - (Optional) Cron blueprint list. See [`cron`](#blueprint-cron) below.
* `calendar` - (Optional) Calendar blueprint list. See [`calendar`](#blueprint-calendar) below.
* `event` - (Optional) Event blueprint list. See [`event`](#blueprint-event) below.

### `blueprint-cron`

The cron block supports:

* `id` - (Optional) Blueprint unique id.
* `display_name` - (Optional) Blueprint display name.
* `description` - (Optional) Blueprint description.
* `prompt` - (Optional) Agent prompt.
* `cron_expression` - (Optional) Cron expression.
* `variables` - (Optional) Variables passed to agent.
* `time_zone` - (Optional) Time zone. Default to `+0800`.
* `delay` - (Optional) Delay in seconds. Default to `0`.
* `run_immediately` - (Optional) Run immediately after create. Default to `false`.
* `priority` - (Optional) Priority 1-10. Default to `1`.
* `concurrency_policy` - (Optional) Concurrency policy: `skip`, `queue`, or `replace`. Default to `skip`.
* `timeout_seconds` - (Optional) Timeout in seconds.
* `enabled` - (Computed) Whether the cron blueprint is enabled.

### `blueprint-calendar`

The calendar block supports:

* `id` - (Optional) Blueprint unique id.
* `display_name` - (Optional) Blueprint display name.
* `description` - (Optional) Blueprint description.
* `prompt` - (Optional) Agent prompt.
* `rrule` - (Optional) iCalendar RRULE expression (RFC 5545).
* `variables` - (Optional) Variables.
* `time_zone` - (Optional) Time zone. Default to `+0800`.
* `delay` - (Optional) Delay in seconds. Default to `0`.
* `run_immediately` - (Optional) Run immediately. Default to `false`.
* `priority` - (Optional) Priority. Default to `1`.
* `concurrency_policy` - (Optional) Concurrency policy. Default to `skip`.
* `timeout_seconds` - (Optional) Timeout in seconds.
* `enabled` - (Computed) Whether the calendar blueprint is enabled.

### `blueprint-event`

The event block supports:

* `id` - (Optional) Blueprint unique id.
* `display_name` - (Optional) Blueprint display name.
* `description` - (Optional) Blueprint description.
* `prompt` - (Optional) Agent prompt.
* `workspace` - (Optional) Workspace name for event consumption.
* `max_concurrency` - (Optional) Max concurrent runs. Default to `1`.
* `debounce_seconds` - (Optional) Debounce in seconds. Default to `0`.
* `variables` - (Optional) Variables.
* `priority` - (Optional) Priority. Default to `1`.
* `concurrency_policy` - (Optional) Concurrency policy. Default to `skip`.
* `timeout_seconds` - (Optional) Timeout in seconds.
* `enabled` - (Computed) Whether the event blueprint is enabled.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the mission. It is the same as the `name`.
* `create_time` - The creation time of the mission.
* `update_time` - The last update time of the mission.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Mission.
* `delete` - (Defaults to 5 mins) Used when delete the Mission.
* `update` - (Defaults to 5 mins) Used when update the Mission.

## Import

Cms Mission can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_mission.example <name>
```
