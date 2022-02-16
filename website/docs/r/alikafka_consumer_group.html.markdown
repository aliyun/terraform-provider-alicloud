---
subcategory: "Alikafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_consumer_group"
sidebar_current: "docs-alicloud-resource-alikafka-consumer-group"
description: |-
  Provides a Alicloud Alikafka Consumer Group resource.
---

# alicloud\_alikafka\_consumer\_group

Provides an ALIKAFKA consumer group resource.

-> **NOTE:** Available in 1.56.0+

-> **NOTE:**  Only the following regions support create alikafka consumer group.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

## Example Usage

Basic Usage

```
variable "consumer_id" {
  default = "CID-alikafkaGroupDatasourceName"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
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

resource "alicloud_alikafka_consumer_group" "default" {
  consumer_id = var.consumer_id
  instance_id = alicloud_alikafka_instance.default.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the ALIKAFKA Instance that owns the groups.
* `consumer_id` - (Required, ForceNew) ID of the consumer group. The length cannot exceed 64 characters.
* `tags` - (Optional, Available in v1.63.0+) A mapping of tags to assign to the resource.
* `description` - (Optional, ForceNew Available in v1.157.0+) The description of the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<consumer_id>`.

## Import

ALIKAFKA GROUP can be imported using the id, e.g.

```
$ terraform import alicloud_alikafka_consumer_group.group alikafka_post-cn-123455abc:consumerId
```
