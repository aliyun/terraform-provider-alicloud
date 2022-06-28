---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_account"
sidebar_current: "docs-alicloud-resource-gpdb-account"
description: |-
  Provides a Alicloud GPDB Account resource.
---

# alicloud\_gpdb\_account

Provides a GPDB Account resource.

For information about GPDB Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/doc-detail/86924.htm).

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tftestacc"
}
data "alicloud_gpdb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.zones.2.id
}

resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_gpdb_zones.default.zones.3.id
  vswitch_name = var.name
}

resource "alicloud_gpdb_elastic_instance" "default" {
  engine                  = "gpdb"
  engine_version          = "6.0"
  seg_storage_type        = "cloud_essd"
  seg_node_num            = 4
  storage_size            = 50
  instance_spec           = "2C16G"
  db_instance_description = "Created by terraform"
  instance_network_type   = "VPC"
  payment_type            = "PayAsYouGo"
  vswitch_id              = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.default.*.id, [""])[0]
}

resource "alicloud_gpdb_account" "default" {
  account_name        = var.name
  db_instance_id      = alicloud_gpdb_elastic_instance.default.id
  account_password    = "TFTest123"
  account_description = var.name
}

```

## Argument Reference

The following arguments are supported:

* `account_description` - (Optional, ForceNew) The description of the account.
  * Starts with a letter.
  * Does not start with `http://` or `https://`.
  * Contains letters, underscores (_), hyphens (-), or digits.
  * Be 2 to 256 characters in length.
* `account_name` - (Required, ForceNew) The name of the account. The account name must be unique and meet the following requirements:
  * Starts with a letter.
  * Contains only lowercase letters, digits, or underscores (_).
  * Be up to 16 characters in length.
  * Contains no reserved keywords.
* `account_password` - (Required) The password of the account. The password must be 8 to 32 characters in length and contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include `! @ # $ % ^ & * ( ) _ + - =`.
* `db_instance_id` - (Required, ForceNew) The ID of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Account. The value formats as `<db_instance_id>:<account_name>`.
* `status` - The status of the account. Valid values: `Active`, `Creating` and `Deleting`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Account.

## Import

GPDB Account can be imported using the id, e.g.

```
$ terraform import alicloud_gpdb_account.example <db_instance_id>:<account_name>
```
