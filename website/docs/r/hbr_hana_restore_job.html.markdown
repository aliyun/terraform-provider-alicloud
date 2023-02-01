---
subcategory: "Hbr"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_hana_restore_job"
sidebar_current: "docs-alicloud-resource-hbr-hana-restore-job"
description: |-
  Provides a Alicloud Hbr Hana Restore Job resource.
---

# alicloud_hbr_hana_restore_job

Provides a Hbr Hana Restore Job resource.

For information about Hbr Hana Restore Job and how to use it, see [What is Hana Restore Job](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available in v1.198.0+.

## Example Usage

Basic Usage

```terraform
没有资源测试用例，请先通过资源测试用例后再生成示例代码。
```

## Argument Reference

The following arguments are supported:
* `backup_id` - (ForceNew,Optional) backup id
* `backup_prefix` - (Required,ForceNew) backup prefix
* `check_access` - (ForceNew,Optional) check access
* `clear_log` - (ForceNew,Optional) clear log
* `cluster_id` - (Required) cluester id
* `database_name` - (Required) database name
* `log_position` - (ForceNew,Optional) log  position
* `mode` - (Required,ForceNew) mode
* `recovery_point_in_time` - (ForceNew,Optional) recovery point time
* `source` - (ForceNew,Optional) source
* `source_cluster_id` - (ForceNew,Optional) source cluster id
* `system_copy` - (ForceNew,Optional) system copy
* `token` - (Computed,Optional) 当前属性没有在镇元上录入属性描述，请补充后再生成代码。
* `use_catalog` - (ForceNew,Optional) user catalog
* `use_delta` - (ForceNew,Optional) usedelta
* `vault_id` - (Computed,Optional) vault id
* `volume_id` - (ForceNew,Optional) volume id



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `current_phase` - current  phase
* `current_progress` - current progress
* `database_restore_id` - data base restore id
* `end_time` - end time
* `max_phase` - max phase
* `max_progress` - max progress
* `message` - message
* `phase` - phase
* `reached_time` - reached time
* `restore_id` - restore ID
* `start_time` - start time
* `state` - state
* `status` - status
* `token` - 当前属性没有在镇元上录入属性描述，请补充后再生成代码。
* `vault_id` - vault id

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Hana Restore Job.
* `update` - (Defaults to 5 mins) Used when update the Hana Restore Job.

## Import

Hbr Hana Restore Job can be imported using the id, e.g.

```shell
$terraform import alicloud_hbr_hana_restore_job.example <id>
```