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
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_alikafka_instance" "default" {
  name           = var.name
  partition_num  = "50"
  disk_type      = "1"
  disk_size      = "500"
  deploy_type    = "5"
  io_max         = "20"
  vswitch_id     = data.alicloud_vswitches.default.ids.0
  security_group = alicloud_security_group.default.id
}

resource "alicloud_alikafka_consumer_group" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  consumer_id = "${var.name}-${random_integer.default.result}"
  description = "${var.name}-${random_integer.default.result}"
}
```

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
