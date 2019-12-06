---
subcategory: "Message Notification Service (MNS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_topic"
sidebar_current: "docs-alicloud-resource-mns-topic"
description: |-
  Provides a Alicloud MNS Topic resource.
---

# alicloud\_mns\_topic

Provides a MNS topic resource.

-> **NOTE:** Terraform will auto build a mns topic  while it uses `alicloud_mns_topic` to build a mns topic resource.

## Example Usage

Basic Usage

```
resource "alicloud_mns_topic" "topic" {
  name                 = "tf-example-mnstopic"
  maximum_message_size = 65536
  logging_enabled      = false
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew)Two topics on a single account in the same region cannot have the same name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters.
* `maximum_message_size` - (Optional)This indicates the maximum length, in bytes, of any message body sent to the topic. Valid value range: 1024-65536, i.e., 1K to 64K. Default value to 65536.
* `logging_enabled` - (Optional) Is logging enabled? true or false. Default value to false.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the topic is equal to name.

## Import

MNS Topic can be imported using the id or name, e.g.

```
$ terraform import alicloud_mns_topic.topic topicName

```
