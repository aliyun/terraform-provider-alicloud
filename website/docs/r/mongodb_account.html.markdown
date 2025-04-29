---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_account"
description: |-
  Provides a Alicloud Mongo D B Account resource.
---

# alicloud_mongodb_account

Provides a Mongo D B Account resource.



For information about Mongo D B Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/en/doc-detail/62154.html).

-> **NOTE:** Available since v1.148.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mongodb_account&exampleId=6b410fdc-06af-0ce6-22f7-149630f8a0982c790f03&activeTab=example&spm=docs.r.mongodb_account.0.6b410fdc06&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}
data "alicloud_mongodb_zones" "default" {}
locals {
  index   = length(data.alicloud_mongodb_zones.default.zones) - 1
  zone_id = data.alicloud_mongodb_zones.default.zones[local.index].id
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = local.zone_id
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.2"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  vswitch_id          = alicloud_vswitch.default.id
  security_ip_list    = ["10.168.1.12", "100.69.7.112"]
  name                = var.name
  tags = {
    Created = "TF"
    For     = "example"
  }
}

resource "alicloud_mongodb_account" "default" {
  account_name        = "root"
  account_password    = "Example_123"
  instance_id         = alicloud_mongodb_instance.default.id
  account_description = var.name
}
```

### Deleting `alicloud_mongodb_account` or removing it from your configuration

Terraform cannot destroy resource `alicloud_mongodb_account`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `account_description` - (Optional) Account comment information.

-> **NOTE:**  Call the [ModifyAccountDescription](~~ 468391 ~~) interface to set the account description information before this parameter is returned.

* `account_name` - (Required, ForceNew) The new password.

  - The password must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include `! # $ % ^ & * ( ) _ + - =`
  - The password must be 8 to 32 characters in length.
* `account_password` - (Required) The password of the database account. The password must be 8 to 32 characters in length. It can contain at least three types of the following characters: uppercase letters, lowercase letters, digits, and special characters. Special characters include ! # $ % ^ & \* ( ) \_ + - =
* `character_type` - (Optional, ForceNew, Available since v1.241.0) The role type of the instance. Value description

  - When the instance type is sharded cluster, charactertype is required. The values are db and cs.
  - When the instance type is a replica set, charactertype can be null or pass in normal.
* `instance_id` - (Required, ForceNew) The account whose password needs to be reset. Set the value to `root`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<account_name>`.
* `status` - Account Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account.
* `update` - (Defaults to 5 mins) Used when update the Account.

## Import

Mongo D B Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_account.example <instance_id>:<account_name>
```