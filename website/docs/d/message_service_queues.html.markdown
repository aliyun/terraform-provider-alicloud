---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_queues"
sidebar_current: "docs-alicloud-datasource-message-service-queues"
description: |-
  Provides a list of Message Notification Service Queues to the user.
---

# alicloud\_message\_service\_queues

This data source provides the Message Notification Service Queues of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.188.0+.

## Example Usage

```terraform
data "alicloud_message_service_queues" "ids" {
  ids = ["example_id"]
}

output "queue_id_1" {
  value = data.alicloud_message_service_queues.ids.queues.0.id
}

data "alicloud_message_service_queues" "name" {
  queue_name = "tf-example"
}

output "queue_id_2" {
  value = data.alicloud_message_service_queues.name.queues.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Queue IDs. Its element value is same as Queue Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Queue name.
* `queue_name` - (Optional, ForceNew) The name of the queue.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Queue names. 
* `queues` - A list of Queues. Each element contains the following attributes:
  * `id` - The ID of the Queue. Its value is same as Queue Name.
  * `queue_name` - The name of the queue.
  * `delay_seconds` - The delay period after which all messages that are sent to the queue can be consumed. Unit: seconds.
  * `maximum_message_size` - The maximum size of a message body that can be sent to the queue. Unit: bytes.
  * `message_retention_period` - The maximum period for which a message can be retained in the queue. A message that is sent to the queue can be retained for a specified period. After the specified period ends, the message is deleted no matter whether it is consumed. Unit: seconds.
  * `visibility_timeout` - The invisibility period for which the received message remains the Inactive state. Unit: seconds.
  * `polling_wait_seconds` - The maximum period for which a ReceiveMessage request waits if no message is available in the queue. Unit: seconds.
  * `logging_enabled` - Indicates whether the log management feature is enabled for the queue.
  * `active_messages` - The total number of messages that are in the Active state in the queue. The value is an approximate number.
  * `inactive_messages` - The total number of the messages that are in the Inactive state in the queue. The value is an approximate number.
  * `delay_messages` - The total number of the messages that are in the Delayed state in the queue. The value is an approximate number.
  * `queue_url` - The url of the queue.
  * `queue_internal_url` - The internal url of the queue.
  * `last_modify_time` - The time when the queue was last modified. This value is a UNIX timestamp representing the number of milliseconds that have elapsed since the epoch time January 1, 1970, 00:00:00 UTC.
  * `create_time` - The time when the queue was created. This value is a UNIX timestamp representing the number of milliseconds that have elapsed since the epoch time January 1, 1970, 00:00:00 UTC.