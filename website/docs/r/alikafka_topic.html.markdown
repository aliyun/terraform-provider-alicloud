---
subcategory: "Alikafka"
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
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_alikafka_instance" "default" {
  name        = "tf-testacc-alikafkainstance"
  topic_quota = "50"
  disk_type   = "1"
  disk_size   = "500"
  deploy_type = "5"
  io_max      = "20"
  vswitch_id  = alicloud_vswitch.default.id
}

variable "topic" {
  default = "alikafkaTopicName"
}

resource "alicloud_alikafka_topic" "default" {
  instance_id   = alicloud_alikafka_instance.default.id
  topic         = var.topic
  local_topic   = "false"
  compact_topic = "false"
  partition_num = "12"
  remark        = "dafault_kafka_topic_remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) InstanceId of your Kafka resource, the topic will create in this instance.
* `topic` - (Required, ForceNew) Name of the topic. Two topics on a single instance cannot have the same name. The length cannot exceed 64 characters.
* `local_topic` - (Optional, ForceNew) Whether the topic is localTopic or not.
* `compact_topic` - (Optional, ForceNew) Whether the topic is compactTopic or not. Compact topic must be a localTopic.
* `partition_num` - (Optional) The number of partitions of the topic. The number should between 1 and 48.
* `remark` - (Required) This attribute is a concise description of topic. The length cannot exceed 64.
* `tags` - (Optional, Available in v1.63.0+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<topic>`.

## Import

ALIKAFKA TOPIC can be imported using the id, e.g.

```shell
$ terraform import alicloud_alikafka_topic.topic alikafka_post-cn-123455abc:topicName
```

### Timeouts

-> **NOTE:** Available in v1.119.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the topic (until it reaches the initial `Running` status). 
