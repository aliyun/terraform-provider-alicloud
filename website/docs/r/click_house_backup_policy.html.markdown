---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_backup_policy"
sidebar_current: "docs-alicloud-resource-click-house-backup-policy"
description: |-
  Provides a Alicloud Click House Backup Policy resource.
---

# alicloud_click_house_backup_policy

Provides a Click House Backup Policy resource.

For information about Click House Backup Policy and how to use it, see [What is Backup Policy](https://www.alibabacloud.com/help/zh/clickhouse/latest/api-clickhouse-2019-11-11-createbackuppolicy).

-> **NOTE:** Available since v1.147.0.

-> **NOTE:** Only the cloud database ClickHouse cluster version `20.3` supports data backup.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_click_house_backup_policy&exampleId=941b444e-b3fd-7248-1589-834e274b2179740808a3&activeTab=example&spm=docs.r.click_house_backup_policy.0.941b444eb3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_click_house_regions" "default" {
  current = true
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_click_house_regions.default.regions.0.zone_ids.1.zone_id
}

resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "22.8.5.29"
  status                  = "Running"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = alicloud_vswitch.default.id
  vpc_id                  = alicloud_vpc.default.id
}

resource "alicloud_click_house_backup_policy" "default" {
  db_cluster_id           = alicloud_click_house_db_cluster.default.id
  preferred_backup_period = ["Monday", "Friday"]
  preferred_backup_time   = "00:00Z-01:00Z"
  backup_retention_period = 7
}
```

## Argument Reference

The following arguments are supported:

* `backup_retention_period` - (Optional) Data backup days. Valid values: `7` to `730`.
* `db_cluster_id` - (Required, ForceNew) The id of the DBCluster.
* `preferred_backup_period` - (Required) DBCluster Backup period. A list of DBCluster Backup period. Valid values: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"].
* `preferred_backup_time` - (Required) DBCluster backup time, in the format of `HH:mmZ-HH:mmZ`. Time setting interval is one hour. China time is 8 hours behind it.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Backup Policy. Its value is same as `db_cluster_id`.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Backup Policy.

## Import

Click House Backup Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_backup_policy.example <db_cluster_id>
```