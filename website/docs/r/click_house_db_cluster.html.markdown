---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_db_cluster"
sidebar_current: "docs-alicloud-resource-click-house-db-cluster"
description: |-
  Provides a Alicloud Click House DBCluster resource.
---

# alicloud\_click\_house\_db\_cluster

Provides a Click House DBCluster resource.

For information about Click House DBCluster and how to use it, see [What is DBCluster](https://www.alibabacloud.com/product/clickhouse).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "20.3.10.75"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = "your_vswitch_id"
}

```

## Argument Reference

The following arguments are supported:

* `category` - (Required, ForceNew) The Category of DBCluster. Valid values: `Basic`,`HighAvailability`.
* `db_cluster_class` - (Required) The DBCluster class. According to the category, db_cluster_class has two value ranges:
  * Under the condition that the category is the `Basic`, Valid values: `S4-NEW`, `S8`, `S16`, `S32`, `S64`, `S104`.
  * Under the condition that the category is the `HighAvailability`, Valid values: `C4-NEW`, `C8`, `C16`, `C32`, `C64`, `C104`.
* `db_cluster_network_type` - (Required, ForceNew) The DBCluster network type. Valid values: `vpc`.
* `db_cluster_version` - (Required, ForceNew) The DBCluster version. Valid values: `19.15.2.2`, `20.3.10.75`, `20.8.7.15`.
* `db_node_storage` - (Required, ForceNew) The db node storage.
* `db_node_group_count` - (Required) The db node group count. The number should between 1 and 48.
* `encryption_key` - (Optional, ForceNew) Key management service KMS key ID.
* `encryption_type` - (Optional, ForceNew) Currently only supports ECS disk encryption, with a value of CloudDisk, not encrypted when empty.
* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `PayAsYouGo`,`Subscription`.
* `period` - (Optional, ForceNew) Pre-paid cluster of the pay-as-you-go cycle. Valid values: `Month`, `Year`.
* `storage_type` - (Required, ForceNew) Storage type of DBCluster. Valid values: `cloud_essd`, `cloud_efficiency`, `cloud_essd_pl2`, `cloud_essd_pl3`.
* `used_time` - (Optional, ForceNew) The used time of DBCluster.
* `vswitch_id` - (Optional, ForceNew) The vswitch id of DBCluster.
* `db_cluster_description` - (Optional) The DBCluster description.
* `status` - (Optional, Computed) The status of the resource. Valid values: `Running`,`Creating`,`Deleting`,`Restarting`,`Preparing`,.
* `maintain_time` - (Optional) The maintenance window of DBCluster. Valid format: `hh:mmZ-hh:mm Z`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DBCluster.

## Import

Click House DBCluster can be imported using the id, e.g.

```
$ terraform import alicloud_click_house_db_cluster.example <id>
```