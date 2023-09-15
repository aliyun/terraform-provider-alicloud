---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_topic"
description: |-
  Provides a Alicloud Message Service Topic resource.
---

# alicloud_message_service_topic

Provides a Message Service Topic resource. 

For information about Message Service Topic and how to use it, see [What is Topic](https://www.alibabacloud.com/help/en/message-service/latest/createtopic).

-> **NOTE:** Available since v1.188.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_message_service_topic" "default" {
  max_message_size = 65536
  topic_name       = var.name
}
```

## Argument Reference

The following arguments are supported:
* `enable_logging` - (Optional) Specifies whether to enable the log management feature. Default value: false. Valid values:
  - `true`: enables the log management feature.
  - `false`: disables the log management feature.
* `max_message_size` - (Optional, Computed) The maximum size of a message body that can be sent to the topic. Unit: bytes. Valid values: 1024-65536. Default value: 65536.
* `topic_name` - (Required, ForceNew) Two topics on a single account in the same region cannot have the same name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.

The following arguments will be discarded. Please use new fields as soon as possible:
* `logging_enabled` - (Deprecated since v1.210.0). Field 'logging_enabled' has been deprecated from provider version 1.210.0. New field 'enable_logging' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Topic.
* `delete` - (Defaults to 5 mins) Used when delete the Topic.
* `update` - (Defaults to 5 mins) Used when update the Topic.

## Import

Message Service Topic can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_topic.example <id>
```