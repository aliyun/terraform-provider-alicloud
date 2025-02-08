---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_db_resource_group"
description: |-
  Provides a Alicloud AnalyticDB for PostgreSQL (GPDB) Db Resource Group resource.
---

# alicloud_gpdb_db_resource_group

Provides a AnalyticDB for PostgreSQL (GPDB) Db Resource Group resource.



For information about AnalyticDB for PostgreSQL (GPDB) Db Resource Group and how to use it, see [What is Db Resource Group](https://next.api.alibabacloud.com/document/gpdb/2016-05-03/CreateDBResourceGroup).

-> **NOTE:** Available since v1.225.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "eu-central-1"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultZc8RD9" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultRv5UXt" {
  vpc_id     = alicloud_vpc.defaultZc8RD9.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultJXWSlW" {
  instance_spec         = "2C8G"
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  engine                = "gpdb"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  zone_id               = data.alicloud_zones.default.zones.0.id
  vswitch_id            = alicloud_vswitch.defaultRv5UXt.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.defaultZc8RD9.id
  db_instance_mode      = "StorageElastic"
  description           = var.name
}


resource "alicloud_gpdb_db_resource_group" "default" {
  resource_group_config = jsonencode({ "CpuRateLimit" : 10, "MemoryLimit" : 10, "MemorySharedQuota" : 80, "MemorySpillRatio" : 0, "Concurrency" : 10 })
  db_instance_id        = alicloud_gpdb_instance.defaultJXWSlW.id
  resource_group_name   = "yb_example_group"
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The instance ID.> You can call the [DescribeDBInstances](~~ 86911 ~~) operation to view the instance IDs of all AnalyticDB PostgreSQL instances in the target region.
* `resource_group_config` - (Required) Resource group configuration.
* `resource_group_name` - (Required, ForceNew) Resource group name.
* `role_list` - (Optional, Set, Available since v1.230.0) Role List

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<resource_group_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Db Resource Group.
* `delete` - (Defaults to 5 mins) Used when delete the Db Resource Group.
* `update` - (Defaults to 5 mins) Used when update the Db Resource Group.

## Import

AnalyticDB for PostgreSQL (GPDB) Db Resource Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_db_resource_group.example <db_instance_id>:<resource_group_name>
```