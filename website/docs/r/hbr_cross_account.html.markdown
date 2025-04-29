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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_hbr_cross_account&exampleId=2d9de846-a7bd-150a-b6c2-a4436942ce650d61451a&activeTab=example&spm=docs.r.hbr_cross_account.0.2d9de846a7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cross Account.
* `delete` - (Defaults to 5 mins) Used when delete the Cross Account.

## Import

Hybrid Backup Recovery (HBR) Cross Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_cross_account.example <cross_account_user_id>:<cross_account_role_name>
```