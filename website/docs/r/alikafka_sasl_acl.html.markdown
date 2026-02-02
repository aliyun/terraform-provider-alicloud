---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_sasl_acl"
description: |-
  Provides a Alicloud Alikafka Sasl Acl resource.
---

# alicloud_alikafka_sasl_acl

Provides a Alikafka Sasl Acl resource.

Kafka access control.

For information about Alikafka Sasl Acl and how to use it, see [What is Sasl Acl](https://next.api.alibabacloud.com/document/alikafka/2019-09-16/CreateAcl).

-> **NOTE:** Available since v1.66.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alikafka_sasl_acl&exampleId=50f8e359-d4d7-2a6e-9d2c-b8675b0d795718c8d3b3&activeTab=example&spm=docs.r.alikafka_sasl_acl.0.50f8e359d4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alikafka_sasl_acl&spm=docs.r.alikafka_sasl_acl.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `acl_operation_type` - (Required, ForceNew) Operation type. Valid values:
  - `Write`: write
  - `Read`: read
  - `Describe`: read TransactionalId
  - `IdempotentWrite`: idempotent write to Cluster
  - `IDEMPOTENT_WRITE`: idempotent write to Cluster, only available for Serverless instances.
  - `DESCRIBE_CONFIGS`: query configuration, only available for Serverless instances.
* `acl_operation_types` - (Optional, Available since v1.270.0) Batch authorization operation types. Multiple operations are separated by commas (,). Valid values:
  - `Write`: write
  - `Read`: read
  - `Describe`: read TransactionalId
  - `IdempotentWrite`: idempotent write to Cluster
  - `IDEMPOTENT_WRITE`: idempotent write to Cluster, only available for Serverless instances.
  - `DESCRIBE_CONFIGS`: query configuration, only available for Serverless instances.
-> **NOTE:**  `acl_operation_types` is only supported for Serverless instances.
* `acl_permission_type` - (Optional, ForceNew, Available since v1.270.0) Authorization method. Value:
  - `DENY`: deny.
  - `ALLOW`: allow.
-> **NOTE:**  `acl_permission_type` is only supported for Serverless instances.
* `acl_resource_name` - (Required, ForceNew) The resource name.
  - The name of the resource, which can be a topic name, Group ID, cluster name, or transaction ID.
  - You can use an asterisk (*) to represent all resources of this type.
* `acl_resource_pattern_type` - (Required, ForceNew) Match the pattern. Valid values:
  - `LITERAL`: exact match
  - `PREFIXED`: prefix matching
* `acl_resource_type` - (Required, ForceNew) The resource type. Valid values:
  - `Topic`: the message Topic.
  - `Group`: consumer Group.
  - `Cluster`: the instance.
  - `TransactionalId`: transaction ID.
* `host` - (Optional, ForceNew) The host of the acl.
-> **NOTE:** From version 1.270.0, `host` can be set.
* `instance_id` - (Required, ForceNew) The instance ID.
* `username` - (Required, ForceNew) The user name.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<instance_id>:<username>:<acl_resource_type>:<acl_resource_name>:<acl_resource_pattern_type>:<acl_operation_type>`.

## Timeouts

-> **NOTE:** Available since v1.270.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Sasl Acl.
* `delete` - (Defaults to 5 mins) Used when delete the Sasl Acl.

## Import

Alikafka Sasl Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_alikafka_sasl_acl.example <instance_id>:<username>:<acl_resource_type>:<acl_resource_name>:<acl_resource_pattern_type>:<acl_operation_type>
```
