---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_nas_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-nas-backup-plan"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Nas Backup Plan resource.
---

# alicloud\_hbr\_nas\_backup\_plan

Provides a HBR Nas Backup Plan resource.

For information about HBR Nas Backup Plan and how to use it, see [What is Nas Backup Plan](https://www.alibabacloud.com/help/doc-detail/132248.htm).

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testAccHBRNas"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = var.name
  encrypt_type  = "1"
}

data "alicloud_nas_file_systems" "default" {
  protocol_type     = "NFS"
  description_regex = alicloud_nas_file_system.default.description
}

resource "alicloud_hbr_nas_backup_plan" "default" {
  depends_on           = ["alicloud_nas_file_system.default"]
  nas_backup_plan_name = var.name
  file_system_id       = alicloud_nas_file_system.default.id
  schedule             = "I|1602673264|PT2H"
  backup_type          = "COMPLETE"
  vault_id             = alicloud_hbr_vault.default.id
  create_time          = data.alicloud_nas_file_systems.default.systems.0.create_time
  retention            = "2"
  path                 = ["/"]
}
```

## Argument Reference

The following arguments are supported:

* `nas_backup_plan_name` - (Required) The name of the backup plan. 1~64 characters, the backup plan name of each data source type in a single warehouse required to be unique.
* `vault_id` - (Required, ForceNew) The ID of Backup vault.
* `retention` - (Required) Backup retention days, the minimum is 1.
* `schedule` - (Required) Backup strategy. Optional format: `I|{startTime}|{interval}`. It means to execute a backup task every `{interval}` starting from `{startTime}`. The backup task for the elapsed time will not be compensated. If the last backup task has not completed yet, the next backup task will not be triggered.
    * `startTime` Backup start time, UNIX time seconds.
    * `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `P1D` means one day apart.
* `disabled` - (Optional) Whether to disable the backup task. Valid values: `true`, `false`.
* `backup_type` - (Required, ForceNew) Backup type. Valid values: `COMPLETE`.
* `file_system_id` - (Required, ForceNew) The File System ID of Nas.
* `create_time` - (Optional, Deprecated) This field has been deprecated from provider version 1.153.0+. The creation time of NAS file system. **Note** The time format of the API adopts the ISO 8601, such as `2021-07-09T15:45:30CST` or `2021-07-09T07:45:30Z`.
* `path` - (Required) List of backup path. Up to 65536 characters. e.g.`["/home", "/var"]`. **Note** You should at least specify a backup path, empty array not allowed here.


-> **Note** `alicloud_hbr_nas_backup_plan` depends on the `alicloud_nas_file_system` and creates a mount point on the file system. If this dependency has not declared, the file system may not be deleted correctly.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Nas Backup Plan.

## Import

HBR Nas Backup Plan can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_nas_backup_plan.example <id>
```
