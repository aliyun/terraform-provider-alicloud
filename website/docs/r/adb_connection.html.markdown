---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_connection"
sidebar_current: "docs-alicloud-resource-adb-connection"
description: |-
  Provides an ADB cluster connection resource.
---

# alicloud_adb_connection

Provides an ADB connection resource to allocate an Internet connection string for ADB cluster.

-> **NOTE:** Each ADB instance will allocate a intranet connnection string automatically and its prifix is ADB instance ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.

-> **NOTE:** Available since v1.81.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_adb_zones" "default" {
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.4.0.0/24"
  zone_id      = data.alicloud_adb_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_adb_db_cluster" "default" {
  db_cluster_category = "Cluster"
  db_node_class       = "C8"
  db_node_count       = "4"
  db_node_storage     = "400"
  mode                = "reserver"
  db_cluster_version  = "3.0"
  payment_type        = "PayAsYouGo"
  vswitch_id          = alicloud_vswitch.default.id
  description         = var.name
  maintain_time       = "23:00Z-00:00Z"
  resource_group_id   = data.alicloud_resource_manager_resource_groups.default.ids.0
  security_ips        = ["10.168.1.12", "10.168.1.11"]
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_adb_connection" "default" {
  db_cluster_id     = alicloud_adb_db_cluster.default.id
  connection_prefix = "example"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `connection_prefix` - (Optional, ForceNew) Prefix of the cluster public endpoint. The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter. Default to `<db_cluster_id> + tf`.

## Attributes Reference

The following attributes are exported:

* `id` - The current cluster connection resource ID. Composed of cluster ID and connection string with format `<db_cluster_id>:<connection_prefix>`.
* `port` - Connection cluster port.
* `connection_string` - Connection cluster string.
* `ip_address` - The ip address of connection string.

## Import

ADB connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_adb_connection.example am-12345678
```
