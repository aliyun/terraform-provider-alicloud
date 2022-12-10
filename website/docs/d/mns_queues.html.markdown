---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_queues"
sidebar_current: "docs-alicloud-datasource-mns-queues"
description: |-
  Provides a list of mns queues available to the user.
---

# alicloud\_mns\_queues

This data source provides a list of MNS queues in an Alibaba Cloud account according to the specified parameters.

-> **DEPRECATED:**  This datasource has been deprecated from version `1.188.0`. Please use new datasource [message_service_queues](https://www.terraform.io/docs/providers/alicloud/d/message_service_queues).

## Example Usage

```terraform
data "alicloud_mns_queues" "queues" {
  name_prefix = "tf-"
}

output "first_queue_id" {
  value = "${data.alicloud_mns_queues.queues.queues.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_prefix` - (Optional) A string to filter resulting queues by their name prefixs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of queue names. 
* `queues` - A list of queues. Each element contains the following attributes:
  * `id` - The id of the queue, The value is set to `name`.
  * `name` - The name of the queue
  * `delay_seconds` - This attribute defines the length of time, in seconds, after which every message sent to the queue is dequeued.
  * `maximum_message_size` - This indicates the maximum length, in bytes, of any message body sent to the queue.
  * `message_retention_period` - Messages are deleted from the queue after a specified length of time, whether they have been activated or not. This attribute defines the viability period, in seconds, for every message in the queue.
  * `visibility_timeouts` - Dequeued messages change from active (visible) status to inactive (invisible) status. This attribute defines the length of time, in seconds, that messages remain invisible. Messages return to active status after the set period.
  * `polling_wait_seconds` - Long polling is measured in seconds. When this attribute is set to 0, long polling is disabled. When it is not set to 0, long polling is enabled and message dequeue requests will be processed only when valid messages are received or when long polling times out.
