---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_connection"
sidebar_current: "docs-alicloud-resource-gpdb-connection"
description: |-
  Provides an AnalyticDB for PostgreSQL instance connection resource.
---

# alicloud_gpdb_connection

Provides a connection resource to allocate an Internet connection string for instance.

-> **NOTE:** Available since v1.48.0.

-> **NOTE:** Each instance will allocate a intranet connection string automatically and its prefix is instance ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_gpdb_zones" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_gpdb_zones.default.ids.0
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
  master_node_num       = 1
  payment_type          = "PayAsYouGo"
  private_ip_address    = "1.1.1.1"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = alicloud_vpc.default.id
  vswitch_id            = alicloud_vswitch.default.id
  ip_whitelist {
    security_ip_list = "127.0.0.1"
  }
}

resource "alicloud_gpdb_connection" "default" {
  instance_id       = alicloud_gpdb_instance.default.id
  connection_prefix = "exampelcon"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `connection_prefix` - (Optional, ForceNew) Prefix of an Internet connection string. It must be checked for uniqueness. It may consist of lowercase letters, numbers, and underlines, and must start with a letter and have no more than 30 characters. Default to <instance_id> + '-tf'.
* `port` - (Optional) Internet connection port. Valid value: [3200-3999]. Default to 3306.
* `connection_string` - (Optional) Connection instance string.
* `ip_address` - (Optional) The ip address of connection string.

## Timeouts

-> **NOTE:** Available since v1.53.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the Internet connection (until DB instance reaches the initial `Running` status). 
* `update` - (Defaults to 10 mins) Used when activating the DB instance during update.
* `delete` - (Defaults to 10 mins) Used when terminating the DB instance. 

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<instance_id>:<connection_prefix>`.

## Import

AnalyticDB for PostgreSQL's connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_connection.example abc12345678
```
