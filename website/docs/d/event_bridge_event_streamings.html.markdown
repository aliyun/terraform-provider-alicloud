---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_streamings"
sidebar_current: "docs-alicloud-datasource-event-bridge-event-streamings"
description: |-
  Provides a list of Event Bridge Event Streamings to the user.
---

# alicloud\_event\_bridge\_event\_streamings

This data source provides the Event Bridge Event Streamings of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.189.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_event_bridge_event_streamings" "ids" {}
output "event_bridge_event_streaming_id_1" {
  value = data.alicloud_event_bridge_event_streamings.ids.streamings.0.id
}

data "alicloud_event_bridge_event_streamings" "nameRegex" {
  name_regex = "^my-EventStreaming"
}
output "event_bridge_event_streaming_id_2" {
  value = data.alicloud_event_bridge_event_streamings.nameRegex.streamings.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Event Streaming IDs. Its element value is same as Event Streaming Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Event Streaming name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `PAUSED`, `READY`, `RUNNING`, `RUNNING_FAILED`, `STARTING`, `STARTING_FAILED`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Event Streaming names.
* `streamings` - A list of Event Bridge Event Streamings. Each element contains the following attributes:
	* `description` - The description of the EventStreaming that is returned.
	* `id` - The id of the resource.
	* `event_streaming_name` - The name of the EventStreaming that is returned.
	* `filter_pattern` - The event rule. If this parameter is left empty, all events are matched.
	* `run_options` - The runtime environment.
		* `batch_window` - Batch push window.
			* `count_based_window` - Bulk Push Data Volume Window.
			* `time_based_window` - Batch push time window.
		* `dead_letter_queue` - Dead letter queue.
			* `arn` - Dead letter queue.
		* `errors_tolerance` - Exception tolerance strategy: NONE (intolerance of exceptions), ALL (tolerance of ALL exceptions).
		* `maximum_tasks` - Represents the number of concurrent tasks.
		* `retry_strategy` - Represents retry policy.
			* `maximum_event_age_in_seconds` - Maximum retry time.
			* `maximum_retry_attempts` - Maximum number of retries.
			* `push_retry_strategy` - Retry strategy: BACKOFF_RETRY (Backoff retry) and expential_decay_retry (exponential decay retry).
	* `sink` - Represents event flow data target.
		* `sink_kafka_parameters` - Represents Kafka event target parameters.
			* `key` - Represents message Key.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
			* `sasl_user` - SaslUser.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
			* `topic` - Represents Kafka Topic.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
			* `value` - Represents message body.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
				* `form` - The transformation method.
			* `acks` - Representative confirmation mode.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
			* `instance_id` - The ID of the Kafka instance.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
		* `sink_mns_parameters` - Represents the MNS event target parameter.
			* `queue_name` - The name of the queue specified when the event target is MNS.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
			* `body` - The content of the message.
				* `value` - The value before the transformation.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
			* `is_base64_encode` - Indicates whether Base64 encoding is enabled.
				* `value` - Indicates whether Base64 encoding is enabled.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
		* `sink_rabbit_mq_parameters` - Represents RabbitMQ event target parameter.
			* `instance_id` - The ID of the Message Queue for RabbitMQ instance specified when the event target is Message Queue for RabbitMQ.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - The ID of the Message Queue for RabbitMQ instance.
			* `queue_name` - The name of the queue to which events are pushed in the destination instance.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - The name of the queue in the Message Queue for RabbitMQ instance.
			* `routing_key` - The routing rule of the message.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - The routing rule of the message.
			* `target_type` - The type of the resource to which events are pushed.
				* `template` - Represents parameter Template.
				* `value` - The type of the resource to which events are pushed.
				* `form` - The transformation method.
			* `virtual_host_name` - The name of the vhost in the Message Queue for RabbitMQ instance.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - The name of the vhost in the Message Queue for RabbitMQ instance.
			* `body` - The content of the message.
				* `form` - The transformation method.
				* `template` - The template based on which events are transformed.
				* `value` - The value before the transformation.
			* `exchange` - The name of the exchange to which events are pushed in the destination instance.
				* `template` - Represents parameter Template.
				* `value` - The name of the exchange in the Message Queue for RabbitMQ instance.
				* `form` - The transformation method.
			* `message_id` - The ID of the message.
				* `form` - The transformation method.
				* `template` - The template based on which events are transformed.
				* `value` - The value before the transformation.
			* `properties` - The properties for filtering.
				* `template` - The template based on which events are transformed.
				* `value` - The value before the transformation.
				* `form` - The transformation method.
		* `sink_rocket_mq_parameters` - Represents the RocketMQ event target parameter.
			* `body` - The content of the message.
				* `template` - The template based on which events are transformed.
				* `value` - Represents parameter values.
				* `form` - TEMPLATE.
			* `instance_id` - The ID of the Message Queue for Apache RocketMQ instance specified when the event target is Message Queue for Apache RocketMQ.
				* `value` - The ID of the Message Queue for Apache RocketMQ instance.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
			* `keys` - The keys for filtering.
				* `form` - The transformation method.
				* `template` - The template based on which events are transformed.
				* `value` - The value before the transformation.
			* `properties` - The properties for filtering.
				* `form` - The transformation method.
				* `template` - The template based on which events are transformed.
				* `value` - The value before the transformation.
			* `tags` - The tags for filtering.
				* `form` - The transformation method.
				* `template` - The template based on which events are transformed.
				* `value` - The value before the transformation.
			* `topic` - The topic in the Message Queue for Apache RocketMQ instance.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - The name of the topic in the Message Queue for Apache RocketMQ instance.
		* `sink_sls_parameters` - Represents the SLS event target parameter.
			* `log_store` - Represents log store.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
			* `project` - Represents log project.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
			* `role_name` - Represents authorized role name.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
			* `topic` - Representative log Topic.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
				* `form` - The transformation method.
			* `body` - Representative event content.
				* `value` - Represents parameter values.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
		* `sink_fc_parameters` - Representative function calculates event target parameters.
			* `invocation_type` - Represents invocation type.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
				* `form` - The transformation method.
			* `qualifier` - Represents service version and alias.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
			* `service_name` - Represents service name.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
			* `body` - Representative event content.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
				* `form` - The transformation method.
			* `function_name` - Represents function name.
				* `form` - The transformation method.
				* `template` - Represents parameter Template.
				* `value` - Represents parameter values.
	* `source` - Represents an event flow data source.
		* `source_mqtt_parameters` - Represents the MQTT event source parameter.
			* `instance_id` - The ID of the MQTT instance.
			* `region_id` - The resource attribute field of region.
			* `topic` - Represents MQTT Topic.
		* `source_rabbit_mq_parameters` - The resource information configured when the event provider is Message Queue for RabbitMQ.
			* `region_id` - The resource attribute field representing the region.
			* `virtual_host_name` - The name of the vhost in the Message Queue for RabbitMQ instance.
			* `instance_id` - The ID of the Message Queue for RabbitMQ instance.
			* `queue_name` - The name of the queue in the Message Queue for RabbitMQ instance.
		* `source_rocket_mq_parameters` - The resource information configured when the event provider is Message Queue for Apache RocketMQ.
			* `topic` - The name of the topic in the Message Queue for Apache RocketMQ instance.
			* `group_id` - The group ID of the Message Queue for Apache RocketMQ instance.
			* `instance_id` - The ID of the Message Queue for Apache RocketMQ instance.
			* `offset` - The consumer offset of the message.
			* `region_id` - The ID of the region where the Message Queue for Apache RocketMQ instance resides.
			* `tag` - The tags for filtering.
			* `timestamp` - The timestamp of the offset from which the consumption starts.
		* `source_sls_parameters` - Represents SLS event source parameters.
			* `consume_position` - Represents the initial consumption offset.
			* `consumer_group` - Represents the consumer group name.
			* `log_store` - The log service that access logs are shipped to.
			* `project` - Represents log project.
			* `role_name` - Represents authorized role name.
		* `source_dts_parameters` - Representative DTS event source parameter.
			* `task_id` - The ID of the DTS task.
			* `topic` - Represents DTS Topic.
			* `username` - Represents the username of consumption group.
			* `broker_url` - Represents the DTS instance access point.
			* `init_check_point` - Represents the initial consumption offset.
			* `password` - Represents the password of consumption group.
			* `sid` - Represents the consumption group ID.
		* `source_kafka_parameters` - Representative Kafka event source parameter.
			* `offset_reset` - Represents the initial consumption locus.
			* `region_id` - The resource attribute field representing the region.
			* `topic` - Represents Kafka topic.
			* `consumer_group` - Represents Kafka consumption Group.
			* `instance_id` - The ID of the Kafka instance.
			* `network` - Represents network type.
			* `security_group_id` - Represents the security group ID.
			* `vswitch_ids` - The ID of the VSwitch.
			* `vpc_id` - The ID of the VPC.
		* `source_mns_parameters` - The resource information configured when the event provider is MNS.
			* `is_base64_decode` - Indicates whether Base64 encoding is enabled.
			* `queue_name` - The name of the queue in the MNS instance.
			* `region_id` - The resource attribute field representing the region.
	* `status` - The status of the resource.