---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_server_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-server-backup-plan"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Server Backup Plan resource.
---

# alicloud\_hbr\_server\_backup\_plan

Provides a Hybrid Backup Recovery (HBR) Server Backup Plan resource.

For information about Hybrid Backup Recovery (HBR) Server Backup Plan and how to use it, see [What is Server Backup Plan](https://www.alibabacloud.com/help/doc-detail/211140.htm).

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_hbr_server_backup_plan&exampleId=4c2af01f-d884-b2c5-fe19-a0e59accc44ef1665bfd&activeTab=example&spm=docs.r.hbr_server_backup_plan.0.4c2af01fd8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "example" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_security_group" "example" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  image_id             = data.alicloud_images.example.images.0.id
  instance_type        = data.alicloud_instance_types.example.instance_types.0.id
  availability_zone    = data.alicloud_zones.example.zones.0.id
  security_groups      = [alicloud_security_group.example.id]
  instance_name        = "terraform-example"
  internet_charge_type = "PayByBandwidth"
  vswitch_id           = alicloud_vswitch.example.id
}

resource "alicloud_hbr_server_backup_plan" "example" {
  ecs_server_backup_plan_name = "terraform-example"
  instance_id                 = alicloud_instance.example.id
  schedule                    = "I|1602673264|PT2H"
  retention                   = 1
  detail {
    app_consistent = true
    snapshot_group = true
  }
  disabled = false
}
```

## Argument Reference

The following arguments are supported:

* `ecs_server_backup_plan_name` - (Required) The name of the backup plan. 1~64 characters, the backup plan name of each data source type in a single warehouse required to be unique.
* `instance_id` - (Required, ForceNew) The ID of ECS instance.
* `retention` - (Required) Backup retention days, the minimum is 1.
* `schedule` - (Required) Backup strategy. Optional format: `I|{startTime}|{interval}`
  * `startTime` Backup start time, UNIX time, in seconds. 
  * `interval` **ISO8601 time interval**. E.g: `PT1H` means one hour apart. `P1D` means one day apart. It means to execute a backup task every `{interval}` starting from `{startTime}`. The backup task for the elapsed time will not be compensated. If the last backup task is not completed, the next backup task will not be triggered.
* `detail` - (Required) ECS server backup plan details.
* `disabled` - (Optional) Whether to disable the backup task. Valid values: `true`, `false`.
* `cross_account_type` - (Optional, ForceNew, Computed, Available in v1.199.0+) The type of the cross account backup. Valid values: `SELF_ACCOUNT`, `CROSS_ACCOUNT`.
* `cross_account_user_id` - (Optional, ForceNew, Available in v1.199.0+) The original account ID of the cross account backup managed by the current account.
* `cross_account_role_name` - (Optional, ForceNew, Available in v1.199.0+) The role name created in the original account RAM backup by the cross account managed by the current account.

#### Block detail

The `detail` supports the following:

* `app_consistent` - (Required) Whether to turn on application consistency. The application consistency snapshot backs up memory data and ongoing database transactions at the time of snapshot creation to ensure the consistency of application system data and database transactions. By applying consistent snapshots, there is no data damage or loss, so as to avoid log rollback during database startup and ensure that the application is in a consistent startup state. Valid values: `true`, `false`.
* `snapshot_group` - (Required) Whether to turn on file system consistency. If SnapshotGroup is true, when AppConsistent is true but the relevant conditions are not met or AppConsistent is false, the resulting snapshot will be a file system consistency snapshot. The file system consistency ensures that the file system memory and disk information are synchronized at the time of snapshot creation, and the file system write operation is frozen to make the file system in a consistent state. The file system consistency snapshot can prevent the operating system from performing disk inspection and repair operations such as CHKDSK or fsck after restart. Valid values: `true`, `false`.
* `enable_fs_freeze` - (Optional) Only the Linux system is valid. Whether to use the Linux FsFreeze mechanism to ensure that the file system is read-only consistent before creating a storage snapshot. The default is True. Valid values: `true`, `false`.
* `pre_script_path` - (Optional) Only vaild for the linux system when AppConsistent is true. Apply the freeze script path (e.g. /tmp/prescript.sh). prescript.sh scripts must meet the following conditions: in terms of permissions, only root, as the owner, has read, write, and execute permissions, that is, 700 permissions. In terms of content, the script content needs to be customized according to the application itself. This indicates that this parameter must be set when creating an application consistency snapshot for a Linux instance. If the script is set incorrectly (for example, permissions, save path, or file name are set incorrectly), the resulting snapshot is a file system consistency snapshot.
* `post_script_path` - (Optional) Only vaild for the linux system when AppConsistent is true. The application thaw script path (e.g. /tmp/postscript.sh). The postscript.sh script must meet the following conditions: in terms of permissions, only the root user as the owner has read, write, and execute permissions, that is, 700 permissions. In terms of content, the script content needs to be customized according to the application itself. This indicates that this parameter must be set when creating an application consistency snapshot for a Linux instance. If the script is set incorrectly (for example, permissions, save path, or file name are set incorrectly), the resulting snapshot is a file system consistency snapshot.
* `timeout_in_seconds` - (Optional) Only the Linux system is valid, and the IO freeze timeout period. The default is 30 seconds.
* `disk_id_list` - (Optional) The list of cloud disks to be backed up in the ECS instance. When not specified, a snapshot is executed for all the disks on the ECS instance.
* `do_copy` - (Optional) Whether replicate to another region. Valid values: `true`, `false`.
* `destination_region_id` - (Optional) Only vaild when DoCopy is true. The destination region ID when replicating to another region. **Note:** Once you set a value of this property, you cannot set it to an empty string anymore.
* `destination_retention` - (Optional) Only vaild when DoCopy is true. The retention days of the destination backup. When not specified, the destination backup will be saved permanently. **Note:** Once you set a value of this property, you cannot set it to an empty string anymore.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Server Backup Plan.

## Import

Hybrid Backup Recovery (HBR) Server Backup Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_server_backup_plan.example <id>
```
