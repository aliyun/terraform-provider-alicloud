---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_rule"
sidebar_current: "docs-alicloud-resource-event-bridge-rule"
description: |-
  Provides a Alicloud Event Bridge Rule resource.
---

# alicloud_event_bridge_rule

Provides a Event Bridge Rule resource.

For information about Event Bridge Rule and how to use it, see [What is Rule](https://www.alibabacloud.com/help/en/eventbridge/latest/createrule-6).

-> **NOTE:** Available since v1.129.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_account" "default" {
}

resource "alicloud_event_bridge_event_bus" "default" {
  event_bus_name = var.name
}

resource "alicloud_mns_queue" "default" {
  name = var.name
}

resource "alicloud_event_bridge_rule" "example" {
  event_bus_name = alicloud_event_bridge_event_bus.default.event_bus_name
  rule_name      = var.name
  filter_pattern = <<EOF
{
    "source": [
        "crmabc.newsletter"
    ],
    "type": [
        "UserSignUp",
        "UserLogin"
    ]
}
    EOF
  description    = var.name
  targets {
    target_id           = "tf-example"
    type                = "http"
    endpoint            = "http://www.aliyun.com"
    push_retry_strategy = "EXPONENTIAL_DECAY_RETRY"
    dead_letter_queue {
      arn = local.mns_endpoint
    }
    param_list {
      resource_key = "Body"
      form         = "ORIGINAL"
    }
    param_list {
      resource_key = "url"
      form         = "CONSTANT"
      value        = "http://www.aliyun.com"
    }
    param_list {
      resource_key = "Network"
      form         = "CONSTANT"
      value        = "PublicNetwork"
    }
  }
}

locals {
  mns_endpoint = format("acs:mns:%s:%s:queues/%s", data.alicloud_regions.default.regions.0.id, data.alicloud_account.default.id, alicloud_mns_queue.default.name)
}
```

## Argument Reference

The following arguments are supported:

* `event_bus_name` - (Required, ForceNew) The name of the event bus.
* `rule_name` - (Required, ForceNew) The name of the event rule.
* `filter_pattern` - (Required) The pattern to match interested events. Event mode, JSON format. The value description is as follows: `stringEqual` mode. `stringExpression` mode. Each field has up to 5 expressions (map structure).
* `description` - (Optional) The description of the event rule.
* `status` - (Optional) The status of the event rule. Valid values: `ENABLE`, `DISABLE`.
* `targets` - (Required, Set) The targets of rule. See [`targets`](#targets) below.

### `targets`

The targets supports the following:

* `target_id` - (Required, ForceNew) The ID of the custom event target.
* `type` - (Required) The type of the event target. Valid values: `acs.alikafka`, `acs.api.destination`, `acs.arms.loki`, `acs.datahub`, `acs.dingtalk`, `acs.eventbridge`, `acs.eventbridge.olap`, `acs.eventbus.SLSCloudLens`, `acs.fc.function`, `acs.fnf`, `acs.k8s`, `acs.mail`, `acs.mns.queue`, `acs.mns.topic`, `acs.openapi`, `acs.rabbitmq`, `acs.rds.mysql`, `acs.rocketmq`, `acs.sae`, `acs.sls`, `acs.sms`, `http`,`https` and `mysql`.
  **NOTE:** From version 1.208.1, `type` can be set to `acs.alikafka`, `acs.api.destination`, `acs.arms.loki`, `acs.datahub`, `acs.eventbridge.olap`, `acs.eventbus.SLSCloudLens`, `acs.fnf`, `acs.k8s`, `acs.openapi`, `acs.rds.mysql`, `acs.sae`, `acs.sls`, `mysql`.
* `endpoint` - (Required) The endpoint of the event target.
* `push_retry_strategy` - (Optional, Available since v1.184.0) The retry policy that is used to push the event. Valid values:
  - `BACKOFF_RETRY`: Backoff retry. The request can be retried up to three times. The interval between two consecutive retries is a random value between 10 and 20 seconds.
  - `EXPONENTIAL_DECAY_RETRY`: Exponential decay retry. The request can be retried up to 176 times. The interval between two consecutive retries exponentially increases to 512 seconds, and the total retry time is one day. The specific retry intervals are 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 512, ..., and 512 seconds, including a maximum of one hundred and sixty-seven 512 seconds in total.
* `dead_letter_queue` - (Optional, Set, Available since v1.184.0) The dead letter queue. Events that are not processed or exceed the number of retries will be written to the dead letter. Support message service MNS and message queue RocketMQ. See [`dead_letter_queue`](#targets-dead_letter_queue) below.
* `param_list` - (Required, Set) The parameters that are configured for the event target. See [`param_list`](#targets-param_list) below.

### `targets-dead_letter_queue`

The dead_letter_queue supports the following:

* `arn` - (Optional) The Alibaba Cloud Resource Name (ARN) of the dead letter queue. Events that are not processed or whose maximum retries are exceeded are written to the dead-letter queue. The ARN feature is supported by the following queue types: MNS and Message Queue for Apache RocketMQ.

### `targets-param_list`

The param_list supports the following:

* `resource_key` - (Required) The resource parameter of the event target. For more information, see [How to use it](https://www.alibabacloud.com/help/en/eventbridge/latest/event-target-parameters)
* `form` - (Required) The format of the event target parameter. Valid values: `ORIGINAL`, `TEMPLATE`, `JSONPATH`, `CONSTANT`.
* `template` - (Optional) The template of the event target parameter.
* `value` - (Optional) The value of the event target parameter.

-> **NOTE:** There exists a potential diff error that the backend service will return a default param as following:

```
param_list {
  resource_key = "IsBase64Encode"
  form         = "CONSTANT"
  value        = "false"
  template     = ""
}
```

In order to fix the diff, from version 1.160.0, this resource has removed the param which `resource_key = "IsBase64Encode"` and `value = "false"`.
If you want to set `resource_key = "IsBase64Encode"`, please avoid to set `value = "false"`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Rule. It formats as `<event_bus_name>:<rule_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Rule.
* `update` - (Defaults to 10 mins) Used when update the Rule.
* `delete` - (Defaults to 10 mins) Used when delete the Rule.

## Import

Event Bridge Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_rule.example <event_bus_name>:<rule_name>
```
