---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_oss_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-oss-backup-plan"
description: |-
  Provides a Alicloud HBR Oss Backup Plan resource.
---

# alicloud\_hbr\_oss\_backup\_plan

Provides a HBR Oss Backup Plan resource.

For information about HBR Oss Backup Plan and how to use it, see [What is Oss Backup Plan](https://www.alibabacloud.com/help/doc-detail/130040.htm).

-> **NOTE:** Available in v1.131.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-test112358"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_hbr_oss_backup_plan" "default" {
  oss_backup_plan_name = var.name
  prefix               = "/"
  bucket               = alicloud_oss_bucket.default.bucket
  vault_id             = alicloud_hbr_vault.default.id
  schedule             = "I|1602673264|PT2H"
  backup_type          = "COMPLETE"
  retention            = "2"
}
```

## Argument Reference

The following arguments are supported:

* `oss_backup_plan_name` - (Required) The name of the backup plan. 1~64 characters, the backup plan name of each data source type in a single warehouse required to be unique.
* `vault_id` - (Required, ForceNew) The ID of backup vault.
* `bucket` - (Required, ForceNew) The name of OSS bucket.
* `retention` - (Required) Backup retention days, the minimum is 1.
* `schedule` - (Required) Backup strategy. Optional format: `I|{startTime}|{interval}`. It means to execute a backup task every `{interval}` starting from `{startTime}`. The backup task for the elapsed time will not be compensated. If the last backup task has not completed yet, the next backup task will not be triggered.
    * `startTime` Backup start time, UNIX time seconds.
    * `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `P1D` means one day apart.
* `backup_type` - (Required, ForceNew) Backup type. Valid values: `COMPLETE`.
* `disabled` - (Optional) Whether to disable the backup task. Valid values: `true`, `false`.
* `prefix` - (Optional) Backup prefix. Once specified, only objects with matching prefixes will be backed up.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Oss Backup Plan.

## Import

HBR Oss Backup Plan can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_oss_backup_plan.example <id>
```
