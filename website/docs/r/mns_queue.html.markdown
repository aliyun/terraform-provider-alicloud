--
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_queue"
sidebar_current: "docs-alicloud-resource-mns-queue"
description: |-
  Provides a Alicloud MNS Queue resource.
---

# alicloud\_mns\_queue

Provides a MNS queue resource.

~> **NOTE:** Terraform will auto build a mns queue  while it uses `alicloud_mns_queue` to build a mns queue resource.

## Example Usage

Basic Usage

```
resource "alicloud_mns_queue" "queue"{
    name="${var.name}"
    delay_seconds=${var.delay_seconds}
    maximum_message_size=${var.maximum_message_size}
    message_retention_period=${var.message_retention_period}
    visibility_timeout=${var.visibility_timeout}
    polling_wait_seconds=${var.polling_wait_seconds}
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)Two queues on a single account in the same region cannot have the same name. A queue name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters.
* `delay_seconds` - (Optional)This attribute defines the length of time, in seconds, after which every message sent to the queue is dequeued. Valid value range: 0-604800 seconds, i.e., 0 to 7 days.
* `maximum_message_size` - (Optional)This indicates the maximum length, in bytes, of any message body sent to the queue. Valid value range: 1024-65536, i.e., 1K to 64K.
* `message_retention_period` - (Optional) Messages are deleted from the queue after a specified length of time, whether they have been activated or not. This attribute defines the viability period, in seconds, for every message in the queue. Valid value range: 60-259200 seconds, i.e., 1 minutes to 3 days.
* `visibility_timeout` - (Optional) The VisibilityTimeout attribute of the queue. A dequeued messages will change from active (visible) status to inactive (invisible) status, and this attribute defines the length of time, in seconds, that messages remain invisible. Messages return to active status after the set period. Valid value range: 1-43200 seconds, i.e., 1 seconds to 12 hours.
* `polling_wait_seconds` - (Optional) Long polling is measured in seconds. When this attribute is set to 0, long polling is disabled. When it is not set to 0, long polling is enabled and message dequeue requests will be processed only when valid messages are received or when long polling times out. Valid value range: 0-30 seconds.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the queue,is equal to name.
* `name` - The name of the queue.

## Import
MNS QUEUE can be imported using the id or name, e.g.

```
$ terraform import alicloud_mns_queue.queue queuename
```