---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_enterprise_snapshot_policy"
description: |-
  Provides a Alicloud EBS Enterprise Snapshot Policy resource.
---

# alicloud_ebs_enterprise_snapshot_policy

Provides a EBS Enterprise Snapshot Policy resource. enterprise snapshot policy.

For information about EBS Enterprise Snapshot Policy and how to use it, see [What is Enterprise Snapshot Policy](https://next.api.aliyun.com/api/ebs/2021-07-30/CreateEnterpriseSnapshotPolicy).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ebs_enterprise_snapshot_policy&exampleId=df91d4f1-3fef-2673-fad9-2c9f38d2a711d347f903&activeTab=example&spm=docs.r.ebs_enterprise_snapshot_policy.0.df91d4f13f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ecs_disk" "defaultJkW46o" {
  category          = "cloud_essd"
  description       = "esp-attachment-test"
  zone_id           = "cn-hangzhou-i"
  performance_level = "PL1"
  size              = "20"
  disk_name         = var.name
}

resource "alicloud_ebs_enterprise_snapshot_policy" "defaultPE3jjR" {
  status = "DISABLED"
  desc   = "DESC"
  schedule {
    cron_expression = "0 0 0 1 * ?"
  }
  enterprise_snapshot_policy_name = var.name

  target_type = "DISK"
  retain_rule {
    time_interval = "120"
    time_unit     = "DAYS"
  }
}
```

## Argument Reference

The following arguments are supported:
* `cross_region_copy_info` - (Optional) Snapshot replication information. See [`cross_region_copy_info`](#cross_region_copy_info) below.
* `desc` - (Optional) Description information representing the resource.
* `enterprise_snapshot_policy_name` - (Required) The name of the resource.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `retain_rule` - (Required) Snapshot retention policy representing resources. See [`retain_rule`](#retain_rule) below.
* `schedule` - (Required) The scheduling plan that represents the resource. See [`schedule`](#schedule) below.
* `special_retain_rules` - (Optional, ForceNew) Snapshot special retention rules. See [`special_retain_rules`](#special_retain_rules) below.
* `status` - (Optional, Computed) The status of the resource.
* `storage_rule` - (Optional) Snapshot storage policy. See [`storage_rule`](#storage_rule) below.
* `tags` - (Optional, Map) The tag of the resource.
* `target_type` - (Required, ForceNew) Represents the target type of resource binding.

### `cross_region_copy_info`

The cross_region_copy_info supports the following:
* `enabled` - (Optional) Enable Snapshot replication.
* `regions` - (Optional) Destination region for Snapshot replication. See [`regions`](#cross_region_copy_info-regions) below.

### `cross_region_copy_info-regions`

The cross_region_copy_info-regions supports the following:
* `region_id` - (Optional) Destination region ID.
* `retain_days` - (Optional) Number of days of snapshot retention for replication.

### `retain_rule`

The retain_rule supports the following:
* `number` - (Optional) Retention based on counting method.
* `time_interval` - (Optional) Time unit.
* `time_unit` - (Optional) Time-based retention.

### `schedule`

The schedule supports the following:
* `cron_expression` - (Required) CronTab expression.

### `special_retain_rules`

The special_retain_rules supports the following:
* `enabled` - (Optional) Whether special reservations are enabled. Value range:
  - true
  - false.
* `rules` - (Optional) List of special retention rules. See [`rules`](#special_retain_rules-rules) below.

### `special_retain_rules-rules`

The special_retain_rules-rules supports the following:
* `special_period_unit` - (Optional) The cycle unit of the special reserved snapshot. If the value is set to WEEKS, the first snapshot of each week is reserved. The retention time is determined by TimeUnit and TimeInterval. The value range is:
  - WEEKS
  - MONTHS
  - YEARS.
* `time_interval` - (Optional) Retention time value. The value range is greater than 1.
* `time_unit` - (Optional) The retention time unit for a particular snapshot. Value range:
  - DAYS
  - WEEKS.

### `storage_rule`

The storage_rule supports the following:
* `enable_immediate_access` - (Optional) Snapshot speed available.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Enterprise Snapshot Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise Snapshot Policy.
* `update` - (Defaults to 5 mins) Used when update the Enterprise Snapshot Policy.

## Import

EBS Enterprise Snapshot Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_enterprise_snapshot_policy.example <id>
```