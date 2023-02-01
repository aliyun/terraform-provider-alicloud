---
subcategory: "Hbr"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_hana_restore_jobs"
sidebar_current: "docs-alicloud-datasource-hbr-hana-restore-jobs"
description: |-
  Provides a list of Hbr Hana Restore Job owned by an Alibaba Cloud account.
---

# alicloud_hbr_hana_restore_jobs

This data source provides Hbr Hana Restore Job available to the user.[What is Hana Restore Job](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available in 1.198.0+

## Example Usage

```
没有资源测试用例，请先通过资源测试用例后再生成示例代码。
```

## Argument Reference

The following arguments are supported:
* `backup_id` - (ForceNew,Optional) backup id
* `cluster_id` - (Required,ForceNew) cluester id
* `database_name` - (ForceNew,Optional) database name
* `restore_id` - (ForceNew,Optional) restore ID
* `token` - (ForceNew,Optional) 当前属性没有在镇元上录入属性描述，请补充后再生成代码。
* `vault_id` - (ForceNew,Optional) vault id
* `ids` - (Optional, ForceNew, Computed) A list of Hana Restore Job IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Hana Restore Job IDs.
* `jobs` - A list of Hana Restore Job Entries. Each element contains the following attributes:
  * `backup_id` - backup id
  * `backup_prefix` - backup prefix
  * `check_access` - check access
  * `clear_log` - clear log
  * `cluster_id` - cluester id
  * `current_phase` - current  phase
  * `current_progress` - current progress
  * `database_name` - database name
  * `database_restore_id` - data base restore id
  * `end_time` - end time
  * `log_position` - log  position
  * `max_phase` - max phase
  * `max_progress` - max progress
  * `message` - message
  * `mode` - mode
  * `phase` - phase
  * `reached_time` - reached time
  * `recovery_point_in_time` - recovery point time
  * `restore_id` - restore ID
  * `source` - source
  * `source_cluster_id` - source cluster id
  * `start_time` - start time
  * `state` - state
  * `status` - status
  * `system_copy` - system copy
  * `use_catalog` - user catalog
  * `use_delta` - usedelta
  * `vault_id` - vault id
  * `volume_id` - volume id
