---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_rule"
sidebar_current: "docs-alicloud-resource-event-bridge-rule"
description: |-
  Provides a Alicloud Event Bridge Rule resource.
---

# alicloud\_event\_bridge\_rule

Provides a Event Bridge Rule resource.

For information about Event Bridge Rule and how to use it, see [What is Rule](https://help.aliyun.com/document_detail/167854.html).

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_event_bridge_event_bus" "example" {
  event_bus_name = "example_value"
}

resource "alicloud_event_bridge_rule" "example" {
  event_bus_name = alicloud_event_bridge_event_bus.example.id
  rule_name      = var.name
  description    = "test"
  filter_pattern = "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}"
  targets {
    target_id = "tf-test"
    endpoint  = "acs:mns:cn-hangzhou:118938335****:queues/tf-test"
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
  }
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of rule.
* `event_bus_name` - (Required, ForceNew) The name of event bus.
* `filter_pattern` - (Required) The pattern to match interested events. Event mode, JSON format. The value description is as follows: `stringEqual` mode. `stringExpression` mode. Each field has up to 5 expressions (map structure).
* `rule_name` - (Required, ForceNew) The name of rule.
* `status` - (Optional, Computed) Rule status, either Enable or Disable. Valid values: `DISABLE`, `ENABLE`.
* `targets` - (Required) The target of rule.

### Block targets

* `endpoint` - (Required) The endpoint of target.
* `param_list` - (Required) A list of param.
* `target_id` - (Required, ForceNew) The ID of target.
* `type` - (Required) The type of target. Valid values: `acs.fc.function`, `acs.mns.topic`, `acs.mns.queue`,`http`,`acs.sms`,`acs.mail`,`acs.dingtalk`,`https`, `acs.eventbridge`,`acs.rabbitmq` and `acs.rocketmq`.

### Block param_list

* `resource_key` - (Required) The resource key of param.  For more information, see [Event target parameters](https://help.aliyun.com/document_detail/185887.htm)
* `form` - (Required) The format of param.  Valid values: `ORIGINAL`, `TEMPLATE`, `JSONPATH`, `CONSTANT`.
* `value` - (Optional) The value of param.
* `template` - (Optional) The template of param.

-> **NOTE:** There exists a potential diff error that the backend service will return a default param as following:
```terraform
param_list {
  resource_key = "IsBase64Encode"
  form         = "CONSTANT"
  value        = "false"
  template     = ""
}
```
In order to fix the diff, from version 1.160.0, 
this resource has removed the param which `resource_key = "IsBase64Encode"` and `value = "false"`.
If you want to set `resource_key = "IsBase64Encode"`, please avoid to set `value = "false"`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Rule. The value is formatted `<event_bus_name>:<rule_name>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Rule.
* `delete` - (Defaults to 10 mins) Used when delete the Rule.
* `update` - (Defaults to 10 mins) Used when update the Rule.

## Import

Event Bridge Rule can be imported using the id, e.g.

```
$ terraform import alicloud_event_bridge_rule.example <event_bus_name>:<rule_name>
```
