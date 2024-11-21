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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_adb_connection&exampleId=28d0f9f6-fdbe-710b-3d91-1e6ffe52d36cbdcfa06e&activeTab=example&spm=docs.r.adb_connection.0.28d0f9f6fd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_adb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_adb_zones.default.ids.0
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_adb_db_cluster" "cluster" {
  db_cluster_category = "MixedStorage"
  mode                = "flexible"
  compute_resource    = "8Core32GB"
  vswitch_id          = local.vswitch_id
  description         = var.name
}

resource "alicloud_adb_connection" "default" {
  db_cluster_id     = alicloud_adb_db_cluster.cluster.id
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
