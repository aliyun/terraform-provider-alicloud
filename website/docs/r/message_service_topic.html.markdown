---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_topic"
sidebar_current: "docs-alicloud-resource-message-service-topic"
description: |-
  Provides a Alicloud Message Notification Service Topic resource.
---

# alicloud\_message\_service\_topic

Provides a Message Notification Service Topic resource.

For information about Message Notification Service Topic and how to use it, see [What is Topic](https://www.alibabacloud.com/help/en/message-service/latest/createtopic).

-> **NOTE:** Available in v1.188.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_message_service_topic" "default" {
  topic_name       = "tf-example-value"
  max_message_size = 12357
  logging_enabled  = true
}
```

## Argument Reference

The following arguments are supported:

* `topic_name` - (Required, ForceNew) Two topics on a single account in the same region cannot have the same name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.
* `max_message_size` - (Optional, Computed) The maximum size of a message body that can be sent to the topic. Unit: bytes. Valid values: 1024-65536. Default value: 65536.
* `logging_enabled` - (Optional) Specifies whether to enable the log management feature. Default value: false. Valid values:
  - `true`: enables the log management feature.
  - `false`: disables the log management feature.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Topic. Its value is same as `topic_name`.

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Topic.
* `update` - (Defaults to 3 mins) Used when update the Topic.
* `delete` - (Defaults to 3 mins) Used when delete the Topic.

## Import

Message Notification Service Topic can be imported using the id or topic_name, e.g.

```shell
$ terraform import alicloud_message_service_topic.example <topic_name>
```
