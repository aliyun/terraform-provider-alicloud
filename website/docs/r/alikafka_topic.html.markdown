---
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_topic"
sidebar_current: "docs-alicloud-resource-alikafka-topic"
description: |-
  Provides a Alicloud ALIKAFKA Topic resource.
---

# alicloud\_alikafka\_topic

Provides an ALIKAFKA topic resource.

-> **NOTE:** Available in 1.56.0+

-> **NOTE:**  Only the following regions support create alikafka topic.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`ap-southeast-1`,`ap-south-1`,`ap-southeast-5`]

## Example Usage

Basic Usage

```
variable "instance_id" {
  default = "yourInstanceId"
}

variable "topic" {
  default = "alikafkaTopicName"
}

resource "alicloud_alikafka_topic" "default" {
  instance_id = "${var.instance_id}"
  topic = "${var.topic}"
  local_topic = "false"
  compact_topic = "false"
  partition_num = "12"
  remark = "dafault_ons_topic_remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) InstanceId of your Kafka resource, the topic will create in this instance.
* `topic` - (Required, ForceNew) Name of the topic. Two topics on a single instance cannot have the same name. The length cannot exceed 64 characters.
* `local_topic` - (Optional, ForceNew) Whether the topic is localTopic or not.
* `compact_topic` - (Optional, ForceNew) Whether the topic is compactTopic or not. Compact topic must be a localTopic.
* `partition_num` - (Optional, ForceNew) The number of partitions of the topic. The number should between 1 and 48.
* `remark` - (Required, ForceNew) This attribute is a concise description of topic. The length cannot exceed 64.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<topic>`.

## Import

ALIKAFKA TOPIC can be imported using the id, e.g.

```
$ terraform import alicloud_alikafka_topic.topic KAFKA_INST_1234567890_Baso1234567:alikafkaTopicDemo
```