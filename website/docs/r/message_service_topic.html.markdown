---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_topic"
sidebar_current: "docs-alicloud-resource-message-service-topic"
description: |-
  Provides a Alicloud Message Notification Service Topic resource.
---

# alicloud_message_service_topic

Provides a Message Notification Service Topic resource.

For information about Message Notification Service Topic and how to use it, see [What is Topic](https://www.alibabacloud.com/help/en/message-service/latest/createtopic).

-> **NOTE:** Available since v1.188.0.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_message_service_topic&exampleId=0fc6d852-0b19-a125-e957-d54e6fafbd179f895ece&activeTab=example&spm=docs.r.message_service_topic.0.0fc6d8520b" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_message_service_topic" "default" {
  topic_name       = var.name
  max_message_size = 12357
  logging_enabled  = true
}
```

## Argument Reference

The following arguments are supported:

* `topic_name` - (Required, ForceNew) Two topics on a single account in the same region cannot have the same name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.
* `max_message_size` - (Optional) The maximum size of a message body that can be sent to the topic. Unit: bytes. Valid values: 1024-65536. Default value: 65536.
* `logging_enabled` - (Optional) Specifies whether to enable the log management feature. Default value: false. Valid values:
  - `true`: enables the log management feature.
  - `false`: disables the log management feature.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Topic. Its value is same as `topic_name`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Topic.
* `update` - (Defaults to 3 mins) Used when update the Topic.
* `delete` - (Defaults to 3 mins) Used when delete the Topic.

## Import

Message Notification Service Topic can be imported using the id or topic_name, e.g.

```shell
$ terraform import alicloud_message_service_topic.example <topic_name>
```
