---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_source_v2"
description: |-
  Provides a Alicloud Event Bridge Event Source V2 resource.
---

# alicloud_event_bridge_event_source_v2

Provides a Event Bridge Event Source V2 resource.



For information about Event Bridge Event Source V2 and how to use it, see [What is Event Source V2](https://www.alibabacloud.com/help/en/eventbridge/latest/api-eventbridge-2020-04-01-createeventsource).

-> **NOTE:** Available since v1.269.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_event_bridge_event_bus" "default" {
  event_bus_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_event_bridge_event_source_v2" "default" {
  event_bus_name         = alicloud_event_bridge_event_bus.default.event_bus_name
  event_source_name      = "${var.name}-${random_integer.default.result}"
  description            = var.name
  linked_external_source = true
  source_http_event_parameters {
    type            = "HTTP"
    security_config = "referer"
    method          = ["GET", "POST", "DELETE"]
    referer         = ["www.aliyun.com", "www.alicloud.com", "www.taobao.com"]
  }
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The detail describe of event source
* `event_bus_name` - (Required, ForceNew) Name of the bus associated with the event source
* `event_source_name` - (Required, ForceNew) The code name of event source
* `linked_external_source` - (Optional, Bool) Whether to connect to an external data source

-> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `source_http_event_parameters` - (Optional, Set) The request parameter SourceHttpEventParameters. See [`source_http_event_parameters`](#source_http_event_parameters) below.
* `source_kafka_parameters` - (Optional, Set) Kafka event source parameter. See [`source_kafka_parameters`](#source_kafka_parameters) below.
* `source_mns_parameters` - (Optional, Set) Lightweight message queue (formerly MNS) event source parameter. See [`source_mns_parameters`](#source_mns_parameters) below.
* `source_oss_event_parameters` - (Optional, Set) OSS event source parameters See [`source_oss_event_parameters`](#source_oss_event_parameters) below.
* `source_rabbit_mq_parameters` - (Optional, Set) The request parameter SourceRabbitMQParameters. See [`source_rabbit_mq_parameters`](#source_rabbit_mq_parameters) below.
* `source_rocketmq_parameters` - (Optional, Set) The request parameter SourceRocketMQParameters. See [`source_rocketmq_parameters`](#source_rocketmq_parameters) below.
* `source_sls_parameters` - (Optional, ForceNew, Set) The request parameter SourceSLSParameters. See [`source_sls_parameters`](#source_sls_parameters) below.
* `source_scheduled_event_parameters` - (Optional, Set) Time event source parameter. See [`source_scheduled_event_parameters`](#source_scheduled_event_parameters) below.

### `source_http_event_parameters`

The source_http_event_parameters supports the following:
* `ip` - (Optional, List) IP segment security configuration. This parameter must be set only when the SecurityConfig value is ip. You can enter an IP address segment or IP address.
* `method` - (Optional, List) The HTTP request method supported by the generated Webhook. Multiple choices are available, with the following options:
  - GET
  - POST
  - PUT
  - PATCH
  - DELETE
  - HEAD
  - OPTIONS
  - TRACE
  - CONNECT
* `referer` - (Optional, List) Security domain name configuration. This parameter must be set only when SecurityConfig is set to referer. You can fill in the domain name.
* `security_config` - (Optional) Select the type of security configuration. The optional range is as follows:
  - none: No configuration is required.
  - ip:IP segment.
  - referer: Security domain name.
* `type` - (Optional) The protocol type supported by the generated Webhook. The value description is as follows:
  - HTTP
  - HTTPS
  - HTTP&HTTPS

### `source_kafka_parameters`

The source_kafka_parameters supports the following:
* `consumer_group` - (Optional) The Group ID of the consumer who subscribes to the Topic.
* `instance_id` - (Optional) The instance ID.
* `network` - (Optional) Network configuration: Default (Default network) and public network (self-built network).
* `offset_reset` - (Optional) Consumption sites.
* `region_id` - (Optional, ForceNew) The region ID.
* `security_group_id` - (Optional) The ID of the security group.
* `topic` - (Optional) The topic name.
* `vswitch_ids` - (Optional) The vSwitch ID.
* `vpc_id` - (Optional) The VPC ID.

### `source_mns_parameters`

The source_mns_parameters supports the following:
* `is_base64_decode` - (Optional) Whether to enable Base64 decoding. By default, it is selected, that is, Base64 decoding is enabled.
* `queue_name` - (Optional) The name of the Queue of the lightweight message Queue (formerly MNS).
* `region_id` - (Optional) The region of the lightweight message queue (formerly MNS).

### `source_oss_event_parameters`

The source_oss_event_parameters supports the following:
* `event_types` - (Optional, ForceNew, List) OSS event type list.
* `match_rules` - (Optional, ForceNew, List) Matching rules. The event source will deliver OSS events that meet the matching requirements to the bus.
* `sts_role_arn` - (Optional) The ARN of the role. EventBridge will use this role to create MNS resources and deliver events to the corresponding bus.

### `source_rabbit_mq_parameters`

The source_rabbit_mq_parameters supports the following:
* `instance_id` - (Optional) The ID of the RabbitMQ instance. For more information, see Usage Restrictions (~~ 163289 ~~).
* `queue_name` - (Optional) The name of the Queue of the RabbitMQ instance. For more information, see Usage Restrictions (~~ 163289 ~~).
* `region_id` - (Optional) The region of the RabbitMQ instance.
* `virtual_host_name` - (Optional) The name of the Vhost of the RabbitMQ instance. For more information, see Usage Restrictions (~~ 163289 ~~).

### `source_rocketmq_parameters`

The source_rocketmq_parameters supports the following:
* `auth_type` - (Optional) ACL or not.
* `group_id` - (Optional) The Group ID of the RocketMQ version of message queue.
* `instance_endpoint` - (Optional) Instance access point.
* `instance_id` - (Optional) The ID of the RocketMQ instance. For more information, see Usage Restrictions (~~ 163289 ~~).
* `instance_network` - (Optional) Instance network.
* `instance_password` - (Optional) The instance password.
* `instance_security_group_id` - (Optional) The ID of the security group.
* `instance_type` - (Optional) The instance type. Only CLOUD_4 (4.0 instance on the cloud), CLOUD_5 (5.0 instance on the cloud), and SELF_BUILT (user-created MQ).
* `instance_username` - (Optional) The instance user name.
* `instance_vswitch_ids` - (Optional) The vSwitch ID.
* `instance_vpc_id` - (Optional) The ID of the VPC.
* `offset` - (Optional) The consumption point of the message. The value description is as follows:
  - `CONSUME_FROM_LAST_OFFSET`: starts consumption from the latest point.
  - `CONSUME_FROM_FIRST_OFFSET`: starts consumption from the earliest point.
  - `CONSUME_FROM_TIMESTAMP`: starts consumption from the specified time point.
Default value: `CONSUME_FROM_LAST_OFFSET`.
* `region_id` - (Optional) The region of the RocketMQ instance.
* `tag` - (Optional) The filter label of the message.
* `timestamp` - (Optional, Float) The timestamp. This parameter is valid only when the value of the Offset parameter is CONSUME_FROM_TIMESTAMP.
* `topic` - (Optional) The Topic name of the RocketMQ instance. For more information, see Usage Restrictions (~~ 163289 ~~).

### `source_sls_parameters`

The source_sls_parameters supports the following:
* `consume_position` - (Optional, ForceNew) Start consumption point, which can be the earliest or latest point corresponding to begin and end respectively, or start consumption from a specified time, measured in seconds.
* `log_store` - (Optional, ForceNew) The logstore of log service SLS.
* `project` - (Optional, ForceNew) The log project of log service SLS.
* `role_name` - (Optional, ForceNew) When authorizing event bus EventBridge to use this role to read SLS log content, the following conditions must be met: when creating the role used by the service in the RAM console, you need to select Alibaba Cloud Service and event bus for trusted service ". For the permissions policy of this role, see custom event source log service SLS.

### `source_scheduled_event_parameters`

The source_scheduled_event_parameters supports the following:
* `schedule` - (Optional) Cron expression
* `time_zone` - (Optional) The Cron execution time zone.
* `user_data` - (Optional) JSON string

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `source_http_event_parameters` - The request parameter SourceHttpEventParameters.
  * `public_web_hook_url` - The public network request URL.
  * `vpc_web_hook_url` - The intranet request URL.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Event Source V2.
* `delete` - (Defaults to 5 mins) Used when delete the Event Source V2.
* `update` - (Defaults to 5 mins) Used when update the Event Source V2.

## Import

Event Bridge Event Source V2 can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_event_source_v2.example <id>
```
