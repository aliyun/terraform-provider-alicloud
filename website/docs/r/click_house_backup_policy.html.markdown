---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_backup_policy"
sidebar_current: "docs-alicloud-resource-click-house-backup-policy"
description: |-
  Provides a Alicloud Click House Backup Policy resource.
---

# alicloud\_click\_house\_backup\_policy

Provides a Click House Backup Policy resource.

For information about Click House Backup Policy and how to use it, see [What is Backup Policy](https://www.alibabacloud.com/help/doc-detail/208840.html).

-> **NOTE:** Available in v1.147.0+.

-> **NOTE:** Only the cloud database ClickHouse cluster version `20.3` supports data backup.

## Example Usage

Basic Usage

```terraform
data "alicloud_click_house_regions" "default" {
  current = true
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
}

resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "20.3.10.75"
  status                  = "Running"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_cluster_description  = var.name
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = data.alicloud_vswitches.default.vswitches.0.id
}

resource "alicloud_click_house_backup_policy" "example" {
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Backup Policy.

## Import

Click House Backup Policy can be imported using the id, e.g.

```
$ terraform import alicloud_click_house_backup_policy.example <db_cluster_id>
```