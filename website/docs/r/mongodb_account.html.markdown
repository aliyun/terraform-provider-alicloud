---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_account"
sidebar_current: "docs-alicloud-resource-mongodb-account"
description: |-
  Provides a Alicloud MongoDB Account resource.
---

# alicloud\_mongodb\_account

Provides a MongoDB Account resource.

For information about MongoDB Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/en/doc-detail/62154.html).

-> **NOTE:** Available in v1.148.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_name = "subnet-for-local-test"
}
resource "alicloud_mongodb_instance" "default" {
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  name                = var.name
  vswitch_id          = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_mongodb_account" "example" {
  account_name        = "root"
  account_password    = "example_value"
  instance_id         = alicloud_mongodb_instance.default.id
  account_description = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `account_description` - (Optional) The description of the account.
  * The description must start with a letter, and cannot start with `http://` or `https://`.
  * It must be `2` to `256` characters in length, and can contain letters, digits, underscores (_), and hyphens (-).
* `account_name` - (Required) The name of the account. Valid values: `root`.
* `account_password` - (Required, Sensitive) The Password of the Account.
  * The password must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include `!#$%^&*()_+-=`.
  * The password must be `8` to `32` characters in length.
* `instance_id` - (Required, ForceNew) The ID of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Account. The value formats as `<instance_id>:<account_name>`.
* `status` - The status of the account. Valid values: `Unavailable`, `Available`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Account.
* `update` - (Defaults to 10 mins) Used when update the Account.

## Import

MongoDB Account can be imported using the id, e.g.

```
$ terraform import alicloud_mongodb_account.example <instance_id>:<account_name>
```