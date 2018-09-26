--
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_topic"
sidebar_current: "docs-alicloud-resource-mns-topic"
description: |-
  Provides a Alicloud MNS Topic resource.
---

# alicloud\_mns\_topic

Provides a MNS topic resource.

~> **NOTE:** Terraform will auto build a mns topic  while it uses `alicloud_mns_topic` to build a mns topic resource.

## Example Usage

Basic Usage

```
resource "alicloud_mns_topic" "topic"{
    name="${var.name}"
    maximum_message_size=${var.maximum_message_size}
    logging_enabled=${var.logging_enabled}
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)Two topics on a single account in the same region cannot have the same name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters.
* `maximum_message_size` - (Optional)This indicates the maximum length, in bytes, of any message body sent to the topic. Valid value range: 1024-65536, i.e., 1K to 64K.
* `logging_enabled` - (Optional) is logging enabled? true or false default value is false

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the topic.
* `name` - The name of the topic.

## Import

MNS Topic can be imported using the id, e.g.

```
$ terraform import alicloud_mns_topic.topic topicName
```