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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_event_bridge_rule&exampleId=7f7ab0bb-cb71-9813-c11f-ff55ad4c6d7f4f9120a4&activeTab=example&spm=docs.r.event_bridge_rule.0.7f7ab0bbcb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "tf-example"
}

data "alicloud_account" "default" {
}

resource "alicloud_event_bridge_event_bus" "default" {
  event_bus_name = var.name
}

resource "alicloud_mns_queue" "queue1" {
  name = var.name
}

locals {
  mns_endpoint_a = format("acs:mns:cn-hangzhou:%s:queues/%s", data.alicloud_account.default.id, alicloud_mns_queue.queue1.name)
  fnf_endpoint   = format("acs:fnf:cn-hangzhou:%s:flow/$${flow}", data.alicloud_account.default.id)
}
resource "alicloud_event_bridge_rule" "example" {
  event_bus_name = alicloud_event_bridge_event_bus.default.event_bus_name
  rule_name      = var.name
  description    = "example"
  filter_pattern = "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}"
  targets {
    target_id = "tf-example1"
    endpoint  = local.mns_endpoint_a
    type      = "acs.mns.queue"
    param_list {
      resource_key = "queue"
      form         = "CONSTANT"
      value        = "tf-testaccEbRule"
    }
    param_list {
      resource_key = "Body"
      form         = "ORIGINAL"
    }
    param_list {
      form         = "CONSTANT"
      resource_key = "IsBase64Encode"
      value        = "true"

    }
  }
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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Rule.
* `update` - (Defaults to 10 mins) Used when update the Rule.
* `delete` - (Defaults to 10 mins) Used when delete the Rule.

## Import

Event Bridge Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_rule.example <event_bus_name>:<rule_name>
```
