---
subcategory: "EBS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_enterprise_snapshot_policy"
sidebar_current: "docs-alicloud-resource-ebs-enterprise-snapshot-policy"
description: |-
  Provides a Alicloud Ebs Enterprise Snapshot Policy resource.
---

# alicloud_ebs_enterprise_snapshot_policy

Provides a Ebs Enterprise Snapshot Policy resource.

For information about Ebs Enterprise Snapshot Policy and how to use it, see [What is Enterprise Snapshot Policy](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/block-storage-overview-elastic-block-storage-devices).

-> **NOTE:** Available in v1.200.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ebs_enterprise_snapshot_policy" "default" {
  status = "ENABLED"
  desc   = var.name
  schedule {
    cron_expression = "0 0 */12 * * *"
  }
  target_type = "DISK"
  retain_rule {
    time_interval = 1
    time_unit     = "DAYS"
  }
  cross_region_copy_info {
    enabled = true
    regions {
      region_id   = "cn-hangzhou"
      retain_days = 7
    }
  }
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
  enterprise_snapshot_policy_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `cross_region_copy_info` - (Optional) Snapshot replication informationSee the following `Block CrossRegionCopyInfo`.
* `desc` - (Optional) Description information representing the resource.
* `enterprise_snapshot_policy_name` - (Required) The name of the resource.
* `resource_group_id` - (ForceNew,Optional) The ID of the resource group.
* `retain_rule` - (Required) Snapshot retention policy representing resourcesSee the following `Block RetainRule`.
* `schedule` - (Required) The scheduling plan that represents the resource.See the following `Block Schedule`.
* `status` - (Optional,Computed) The status of the resource. Valid values: `ENABLED`, `DISABLED`.
* `storage_rule` - (Optional) Snapshot storage policySee the following `Block StorageRule`.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `target_type` - (Required,ForceNew) Represents the target type of resource binding

#### Block Regions

The Regions supports the following:
* `region_id` - (Optional) Destination region ID.
* `retain_days` - (Optional) Number of days of snapshot retention for replication.

#### Block CrossRegionCopyInfo

The CrossRegionCopyInfo supports the following:
* `enabled` - (Optional) Enable Snapshot replication.
* `regions` - (Optional) Destination region for Snapshot replication.See the following `Block Regions`.

#### Block RetainRule

The RetainRule supports the following:
* `number` - (Optional) Retention based on counting method.
* `time_interval` - (Optional) Time unit.
* `time_unit` - (Optional) Time-based retention.

#### Block Schedule

The Schedule supports the following:
* `cron_expression` - (Required) CronTab expression.

#### Block StorageRule

The StorageRule supports the following:
* `enable_immediate_access` - (Optional) Snapshot speed available.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - The creation time of the resource

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Enterprise Snapshot Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise Snapshot Policy.
* `update` - (Defaults to 5 mins) Used when update the Enterprise Snapshot Policy.

## Import

Ebs Enterprise Snapshot Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_enterprise_snapshot_policy.example <id>
```