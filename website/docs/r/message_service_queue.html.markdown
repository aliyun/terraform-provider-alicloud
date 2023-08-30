---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_queue"
sidebar_current: "docs-alicloud-resource-message-service-queue"
description: |-
  Provides a Alicloud Message Notification Service Queue resource.
---

# alicloud_message_service_queue

Provides a Message Notification Service Queue resource.

For information about Message Notification Service Queue and how to use it, see [What is Queue](https://www.alibabacloud.com/help/en/message-service/latest/createqueue).

-> **NOTE:** Available since v1.188.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_message_service_queue" "queue" {
  queue_name               = var.name
  delay_seconds            = 60478
  maximum_message_size     = 12357
  message_retention_period = 256000
  visibility_timeout       = 30
  polling_wait_seconds     = 3
  logging_enabled          = true
}
```

## Argument Reference

The following arguments are supported:

* `queue_name` - (Required, ForceNew) Two queues on a single account in the same region cannot have the same name. A queue name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 120 characters.
* `delay_seconds` - (Optional) The delay period after which a message sent to the queue can be consumed. Unit: seconds. Valid values: 0-604800 seconds. Default value: 0.
* `maximum_message_size` - (Optional) The maximum size of a message body that can be sent to the queue. Unit: bytes. Valid value range: 1024-65536. Default value: 65536.
* `message_retention_period` - (Optional) The maximum period for which a message can be retained in the queue. After the specified period, the message is deleted no matter whether the message is consumed. Unit: seconds. Valid values: 60-604800. Default value: 345600.
* `visibility_timeout` - (Optional) The invisibility period for which the received message remains the Inactive state. Unit: seconds. Valid values: 1-43200. Default value: 30.
* `polling_wait_seconds` - (Optional) The maximum period for which a ReceiveMessage request waits if no message is available in the queue. Unit: seconds. Valid values: 0-30. Default value: 0.
* `logging_enabled` - (Optional) Specifies whether to enable the log management feature. Default value: false. Valid values:
  - `true`: enables the log management feature.
  - `false`: disables the log management feature.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Queue. Its value is same as `queue_name`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Queue.
* `update` - (Defaults to 3 mins) Used when update the Queue.
* `delete` - (Defaults to 3 mins) Used when delete the Queue.

## Import

Message Notification Service Queue can be imported using the id or queue_name, e.g.

```shell
$ terraform import alicloud_message_service_queue.example <queue_name>
```
