---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_sasl_user"
sidebar_current: "docs-alicloud-resource-alikafka-sasl-user"
description: |-
  Provides a Alicloud Alikafka Sasl User resource.
---

# alicloud_alikafka_sasl_user

Provides an Alikafka Sasl User resource.

-> **NOTE:** Available since v1.66.0.

-> **NOTE:**  Only the following regions support create alikafka Sasl User.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

For information about Alikafka Sasl User and how to use it, see [What is Sasl User](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/api-alikafka-2019-09-16-createsasluser).

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
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
  vswitch_id      = alicloud_vswitch.default.id
  security_group  = alicloud_security_group.default.id
  config          = <<EOF
  {
    "enable.acl": "true"
  }
  EOF
}

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  username    = var.name
  password    = "tf_example123"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the ALIKAFKA Instance that owns the groups.
* `username` - (Required, ForceNew) The name of the SASL user. The length should between `1` to `64` characters. The characters can only contain `a`-`z`, `A`-`Z`, `0`-`9`, `_` and `-`.
* `type` - (Optional, ForceNew, Available since v1.159.0) The authentication mechanism. Default value: `plain`. Valid values: `plain`, `scram`.
* `password` - (Optional, Sensitive) The password of the SASL user. It may consist of letters, digits, or underlines, with a length of 1 to 64 characters. You have to specify one of `password` and `kms_encrypted_password` fields.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. You have to specify one of `password` and `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, MapString) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a user with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Sasl User. It formats as `<instance_id>:<username>`.

## Import

Alikafka Sasl User can be imported using the id, e.g.

```shell
terraform import alicloud_alikafka_sasl_user.example <instance_id>:<username>
```
