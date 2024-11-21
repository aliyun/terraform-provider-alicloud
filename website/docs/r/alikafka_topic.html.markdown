---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_topic"
sidebar_current: "docs-alicloud-resource-alikafka-topic"
description: |-
  Provides a Alicloud ALIKAFKA Topic resource.
---

# alicloud_alikafka_topic

Provides an ALIKAFKA topic resource, see [What is Alikafka topic ](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/api-alikafka-2019-09-16-createtopic).

-> **NOTE:** Available since v1.56.0.

-> **NOTE:**  Only the following regions support create alikafka topic.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alikafka_topic&exampleId=5d022a37-694a-dbde-15f1-c57d901ec512d1e4c7c8&activeTab=example&spm=docs.r.alikafka_topic.0.5d022a3769&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "instance_name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_alikafka_instance" "default" {
  name           = "${var.instance_name}-${random_integer.default.result}"
  partition_num  = "50"
  disk_type      = "1"
  disk_size      = "500"
  deploy_type    = "5"
  io_max         = "20"
  vswitch_id     = alicloud_vswitch.default.id
  security_group = alicloud_security_group.default.id
}

resource "alicloud_alikafka_topic" "default" {
  instance_id   = alicloud_alikafka_instance.default.id
  topic         = "example-topic"
  local_topic   = "false"
  compact_topic = "false"
  partition_num = "12"
  remark        = "dafault_kafka_topic_remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) InstanceId of your Kafka resource, the topic will create in this instance.
* `topic` - (Required, ForceNew) Name of the topic. Two topics on a single instance cannot have the same name. The length cannot exceed 249 characters.
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

## Timeouts

-> **NOTE:** Available since v1.119.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the topic (until it reaches the initial `Running` status). 
