---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_sasl_acl"
sidebar_current: "docs-alicloud-resource-alikafka-sasl_acl"
description: |-
  Provides a Alicloud Alikafka Sasl Acl resource.
---

# alicloud_alikafka_sasl_acl

Provides an ALIKAFKA sasl acl resource, see [What is alikafka sasl acl](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/api-alikafka-2019-09-16-createacl).

-> **NOTE:** Available since v1.66.0.

-> **NOTE:**  Only the following regions support create alikafka sasl user.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alikafka_sasl_acl&exampleId=f5bd67e8-ea17-613c-1911-008df8d1de584edf1b34&activeTab=example&spm=docs.r.alikafka_sasl_acl.0.f5bd67e8ea&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_alikafka_instance" "default" {
  name            = "${var.name}-${random_integer.default.result}"
  partition_num   = 50
  disk_type       = "1"
  disk_size       = "500"
  deploy_type     = "5"
  io_max          = "20"
  spec_type       = "professional"
  service_version = "2.2.0"
  config          = "{\"enable.acl\":\"true\"}"
  vswitch_id      = alicloud_vswitch.default.id
  security_group  = alicloud_security_group.default.id
}

resource "alicloud_alikafka_topic" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  topic       = "example-topic"
  remark      = "topic-remark"
}

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  username    = var.name
  password    = "tf_example123"
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

```shell
$ terraform import alicloud_alikafka_sasl_acl.acl alikafka_post-cn-123455abc:username:Topic:test-topic:LITERAL:Write
```
