---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_queue"
description: |-
  Provides a Alicloud Message Service Queue resource.
---

# alicloud_message_service_queue

Provides a Message Service Queue resource. 

For information about Message Service Queue and how to use it, see [What is Queue](https://www.alibabacloud.com/help/en/message-service/latest/createqueue).

-> **NOTE:** Available since v1.188.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_message_service_queue" "default" {
  delay_seconds            = 2
  polling_wait_seconds     = 2
  message_retention_period = 566
  maximum_message_size     = 1123
  visibility_timeout       = 30
  queue_name               = var.name
}
```

## Argument Reference

The following arguments are supported:
* `delay_seconds` - (Optional, Computed) The delay period after which a message sent to the queue can be consumed. Unit: seconds. Valid values: 0-604800 seconds. Default value: 0.
* `logging_enabled` - (Optional) Specifies whether to enable the log management feature. Default value: false. Valid values:
  - true: enables the log management feature.
  - false: disables the log management feature.
* `maximum_message_size` - (Optional, Computed) The maximum size of a message body that can be sent to the queue. Unit: bytes. Valid value range: 1024-65536. Default value: 65536.
* `message_retention_period` - (Optional, Computed) The maximum period for which a message can be retained in the queue. After the specified period, the message is deleted no matter whether the message is consumed. Unit: seconds. Valid values: 60-604800. Default value: 345600.
* `polling_wait_seconds` - (Optional, Computed)  The maximum period for which a ReceiveMessage request waits if no message is available in the queue. Unit: seconds. Valid values: 0-30. Default value: 0.
* `queue_name` - (Required, ForceNew) Two queues on a single account in the same region cannot have the same name. A queue name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 120 characters.
* `visibility_timeout` - (Optional, Computed) The invisibility period for which the received message remains the Inactive state. Unit: seconds. Valid values: 1-43200. Default value: 30.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Represents the time when the Queue was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Queue.
* `delete` - (Defaults to 5 mins) Used when delete the Queue.
* `update` - (Defaults to 5 mins) Used when update the Queue.

## Import

Message Service Queue can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_queue.example <id>
```