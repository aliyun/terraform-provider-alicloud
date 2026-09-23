---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_server_backup_plans"
sidebar_current: "docs-alicloud-datasource-hbr-server-backup-plans"
description: |-
  Provides a list of Hbr Server Backup Plans to the user.
---

# alicloud\_hbr\_server\_backup\_plans

This data source provides the Hbr Server Backup Plans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_instances" "default" {
  name_regex = "no-deleteing-hbr-ecs-server-backup-plan"
  status     = "Running"
}

data "alicloud_hbr_server_backup_plans" "ids" {
  filters {
    key    = "instanceId"
    values = [data.alicloud_instances.default.instances.0.id]
  }
}

output "hbr_server_backup_plan_id_1" {
  value = data.alicloud_hbr_server_backup_plans.ids.plans.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Server Backup Plan IDs.
* `filters` - (Optional, ForceNew) The filters.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

### Block filter

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:
* `key` - (Required) The key of the field to filter. Valid values: `planId`, `instanceId`, `planName`.
* `values` - (Required) Set of values that are accepted for the given field.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `plans` - A list of HBR server backup plan. Each element contains the following attributes:
  * `id` - The ID of the server backup plan.
  * `ecs_server_backup_plan_id` - The ID of the server backup plan.
  * `ecs_server_backup_plan_name` - The name of the backup plan. 1~64 characters, the backup plan name of each data source type in a single warehouse required to be unique.
  * `instance_id` - The ID of ECS Instance.
  * `retention` - Backup retention days, the minimum is 1.
  * `schedule` - Backup strategy.
  * `detail` - ECS server backup plan details.
    * `app_consistent` - Whether to turn on application consistency. The application consistency snapshot backs up memory data and ongoing database transactions at the time of snapshot creation to ensure the consistency of application system data and database transactions. By applying consistent snapshots, there is no data damage or loss, so as to avoid log rollback during database startup and ensure that the application is in a consistent startup state. Valid values: `true`, `false`.
    * `snapshot_group` - Whether to turn on file system consistency. If SnapshotGroup is true, when AppConsistent is true but the relevant conditions are not met or AppConsistent is false, the resulting snapshot will be a file system consistency snapshot. The file system consistency ensures that the file system memory and disk information are synchronized at the time of snapshot creation, and the file system write operation is frozen to make the file system in a consistent state. The file system consistency snapshot can prevent the operating system from performing disk inspection and repair operations such as CHKDSK or fsck after restart. Valid values: `true`, `false`.
    * `enable_fs_freeze` - Only the Linux system is valid. Whether to use the Linux FsFreeze mechanism to ensure that the file system is read-only consistent before creating a storage snapshot. The default is True. Valid values: `true`, `false`.
    * `pre_script_path` - Only vaild for the linux system when AppConsistent is true. Apply the freeze script path (e.g. /tmp/prescript.sh). prescript.sh scripts must meet the following conditions: in terms of permissions, only root, as the owner, has read, write, and execute permissions, that is, 700 permissions. In terms of content, the script content needs to be customized according to the application itself. This indicates that this parameter must be set when creating an application consistency snapshot for a Linux instance. If the script is set incorrectly (for example, permissions, save path, or file name are set incorrectly), the resulting snapshot is a file system consistency snapshot.
    * `post_script_path` - Only vaild for the linux system when AppConsistent is true. The application thaw script path (e.g. /tmp/postscript.sh). The postscript.sh script must meet the following conditions: in terms of permissions, only the root user as the owner has read, write, and execute permissions, that is, 700 permissions. In terms of content, the script content needs to be customized according to the application itself. This indicates that this parameter must be set when creating an application consistency snapshot for a Linux instance. If the script is set incorrectly (for example, permissions, save path, or file name are set incorrectly), the resulting snapshot is a file system consistency snapshot.
    * `timeout_in_seconds` - Only the Linux system is valid, and the IO freeze timeout period. The default is 30 seconds.
    * `disk_id_list` - The list of cloud disks to be backed up in the ECS instance. When not specified, a snapshot is executed for all the disks on the ECS instance.
    * `do_copy` - Whether replicate to another region. Valid values: `true`, `false`.
    * `destination_region_id` - Only vaild when DoCopy is true. The destination region ID when replicating to another region. **Note:** Once you set a value of this property, you cannot set it to an empty string anymore.
    * `destination_retention` - Only vaild when DoCopy is true. The retention days of the destination backup. When not specified, the destination backup will be saved permanently. **Note:** Once you set a value of this property, you cannot set it to an empty string anymore.

  * `create_time` - The creation time of backup plan.
  * `disabled` - Whether to disable the backup task. Valid values: `true`, `false`.
