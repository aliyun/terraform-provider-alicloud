---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_account"
description: |-
  Provides a Alicloud Mongodb Account resource.
---

# alicloud_mongodb_account

Provides a Mongodb Account resource.



For information about Mongodb Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/en/doc-detail/62154.html).

-> **NOTE:** Available since v1.148.0.

## Example Usage

Basic Usage

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
* `account_description` - (Optional) Set the comment information of the account.
  - Cannot start with http:// or https.
  - Start with Chinese and English letters.
  - Can contain Chinese characters, English characters, underscores (_), dashes (-), and numbers, and can be 2 to 256 characters in length.
* `account_name` - (Required, ForceNew) Account Name
* `account_password` - (Required) Account Password
* `character_type` - (Optional, ForceNew, Computed, Available since v1.241.0) The account Comment Information type. Value:
  - `db`: shard account.
  - `normal`: The replica set account. (The default value).

-> **NOTE:** When setting the account type to normal, the AccountName must be root. This selection does not create a new account; instead, Terraform will update the password for the existing root account.

* `instance_id` - (Required, ForceNew) Instance Id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<account_name>`.
* `status` - Account Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account.
* `update` - (Defaults to 5 mins) Used when update the Account.

## Import

Mongodb Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_account.example <instance_id>:<account_name>
```