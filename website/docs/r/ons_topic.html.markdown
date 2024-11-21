---
subcategory: "RocketMQ (Ons)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ons_topic"
sidebar_current: "docs-alicloud-resource-ons-topic"
description: |-
  Provides a Alicloud ONS Topic resource.
---

# alicloud_ons_topic

Provides an ONS topic resource.

For more information about how to use it, see [RocketMQ Topic Management API](https://www.alibabacloud.com/help/doc-detail/29591.html). 

-> **NOTE:** Available in 1.53.0+

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ons_topic&exampleId=51c8c410-42ff-d122-f626-491fa574ed62726e2f4c&activeTab=example&spm=docs.r.ons_topic.0.51c8c41042&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "onsInstanceName"
}

variable "topic" {
  default = "onsTopicName"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}


resource "alicloud_ons_instance" "default" {
  instance_name = "${var.name}-${random_integer.default.result}"
  remark        = "default_ons_instance_remark"
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

* `instance_id` - (Required, ForceNew) ID of the ONS Instance that owns the topics.
* `topic` - (Optional, Deprecated from v1.97.0+) Replaced by `topic_name` after version 1.97.0.
* `topic_name` - (Optional, ForceNew, Available in v1.97.0+) Name of the topic. Two topics on a single instance cannot have the same name and the name cannot start with 'GID' or 'CID'. The length cannot exceed 64 characters.
* `message_type` - (Required, ForceNew) The type of the message. Read [Ons Topic Create](https://www.alibabacloud.com/help/doc-detail/29591.html) for further details.
* `remark` - (Optional, ForceNew) This attribute is a concise description of topic. The length cannot exceed 128.
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

```shell
$ terraform import alicloud_ons_topic.topic MQ_INST_1234567890_Baso1234567:onsTopicDemo
```
