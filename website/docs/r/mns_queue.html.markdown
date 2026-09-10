---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_queue"
sidebar_current: "docs-alicloud-resource-mns-queue"
description: |-
  Provides a Alicloud MNS Queue resource.
---

# alicloud\_mns\_queue

Provides a MNS queue resource.

-> **NOTE:** Terraform will auto build a mns queue  while it uses `alicloud_mns_queue` to build a mns queue resource.

-> **DEPRECATED:**  This resource has been deprecated from version `1.188.0`. Please use new resource [message_service_queue](https://www.terraform.io/docs/providers/alicloud/r/message_service_queue).

## Example Usage

Basic Usage

```terraform
resource "alicloud_mns_queue" "queue" {
  name                     = "tf-example-mnsqueue"
  delay_seconds            = 0
  maximum_message_size     = 65536
  message_retention_period = 345600
  visibility_timeout       = 30
  polling_wait_seconds     = 0
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_mns_queue&spm=docs.r.mns_queue.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForcesNew)Two queues on a single account in the same region cannot have the same name. A queue name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters .
* `delay_seconds` - (Optional)This attribute defines the length of time, in seconds, after which every message sent to the queue is dequeued. Valid value range: 0-604800 seconds, i.e., 0 to 7 days. Default value to 0.
* `maximum_message_size` - (Optional)This indicates the maximum length, in bytes, of any message body sent to the queue. Valid value range: 1024-65536, i.e., 1K to 64K. Default value to 65536.
* `message_retention_period` - (Optional) Messages are deleted from the queue after a specified length of time, whether they have been activated or not. This attribute defines the viability period, in seconds, for every message in the queue. Valid value range: 60-604800 seconds, i.e., 1 minutes to 7 days. Default value to 345600.
* `visibility_timeout` - (Optional) The VisibilityTimeout attribute of the queue. A dequeued messages will change from active (visible) status to inactive (invisible) status, and this attribute defines the length of time, in seconds, that messages remain invisible. Messages return to active status after the set period. Valid value range: 1-43200 seconds, i.e., 1 seconds to 12 hours. Default value to 30.
* `polling_wait_seconds` - (Optional) Long polling is measured in seconds. When this attribute is set to 0, long polling is disabled. When it is not set to 0, long polling is enabled and message dequeue requests will be processed only when valid messages are received or when long polling times out. Valid value range: 0-30 seconds. Default value to 0.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the queue is equal to name.

## Import

MNS QUEUE can be imported using the id or name, e.g.

```shell
$ terraform import alicloud_mns_queue.queue queuename
```
