---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_remote_adb_data_source"
description: |-
  Provides a Alicloud GPDB Remote ADB Data Source resource.
---

# alicloud_gpdb_remote_adb_data_source

Provides a GPDB Remote ADB Data Source resource.

RemoteADBDataSource is the data external table call method between greenplums, which will be used for data external table access between ADB-PG.

For information about GPDB Remote ADB Data Source and how to use it, see [What is Remote ADB Data Source](https://www.alibabacloud.com/help/en/analyticdb/analyticdb-for-postgresql/developer-reference/api-gpdb-2016-05-03-createremoteadbdatasource).

-> **NOTE:** Available since v1.227.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_remote_adb_data_source&exampleId=29ccb4f3-21f7-217c-d64b-6ccddea267def0bb2562&activeTab=example&spm=docs.r.gpdb_remote_adb_data_source.0.29ccb4f321&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default4Mf0nY" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultwSAVpf" {
  vpc_id     = alicloud_vpc.default4Mf0nY.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultEtEzMF" {
  instance_spec         = "2C8G"
  description           = var.name
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  zone_id               = "cn-beijing-h"
  vswitch_id            = alicloud_vswitch.defaultwSAVpf.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.default4Mf0nY.id
  db_instance_mode      = "StorageElastic"
  engine                = "gpdb"
}

resource "alicloud_gpdb_instance" "defaultEY7t9t" {
  instance_spec         = "2C8G"
  description           = var.name
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  zone_id               = "cn-beijing-h"
  vswitch_id            = alicloud_vswitch.defaultwSAVpf.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.default4Mf0nY.id
  db_instance_mode      = "StorageElastic"
  engine                = "gpdb"
}

resource "alicloud_gpdb_account" "default26qpEo" {
  account_description = "example_001"
  db_instance_id      = alicloud_gpdb_instance.defaultEtEzMF.id
  account_name        = "example_001"
  account_password    = "example_001"
}

resource "alicloud_gpdb_account" "defaultwXePof" {
  account_description = "example_001"
  db_instance_id      = alicloud_gpdb_instance.defaultEY7t9t.id
  account_name        = "example_001"
  account_password    = "example_001"
}


resource "alicloud_gpdb_remote_adb_data_source" "default" {
  remote_database       = "example_001"
  manager_user_name     = "example_001"
  user_name             = "example_001"
  remote_db_instance_id = alicloud_gpdb_account.defaultwXePof.db_instance_id
  local_database        = "example_001"
  data_source_name      = "myexample"
  user_password         = "example_001"
  manager_user_password = "example_001"
  local_db_instance_id  = alicloud_gpdb_instance.defaultEtEzMF.id
}
```

## Argument Reference

The following arguments are supported:
* `data_source_name` - (Optional) Data Source Name
* `local_database` - (Required, ForceNew) The database of the local instance which connection data.
* `local_db_instance_id` - (Required, ForceNew) The instanceId of the local instance which connection data.
* `manager_user_name` - (Required, ForceNew) The Management user name of the local instance.
* `manager_user_password` - (Required, ForceNew) Password of the Manager user of the local instance
* `remote_database` - (Required, ForceNew) The database of the remote instance which provide data.
* `remote_db_instance_id` - (Required, ForceNew) The instanceId of the remote instance which provide data.
* `user_name` - (Required) The user name used to connect to the remote instance
* `user_password` - (Required) The user password used to connect to the remote instance

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<local_db_instance_id>:<remote_adb_data_source_id>`.
* `remote_adb_data_source_id` - The first ID of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Remote ADB Data Source.
* `delete` - (Defaults to 5 mins) Used when delete the Remote ADB Data Source.
* `update` - (Defaults to 5 mins) Used when update the Remote ADB Data Source.

## Import

GPDB Remote ADB Data Source can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_remote_adb_data_source.example <local_db_instance_id>:<remote_adb_data_source_id>
```