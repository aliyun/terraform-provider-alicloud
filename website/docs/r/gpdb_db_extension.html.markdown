---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_db_extension"
description: |-
  Provides a Alicloud AnalyticDB for PostgreSQL (GPDB) Db Extension resource.
---

# alicloud_gpdb_db_extension

Provides a AnalyticDB for PostgreSQL (GPDB) Db Extension resource.



For information about AnalyticDB for PostgreSQL (GPDB) Db Extension and how to use it, see [What is Db Extension](https://next.api.alibabacloud.com/document/gpdb/2016-05-03/CreateExtensions).

-> **NOTE:** Available since v1.291.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "defaultiYyNGW" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultPlruct" {
  vpc_id     = alicloud_vpc.defaultiYyNGW.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultqsmpIy" {
  instance_spec         = "2C8G"
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  engine                = "gpdb"
  zone_id               = "cn-beijing-h"
  vswitch_id            = alicloud_vswitch.defaultPlruct.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.defaultiYyNGW.id
  db_instance_mode      = "StorageElastic"
}

resource "alicloud_gpdb_account" "defaultOwner" {
  account_name        = "tf_example"
  account_password    = "Example1234"
  account_description = "tf_example"
  db_instance_id      = alicloud_gpdb_instance.defaultqsmpIy.id
}

resource "alicloud_gpdb_database" "defaultPPmRVa" {
  owner              = alicloud_gpdb_account.defaultOwner.account_name
  database_name      = "seagull"
  db_instance_id     = alicloud_gpdb_instance.defaultqsmpIy.id
  character_set_name = "UTF8"
  collate            = "en_US.utf8"
  ctype              = "en_US.utf8"
}


resource "alicloud_gpdb_db_extension" "default" {
  extension_name       = "uuid-ossp"
  db_instance_id       = alicloud_gpdb_instance.defaultqsmpIy.id
  database_name        = alicloud_gpdb_database.defaultPPmRVa.database_name
  is_latest_version = true
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The instance ID.

-> **NOTE:**   You can call the [DescribeDBInstances](https://www.alibabacloud.com/help/en/doc-detail/86911.html) operation to query the information about all AnalyticDB for PostgreSQL instances within a region, including instance IDs.

* `database_name` - (Required, ForceNew) The name of the database.
* `extension_name` - (Required, ForceNew) The name of the extension to install. Each resource manages exactly one extension; multiple extension names are not supported.
* `is_latest_version` - (Optional) Whether the extension is at its latest version. Setting it to true calls UpgradeExtensions to upgrade the extension to the latest version. Leaving it unset or setting it to false triggers no upgrade and does not downgrade the extension. On read, this field is refreshed with the actual value returned by the server.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<db_instance_id>:<database_name>:<extension_name>`.
* `current_version` - Plug-in current version.
* `description` - Plug-in description.
* `extension_id` - Plug-in id.
* `is_install_need_restart` - Whether the instance needs to be restarted for installation.
* `latest_version` - Plug-in latest version.
* `status` - The status of the extension.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Db Extension.
* `delete` - (Defaults to 5 mins) Used when delete the Db Extension.
* `update` - (Defaults to 5 mins) Used when update the Db Extension.

## Import

AnalyticDB for PostgreSQL (GPDB) Db Extension can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_db_extension.example <db_instance_id>:<database_name>:<extension_name>
```