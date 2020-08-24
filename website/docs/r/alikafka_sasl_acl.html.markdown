---
subcategory: "Alikafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_sasl_acl"
sidebar_current: "docs-alicloud-resource-alikafka-sasl_acl"
description: |-
  Provides a Alicloud Alikafka Sasl Acl resource.
---

# alicloud\_alikafka\_sasl\_acl

Provides an ALIKAFKA sasl acl resource.

-> **NOTE:** Available in 1.66.0+

-> **NOTE:**  Only the following regions support create alikafka sasl user.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`ap-southeast-1`,`ap-south-1`,`ap-southeast-5`]

## Example Usage

Basic Usage

```
variable "username" {
  default = "testusername"
}

variable "password" {
  default = "testpassword"
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
  availability_zone = data.alicloud_zones.default.zones[0].id
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

resource "alicloud_alikafka_topic" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  topic       = "test-topic"
  remark      = "topic-remark"
}

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  username    = var.username
  password    = var.password
}

resource "alicloud_alikafka_sasl_acl" "default" {
  instance_id               = alicloud_alikafka_instance.default.id
  username                  = alicloud_alikafka_sasl_user.default.username
  acl_resource_type         = "Topic"
  acl_resource_name         = alicloud_alikafka_topic.default.topic
  acl_resource_pattern_type = "LITERAL"
  acl_operation_type        = "Write"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the ALIKAFKA Instance that owns the groups.
* `username` - (Required, ForceNew) Username for the sasl user. The length should between 1 to 64 characters. The user should be an existed sasl user.
* `acl_resource_type` - (Required, ForceNew) Resource type for this acl. The resource type can only be "Topic" and "Group".
* `acl_resource_name` - (Required, ForceNew) Resource name for this acl. The resource name should be a topic or consumer group name.
* `acl_resource_pattern_type` - (Required, ForceNew) Resource pattern type for this acl. The resource pattern support two types "LITERAL" and "PREFIXED". "LITERAL": A literal name defines the full name of a resource. The special wildcard character "*" can be used to represent a resource with any name. "PREFIXED": A prefixed name defines a prefix for a resource.
* `acl_operation_type` - (Required, ForceNew) Operation type for this acl. The operation type can only be "Write" and "Read".

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<username>:<acl_resource_type>:<acl_resource_name>:<acl_resource_pattern_type>:<acl_operation_type>`.
* `host` - The host of the acl.

## Import

ALIKAFKA GROUP can be imported using the id, e.g.

```
$ terraform import alicloud_alikafka_sasl_acl.acl alikafka_post-cn-123455abc:username:Topic:test-topic:LITERAL:Write
```
