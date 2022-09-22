---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_streaming"
sidebar_current: "docs-alicloud-resource-event-bridge-event-streaming"
description: |-
  Provides a Alicloud Event Bridge Event Streaming resource.
---

# alicloud\_event\_bridge\_event\_streaming

Provides a Event Bridge Event Streaming resource.

For information about Event Bridge Event Streaming and how to use it, see [What is Event Streaming](https://www.alibabacloud.com/help/en/eventbridge/latest/createeventstreaming).

-> **NOTE:** Available in v1.189.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_event_bridge_event_streaming" "default" {
  event_streaming_name = var.name
  description          = var.name
  source {
    source_mns_parameters {
      queue_name       = "test"
      is_base64_decode = "true"
    }
  }
  sink {
    sink_mns_parameters {
      queue_name {
        value = "test"
        form  = "CONSTANT"
      }
      body {
        value = "$.data"
        form  = "JSONPATH"
      }
      is_base64_encode {
        value = "true"
        form  = "CONSTANT"
      }
    }
  }
  run_options {
    errors_tolerance = "ALL"
    retry_strategy {
      push_retry_strategy = "BACKOFF_RETRY"
    }
  }
  filter_pattern = "{}"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Required) The description of the EventStreaming.
* `event_streaming_name` - (Required, ForceNew) The name of the EventStreaming.
* `filter_pattern` - (Required) The event rule. If you leave this parameter empty, all events are matched.
* `run_options` - (Optional) The runtime environment. See the following `Block run_options`.
* `sink` - (Required) The information about the event target. See the following `Block sink`.
* `source` - (Required) The information about the event provider. See the following `Block source`.

#### Block source

The source supports the following: 

* `source_dts_parameters` - (Optional) The resource information that you configure when the event provider is Dts.
* `source_kafka_parameters` - (Optional) The resource information that you configure when the event provider is Kafka. See the following `Block source_kafka_parameters`.
* `source_mns_parameters` - (Optional) The resource information that you configure when the event provider is MNS. See the following `Block source_mns_parameters`.
* `source_mqtt_parameters` - (Optional) The resource information that you configure when the event provider is Message Queue for MQTT. See the following `Block source_mqtt_parameters`.
* `source_rabbit_mq_parameters` - (Optional) The resource information that you configure when the event provider is Message Queue for RabbitMQ. See the following `Block source_rabbit_mq_parameters`.
* `source_rocket_mq_parameters` - (Optional) The resource information that you configure when the event provider is Message Queue for Apache RocketMQ. See the following `Block source_rocket_mq_parameters`.
* `source_sls_parameters` - (Optional) The resource information that you configure when the event provider is Message Queue for SLS. See the following `Block source_sls_parameters`.

#### Block source_dts_parameters

The source_dts_parameters supports the following:

* `broker_url` - (Optional) Represents the DTS instance access point.
* `init_check_point` - (Optional) Represents the initial consumption offset.
* `password` - (Optional) The password of the consumer group account.
* `sid` - (Optional) Represents the consumption group ID.
* `task_id` - (Optional) The ID of the DTS task.
* `topic` - (Optional) The DTS topic.
* `username` - (Optional) The username of the consumer group.

#### Block source_kafka_parameters

The source_kafka_parameters supports the following: 

* `consumer_group` - (Optional) Represents Kafka consumption Group.
* `instance_id` - (Optional) The ID of the Kafka instance.
* `network` - (Optional) Represents network type. Valid values: `PublicNetwork`, `Default`.
* `offset_reset` - (Optional) Represents the initial consumption locus.
* `region_id` - (Optional) The resource attribute field representing the region.
* `security_group_id` - (Optional) The ID of the security group ID.
* `topic` - (Optional) The Kafka topic name.
* `vpc_id` - (Optional) The ID of the VPC.
* `vswitch_ids` - (Optional) The vswitch ids.

#### Block source_mns_parameters

The source_mns_parameters supports the following:

* `is_base64_decode` - (Optional) Represents whether base64 encoding is turned on.
* `region_id` - (Optional) The resource attribute field representing the region.
* `queue_name` - (Optional) Represents MNS queue name.

#### Block source_mqtt_parameters

The source_mqtt_parameters supports the following:

* `instance_id` - (Optional) The ID of the MQTT instance.
* `region_id` - (Optional) The resource attribute field representing the region.
* `topic` - (Optional) The MQTT topic name.

#### Block source_rabbit_mq_parameters

The source_rabbit_mq_parameters supports the following:

* `region_id` - (Optional) The resource attribute field representing the region.
* `virtual_host_name` - (Optional) Represents the virtual host name.
* `instance_id` - (Optional) Represents the RabbitMQ instance ID.
* `queue_name` - (Optional) Represents RabbitMQ queue name.

#### Block source_rocket_mq_parameters

The source_rocket_mq_parameters supports the following:

* `topic` - (Optional) Represents RocketMQ topic.
* `group_id` - (Optional) Represents RocketMQ Group.
* `instance_id` - (Optional) The ID of the RocketMQ instance.
* `offset` - (Optional) Representative consumption offset. Valid values: `CONSUME_FROM_LAST_OFFSET`, `CONSUME_FROM_FIRST_OFFSET`, `CONSUME_FROM_TIMESTAMP`.
* `region_id` - (Optional) Resource attribute field of region.
* `tag` - (Optional) Representative tag.
* `timestamp` - (Optional) Represents consumption timestamp.

#### Block source_sls_parameters

The source_sls_parameters supports the following:

* `consume_position` - (Optional) Represents the initial consumption offset.
* `log_store` - (Optional) The log service that access logs are shipped to.
* `project` - (Optional) Represents log project.
* `role_name` - (Optional) Represents authorized role name.

#### Block sink

The sink supports the following: 

* `sink_fc_parameters` - (Optional) Representative function calculates event target parameters. See the following `Block sink_fc_parameters`.
* `sink_kafka_parameters` - (Optional) The resource information that you configure when the event target is Message Queue for Kafka. See the following `Block sink_kafka_parameters`.
* `sink_mns_parameters` - (Optional) The resource information that you configure when the event target is MNS. See the following `Block sink_mns_parameters`.
* `sink_rabbit_mq_parameters` - (Optional) The resource information that you configure when the event target is Message Queue for RabbitMQ. See the following `Block sink_rabbit_mq_parameters`.
* `sink_rocket_mq_parameters` - (Optional) The resource information that you configure when the event target is Message Queue for Apache RocketMQ. See the following `Block sink_rocket_mq_parameters`.
* `sink_sls_parameters` - (Optional) The resource information that you configure when the event target is Message Queue for SLS. See the following `Block sink_sls_parameters`.

#### Block sink_fc_parameters

The sink_fc_parameters supports the following:

* `body` - (Optional) Representative event content. See the following `Block body`.
* `function_name` - (Optional) Represents function name. See the following `Block body`.
* `invocation_type` - (Optional) Represents invocation type. See the following `Block body`.
* `qualifier` - (Optional) Represents service version and alias. See the following `Block body`.
* `service_name` - (Optional) Represents service name. See the following `Block body`.

#### Block sink_kafka_parameters

The sink_kafka_parameters supports the following: 

* `acks` - (Optional) Representative confirmation mode. See the following `Block body`.
* `instance_id` - (Optional) The ID of the Kafka instance. See the following `Block body`.
* `key` - (Optional) Represents message Key. See the following `Block body`.
* `sasl_user` - (Optional) SaslUser. See the following `Block body`.
* `topic` - (Optional) Represents Kafka Topic. See the following `Block body`.
* `value` - (Optional) Represents message body. See the following `Block body`.

#### Block sink_mns_parameters

The sink_mns_parameters supports the following: 

* `body` - (Optional) Representative event content. See the following `Block body`.
* `queue_name` - (Optional) Represents MNS queue name. See the following `Block body`.
* `is_base64_encode` - (Optional) Represents whether Base64 encoding is turned on. See the following `Block body`.

#### Block sink_rabbit_mq_parameters

The sink_rabbit_mq_parameters supports the following:

* `body` - (Optional) Representative event content. See the following `Block body`.
* `instance_id` - (Optional) Represents the RabbitMQ instance ID. See the following `Block body`.
* `queue_name` - (Optional) Represents RabbitMQ queue name. See the following `Block body`.
* `routing_key` - (Optional) Represents routing keyword. See the following `Block body`.
* `target_type` - (Optional) Represents target type. See the following `Block body`.
* `virtual_host_name` - (Optional) Represents the virtual host name. See the following `Block body`.
* `exchange` - (Optional) Represent the message routing agent. See the following `Block body`.
* `message_id` - (Optional) Representative message ID. See the following `Block body`.
* `properties` - (Optional) Represents a custom attribute. See the following `Block body`.

#### Block sink_rocket_mq_parameters

The sink_rocket_mq_parameters supports the following:

* `body` - (Optional) Representative event content. See the following `Block body`.
* `instance_id` - (Optional) The ID of the RocketMQ instance. See the following `Block body`.
* `keys` - (Optional) Represents message Key. See the following `Block body`.
* `properties` - (Optional) Represents message attributes. See the following `Block body`.
* `tags` - (Optional) Representative tag. See the following `Block body`.
* `topic` - (Optional) Represents RocketMQ Topic. See the following `Block body`.

#### Block sink_sls_parameters

The sink_sls_parameters supports the following:

* `body` - (Optional) Representative event content. See the following `Block body`.
* `log_store` - (Optional) Represents logstore. See the following `Block body`.
* `project` - (Optional) Represents log project. See the following `Block body`.
* `role_name` - (Optional) Represents authorized role name. See the following `Block body`.
* `topic` - (Optional) Representative log Topic. See the following `Block body`.

#### Block body

The body supports the following: 

* `form` - (Optional) Represents parameter format.
* `template` - (Optional) Represents parameter Template.
* `value` - (Optional) Represents parameter values.

#### Block run_options

The run_options supports the following: 

* `batch_window` - (Optional) Batch push window. See the following `Block batch_window`.
* `dead_letter_queue` - (Optional) Specifies whether to enable dead-letter queues. See the following `Block dead_letter_queue`.
* `errors_tolerance` - (Optional) Exception tolerance strategy. Valid values: `NONE` (intolerance of exceptions), `ALL` (tolerance of ALL exceptions).
* `maximum_tasks` - (Optional) Represents the number of concurrent tasks.
* `retry_strategy` - (Optional) Represents retry policy. See the following `Block retry_strategy`.

#### Block retry_strategy

The retry_strategy supports the following: 

* `maximum_event_age_in_seconds` - (Optional) Maximum retry time.
* `maximum_retry_attempts` - (Optional) Maximum number of retries.
* `push_retry_strategy` - (Optional) The retry policy to be used when an event fails to be pushed. Valid values: `BACKOFF_RETRY`, `EXPONENTIAL_DECAY_RETRY`.

#### Block dead_letter_queue

The dead_letter_queue supports the following:

* `arn` - (Optional) Dead letter queue.

#### Block batch_window

The batch_window supports the following: 

* `count_based_window` - (Optional) Bulk Push Data Volume Window.
* `time_based_window` - (Optional) Batch push time window.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Event Streaming. Its value is same as `event_streaming_name`.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Event Streaming.
* `update` - (Defaults to 1 mins) Used when update the Event Streaming.
* `delete` - (Defaults to 1 mins) Used when delete the Event Streaming.


## Import

Event Bridge Event Streaming can be imported using the id, e.g.

```
$ terraform import alicloud_event_bridge_event_streaming.example <event_streaming_name>
```