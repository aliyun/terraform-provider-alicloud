---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_consumer_group"
sidebar_current: "docs-alicloud-resource-alikafka-consumer-group"
description: |-
  Provides a Alicloud Alikafka Consumer Group resource.
---

# alicloud_alikafka_consumer_group

Provides an ALIKAFKA consumer group resource, see [What is alikafka consumer group](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/api-alikafka-2019-09-16-createconsumergroup).

-> **NOTE:** Available since v1.56.0.

-> **NOTE:**  Only the following regions support create alikafka consumer group.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alikafka_consumer_group&exampleId=e0c01fec-c5b8-1e8d-7fd6-c01a9798ede2dfa2cf90&activeTab=example&spm=docs.r.alikafka_consumer_group.0.e0c01fecc5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

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

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_alikafka_instance" "default" {
  name           = "${var.name}-${random_integer.default.result}"
  partition_num  = "50"
  disk_type      = "1"
  disk_size      = "500"
  deploy_type    = "5"
  io_max         = "20"
  vswitch_id     = alicloud_vswitch.default.id
  security_group = alicloud_security_group.default.id
}

resource "alicloud_alikafka_consumer_group" "default" {
  consumer_id = var.name
  instance_id = alicloud_alikafka_instance.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alikafka_consumer_group&spm=docs.r.alikafka_consumer_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the ALIKAFKA Instance that owns the groups.
* `consumer_id` - (Required, ForceNew) ID of the consumer group. The length cannot exceed 64 characters.
* `tags` - (Optional, Available in v1.63.0+) A mapping of tags to assign to the resource.
* `description` - (Optional, ForceNew, Available in v1.157.0+) The description of the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<consumer_id>`.

## Import

ALIKAFKA GROUP can be imported using the id, e.g.

```shell
$ terraform import alicloud_alikafka_consumer_group.group alikafka_post-cn-123455abc:consumerId
```
