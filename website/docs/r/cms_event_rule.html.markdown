---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_event_rule"
sidebar_current: "docs-alicloud-resource-cms-event-rule"
description: |-
  Provides a Alicloud Cloud Monitor Service Event Rule resource.
---

# alicloud_cms_event_rule

Provides a Cloud Monitor Service Event Rule resource.

For information about Cloud Monitor Service Event Rule and how to use it, see [What is Event Rule](https://www.alibabacloud.com/help/en/cloudmonitor/latest/puteventrule).

-> **NOTE:** Available since v1.182.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_event_rule&exampleId=0d4c7024-3112-8273-2782-101ae743de50af352f78&activeTab=example&spm=docs.r.cms_event_rule.0.0d4c702431&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_cms_monitor_group" "default" {
  monitor_group_name = var.name
}

resource "alicloud_cms_event_rule" "example" {
  rule_name    = var.name
  group_id     = alicloud_cms_monitor_group.default.id
  silence_time = 100
  description  = var.name
  status       = "ENABLED"
  event_pattern {
    product         = "ecs"
    sql_filter      = "example_value"
    name_list       = ["example_value"]
    level_list      = ["CRITICAL"]
    event_type_list = ["StatusNotification"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `rule_name` - (Required, ForceNew) The name of the event-triggered alert rule.
* `group_id` - (Optional) The ID of the application group to which the event-triggered alert rule belongs.
* `silence_time` - (Optional, Int) The silence time.
* `description` - (Optional) The description of the event-triggered alert rule.
* `status` - (Optional) The status of the resource. Valid values: `ENABLED`, `DISABLED`.
* `event_pattern` - (Required, Set) Event mode, used to describe the trigger conditions for this event. See [`event_pattern`](#event_pattern) below.
* `contact_parameters` - (Optional, Set, Available since v1.211.1) The information about the alert contact groups that receive alert notifications. See [`contact_parameters`](#contact_parameters) below.
* `webhook_parameters` - (Optional, Set, Available since v1.211.1) The information about the callback URLs that are used to receive alert notifications. See [`webhook_parameters`](#webhook_parameters) below.
* `fc_parameters` - (Optional, Set, Available since v1.211.1) The information about the recipients in Function Compute. See [`fc_parameters`](#fc_parameters) below.
* `mns_parameters` - (Optional, Set, Available since v1.211.1) The information about the recipients in Message Service (MNS). See [`mns_parameters`](#mns_parameters) below.
* `sls_parameters` - (Optional, Set, Available since v1.211.1) The information about the recipients in Simple Log Service. See [`sls_parameters`](#sls_parameters) below.
* `open_api_parameters` - (Optional, Set, Available since v1.211.1) The parameters of API callback notification. See [`open_api_parameters`](#open_api_parameters) below.

### `event_pattern`

The event_pattern supports the following: 

* `product` - (Required) The type of the cloud service.
* `sql_filter` - (Optional) The SQL condition that is used to filter events. If the content of an event meets the specified SQL condition, an alert is automatically triggered.
* `name_list` - (Optional, List) The name of the event-triggered alert rule.
* `level_list` - (Optional, List) The level of the event-triggered alert rule. Valid values:
  - `CRITICAL`: critical.
  - `WARN`: warning.
  - `INFO`: information.
  - `*`: all types.
* `event_type_list` - (Optional, List) The type of the event-triggered alert rule. Valid values:
  - `StatusNotification`: fault notifications.
  - `Exception`: exceptions.
  - `Maintenance`: O&M.
  - `*`: all types.

### `contact_parameters`

The contact_parameters supports the following:

* `contact_parameters_id` (Optional) The ID of the recipient that receives alert notifications.
* `contact_group_name` (Optional) The name of the alert contact group.
* `level` (Optional) The alert level and the corresponding notification methods.

### `webhook_parameters`

The webhook_parameters supports the following:

* `webhook_parameters_id` (Optional) The ID of the recipient that receives alert notifications.
* `protocol` (Optional) The name of the protocol.
* `method` (Optional) The HTTP request method.
* `url` (Optional) The callback URL.

### `fc_parameters`

The fc_parameters supports the following:

* `fc_parameters_id` (Optional) The ID of the recipient that receives alert notifications.
* `service_name` (Optional) The name of the Function Compute service.
* `function_name` (Optional) The name of the function.
* `region` (Optional) The region where Function Compute is deployed.

### `sls_parameters`

The sls_parameters supports the following:

* `sls_parameters_id` (Optional) The ID of the recipient that receives alert notifications.
* `project` (Optional) The name of the Simple Log Service project.
* `log_store` (Optional) The name of the Simple Log Service Logstore.
* `region` (Optional) The region where Simple Log Service is deployed.

### `mns_parameters`

The mns_parameters supports the following:

* `mns_parameters_id` (Optional) The ID of the recipient that receives alert notifications.
* `queue` (Optional) The name of the MNS queue.
* `topic` (Optional) The MNS topic.
* `region` (Optional) The region where Message Service (MNS) is deployed.

### `open_api_parameters`

The open_api_parameters supports the following:

* `open_api_parameters_id` (Optional) The ID of the recipient that receives alert notifications sent by an API callback.
* `product` (Optional) The ID of the cloud service to which the API operation belongs.
* `action` (Optional) The API name.
* `version` (Optional) The version of the API.
* `role` (Optional) The name of the role.
* `region` (Optional) The region where the resource resides.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Event Rule. Its value is same as `rule_name`.
* `fc_parameters` - The information about the recipients in Function Compute.
  * `arn` - (Available since v1.211.1) The Alibaba Cloud Resource Name (ARN) of the function.
* `sls_parameters` - The information about the recipients in Log Service.
  * `arn` - (Available since v1.211.1) The ARN of the Log Service Logstore.
* `mns_parameters` - The information about the recipients in Message Service (MNS).
  * `arn` - (Available since v1.211.1) The ARN of the MNS queue.
* `open_api_parameters` - The information about the recipients in OpenAPI Explorer.
  * `arn` - (Available since v1.211.1) The ARN of the API operation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Cms Event Rule.
* `update` - (Defaults to 3 mins) Used when update the Cms Event Rule.
* `delete` - (Defaults to 3 mins) Used when delete the Cms Event Rule.

## Import

Cloud Monitor Service Event Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_event_rule.example <rule_name>
```
