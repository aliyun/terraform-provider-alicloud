---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_db_cluster"
sidebar_current: "docs-alicloud-resource-click-house-db-cluster"
description: |-
  Provides a Alicloud Click House DBCluster resource.
---

# alicloud_click_house_db_cluster

Provides a Click House DBCluster resource.

For information about Click House DBCluster and how to use it, see [What is DBCluster](https://www.alibabacloud.com/help/zh/clickhouse/latest/api-clickhouse-2019-11-11-createdbinstance).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  region = var.region
}

data "alicloud_click_house_regions" "default" {
  region_id = var.region
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
}

resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "23.8"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "100"
  storage_type            = "cloud_essd"
  vswitch_id              = alicloud_vswitch.default.id
  vpc_id                  = alicloud_vpc.default.id
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Required, ForceNew) The Category of DBCluster. Valid values: `Basic`,`HighAvailability`.
* `db_cluster_class` - (Required, ForceNew) The DBCluster class. According to the category, db_cluster_class has two value ranges:
  * Under the condition that the category is the `Basic`, Valid values: `LS20`, `LS40`, `LS80`,`S8`, `S16`, `S32`, `S64`,`S80`, `S104`.
  * Under the condition that the category is the `HighAvailability`, Valid values: `LC20`, `LC40`, `LC80`,`C8`, `C16`, `C32`, `C64`, `C80`, `C104`.
* `db_cluster_network_type` - (Required, ForceNew) The DBCluster network type. Valid values: `vpc`.
* `db_cluster_version` - (Required, ForceNew) The DBCluster version. Valid values: `20.3.10.75`, `20.8.7.15`, `21.8.10.19`, `22.8.5.29`, `23.8`. **NOTE:** `19.15.2.2` is no longer supported. From version 1.191.0, `db_cluster_version` can be set to `22.8.5.29`.
* `db_node_storage` - (Required) The db node storage.
* `db_node_group_count` - (Required, ForceNew) The db node group count. The number should between 1 and 48.
* `encryption_key` - (Optional, ForceNew) Key management service KMS key ID. It is valid and required when encryption_type is `CloudDisk`.
* `encryption_type` - (Optional, ForceNew) Currently only supports ECS disk encryption, with a value of CloudDisk, not encrypted when empty.
* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `PayAsYouGo`,`Subscription`.
* `renewal_status` - (Optional, Computed, Available since v1.215.0) The renewal status of the resource. Valid values: `AutoRenewal`,`Normal`. It is valid and required when payment_type is `Subscription`. When `renewal_status` is set to `AutoRenewal`, the resource is renewed automatically.
* `period` - (Optional) Pre-paid cluster of the pay-as-you-go cycle. It is valid and required when payment_type is `Subscription`. Valid values: `Month`, `Year`.
* `storage_type` - (Required, ForceNew) Storage type of DBCluster. Valid values: `cloud_essd`, `cloud_efficiency`, `cloud_essd_pl2`, `cloud_essd_pl3`.
* `used_time` - (Optional) The used time of DBCluster. It is valid and required when payment_type is `Subscription`. item choices: [1-9] when period is `Month`, [1-3] when period is `Year`.
* `vswitch_id` - (Optional, ForceNew) The vswitch id of DBCluster.
* `db_cluster_description` - (Optional) The DBCluster description.
* `status` - (Optional) The status of the resource. Valid values: `Running`,`Creating`,`Deleting`,`Restarting`,`Preparing`.
* `maintain_time` - (Optional) The maintenance window of DBCluster. Valid format: `hh:mmZ-hh:mm Z`.
* `db_cluster_access_white_list` - (Optional, Available since v1.145.0) The db cluster access white list. See [`db_cluster_access_white_list`](#db_cluster_access_white_list) below.
* `vpc_id` - (Optional, ForceNew, Available since v1.185.0) The id of the VPC.
* `zone_id` - (Optional, ForceNew, Available since v1.185.0) The zone ID of the instance.
* `multi_zone_vswitch_list` - (Optional, ForceNew, Available since v1.228.0) The zone IDs and 
corresponding vswitch IDs and zone IDs of multi-zone setup. if set, a multi-zone DBCluster will be created. Currently only support 2 available zones, primary zone not included. See [`multi_zone_vswitch_list`](#multi_zone_vswitch_list) below.

### `db_cluster_access_white_list`

The db_cluster_access_white_list supports the following:

* `db_cluster_ip_array_attribute` - (Optional, Removed) Field `db_cluster_ip_array_attribute` has been removed from provider.
* `db_cluster_ip_array_name` - (Optional) Whitelist group name.
* `security_ip_list` - (Optional) The IP address list under the whitelist group.

### `multi_zone_vswitch_list`

The multi_zone_vswitch_list supports the following:
* `zone_id` - (Optional, ForceNew, Available since v1.228.0) The zone ID of the vswitch.
* `vswitch_id` - (Required, ForceNew, Available since v1.228.0) The ID of the vswitch.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DBCluster.
* `connection_string` (Available since v1.196.0) - The connection string of the cluster.
* `port` - (Available since v1.196.0) The connection port of the cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the Click House DBCluster (until it reaches the initial `Running` status).
* `update` - (Defaults to 60 mins) Used when update the Click House DBCluster.

## Import

Click House DBCluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_db_cluster.example <id>
```