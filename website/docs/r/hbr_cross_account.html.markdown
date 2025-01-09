---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_cross_account"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Cross Account resource.
---

# alicloud_hbr_cross_account

Provides a Hybrid Backup Recovery (HBR) Cross Account resource.

The cross account is used for the cross-account backup in the Cloud Backup. The management account can back up the resources under the cross account.

For information about Hybrid Backup Recovery (HBR) Cross Account and how to use it, see [What is Cross Account](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-guangzhou"
}


resource "alicloud_hbr_cross_account" "default" {
  cross_account_user_id   = "1"
  cross_account_role_name = var.name
  alias                   = var.name
}
```

## Argument Reference

The following arguments are supported:
* `alias` - (Optional, ForceNew) Backup account alias
* `cross_account_role_name` - (Required, ForceNew) The name of RAM role that the backup account authorizes the management account to manage its resources
* `cross_account_user_id` - (Required, ForceNew, Int) The uid of the backup account.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<cross_account_user_id>:<cross_account_role_name>`.
* `create_time` - Timestamp of the creation time

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cross Account.
* `delete` - (Defaults to 5 mins) Used when delete the Cross Account.

## Import

Hybrid Backup Recovery (HBR) Cross Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_cross_account.example <cross_account_user_id>:<cross_account_role_name>
```