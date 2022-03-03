---
subcategory: "Alikafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_sasl_user"
sidebar_current: "docs-alicloud-resource-alikafka-sasl_user"
description: |-
  Provides a Alicloud Alikafka Sasl User resource.
---

# alicloud\_alikafka\_sasl\_user

Provides an Alikafka sasl user resource.

-> **NOTE:** Available in 1.66.0+

-> **NOTE:**  Only the following regions support create alikafka sasl user.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

For information about Alikafka sasl user and how to use it, see [What is Alikafka sasl user a](https://www.alibabacloud.com/help/en/doc-detail/162221.html)

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

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  username    = var.username
  password    = var.password
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the ALIKAFKA Instance that owns the groups.
* `username` - (Required, ForceNew) Username for the sasl user. The length should between 1 to 64 characters. The characters can only contain 'a'-'z', 'A'-'Z', '0'-'9', '_' and '-'.
* `password` - (Optional, Sensitive) Operation password. It may consist of letters, digits, or underlines, with a length of 1 to 64 characters. You have to specify one of `password` and `kms_encrypted_password` fields.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. You have to specify one of `password` and `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, MapString) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a user with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `type` - (Optional, ForceNew, Available in 1.159.0+) The authentication mechanism. Valid values: `plain`, `scram`. Default value: `plain`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value is formate as `<instance_id>:<username>`.

## Import

Alikafka Sasl User can be imported using the id, e.g.

```
terraform import alicloud_alikafka_sasl_user.example <instance_id>:<username>
```

