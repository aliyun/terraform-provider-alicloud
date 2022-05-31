---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_ons_topic"
sidebar_current: "docs-alicloud-resource-ons-topic"
description: |-
  Provides a Alicloud ONS Topic resource.
---

# alicloud\_ons\_topic

Provides an ONS topic resource.

For more information about how to use it, see [RocketMQ Topic Management API](https://www.alibabacloud.com/help/doc-detail/29591.html). 

-> **NOTE:** Available in 1.53.0+

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "onsInstanceName"
}

variable "topic" {
  default = "onsTopicName"
}

resource "alicloud_ons_instance" "default" {
  name   = var.name
  remark = "default_ons_instance_remark"
}

resource "alicloud_ons_topic" "default" {
  topic_name   = var.topic
  instance_id  = alicloud_ons_instance.default.id
  message_type = 0
  remark       = "dafault_ons_topic_remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ONS Instance that owns the topics.
* `topic` - (Optional, Deprecated in v1.97.0+) Replaced by `topic_name` after version 1.97.0.
* `topic_name` - (Optional, Available in v1.97.0+) Name of the topic. Two topics on a single instance cannot have the same name and the name cannot start with 'GID' or 'CID'. The length cannot exceed 64 characters.
* `message_type` - (Required) The type of the message. Read [Ons Topic Create](https://www.alibabacloud.com/help/doc-detail/29591.html) for further details.
* `remark` - (Optional) This attribute is a concise description of topic. The length cannot exceed 128.
* `perm` - (Deprecated) This attribute has been deprecated.
* `tags` - (Optional, Available in v1.97.0+) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

-> **NOTE:** At least one of `topic_name` and `topic` should be set.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<topic>`.

## Import

ONS TOPIC can be imported using the id, e.g.

```
$ terraform import alicloud_ons_topic.topic MQ_INST_1234567890_Baso1234567:onsTopicDemo
```
