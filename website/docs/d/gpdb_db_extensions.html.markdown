---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_db_extensions"
sidebar_current: "docs-alicloud-datasource-gpdb-db-extensions"
description: |-
  Provides a list of Gpdb Db Extension owned by an Alibaba Cloud account.
---

# alicloud_gpdb_db_extensions

This data source provides Gpdb Db Extension available to the user.[What is Db Extension](https://next.api.alibabacloud.com/document/gpdb/2016-05-03/CreateExtensions)

-> **NOTE:** Available since v1.291.0.

## Example Usage

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
  extension_name    = "uuid-ossp"
  db_instance_id    = alicloud_gpdb_instance.defaultqsmpIy.id
  database_name     = alicloud_gpdb_database.defaultPPmRVa.database_name
  is_latest_version = true
}

data "alicloud_gpdb_db_extensions" "default" {
  ids            = ["${alicloud_gpdb_db_extension.default.id}"]
  db_instance_id = alicloud_gpdb_instance.defaultqsmpIy.id
  database_name  = alicloud_gpdb_database.defaultPPmRVa.database_name
}

output "alicloud_gpdb_db_extension_example_id" {
  value = data.alicloud_gpdb_db_extensions.default.extensions.0.id
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required) The instance ID.

-> **NOTE:**   You can call the [DescribeDBInstances](https://www.alibabacloud.com/help/en/doc-detail/86911.html) operation to query the information about all AnalyticDB for PostgreSQL instances within a region, including instance IDs.

* `database_name` - (Required) The name of the database.
* `ids` - (Optional, Computed) A list of Db Extension IDs. The value is formulated as `<db_instance_id>:<database_name>:<extension_name>`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Db Extension IDs.
* `extensions` - A list of Db Extension Entries. Each element contains the following attributes:
  * `current_version` - **NOTE:** This field is only available when `enable_details` is `true`. Plug-in current version.
  * `description` - Plug-in description.
  * `extension_id` - **NOTE:** This field is only available when `enable_details` is `true`. Plug-in id.
  * `extension_name` - The name of the extension to install.
  * `is_install_need_restart` - **NOTE:** This field is only available when `enable_details` is `true`. Whether the instance needs to be restarted for installation.
  * `is_latest_version` - **NOTE:** This field is only available when `enable_details` is `true`. Whether the extension is at its latest version.
  * `latest_version` - **NOTE:** This field is only available when `enable_details` is `true`. Plug-in latest version.
  * `status` - The status of the extension.
  * `id` - The ID of the resource supplied above.
