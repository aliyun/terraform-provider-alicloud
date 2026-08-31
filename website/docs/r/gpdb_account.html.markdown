---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_account"
description: |-
  Provides a Alicloud GPDB Account resource.
---

# alicloud_gpdb_account

Provides a GPDB Account resource.

For information about GPDB Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/doc-detail/86924.htm).

-> **NOTE:** Available since v1.142.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_account&exampleId=32410332-99b6-368a-bfc8-d78a60558e39415edfab&activeTab=example&spm=docs.r.gpdb_account.0.3241033299&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_gpdb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}

resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  payment_type          = "PayAsYouGo"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = data.alicloud_vswitches.default.ids[0]
  ip_whitelist {
    security_ip_list = "127.0.0.1"
  }
}

resource "alicloud_gpdb_account" "default" {
  account_name        = "tf_example"
  db_instance_id      = alicloud_gpdb_instance.default.id
  account_password    = "Example1234"
  account_description = "tf_example"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_gpdb_account&spm=docs.r.gpdb_account.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `account_description` - (Optional) The description of the account.
* `account_name` - (Required, ForceNew) The account name.
* `account_password` - (Required) AccountPassword
* `account_type` - (Optional, ForceNew, Available since v1.230.0) Account type. The value range is as follows:

  Normal: Normal account number.

  Super: a high-privilege account.
* `db_instance_id` - (Required, ForceNew) The Adb pg instance ID.
* `database_name` - (Optional, Available since v1.230.0) Database name, with the following restrictions:
  - Can only contain letters, numbers and underscores.
  - Must start with a letter.
  - Length cannot exceed 63 characters.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<account_name>`.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account.
* `delete` - (Defaults to 5 mins) Used when delete the Account.
* `update` - (Defaults to 5 mins) Used when update the Account.

## Import

GPDB Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_account.example <db_instance_id>:<account_name>
```