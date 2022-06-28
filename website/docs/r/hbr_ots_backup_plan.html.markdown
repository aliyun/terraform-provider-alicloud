---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ots_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-ots-backup-plan"
description: |-
  Provides a Alicloud HBR Ots Backup Plan resource.
---

# alicloud\_hbr\_ots\_backup\_plan

Provides a HBR Ots Backup Plan resource.

For information about HBR Ots Backup Plan and how to use it, see [What is Ots Backup Plan](https://www.alibabacloud.com/help/en/hybrid-backup-recovery/latest/overview).

-> **NOTE:** Available in v1.163.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "testAcc"
}
resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
  vault_type = "OTS_BACKUP"
}

resource "alicloud_ots_instance" "foo" {
  name        = var.name
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_ots_table" "basic" {
  instance_name = alicloud_ots_instance.foo.name
  table_name    = var.name
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  time_to_live                  = -1
  max_version                   = 1
  deviation_cell_version_in_sec = 1
}

resource "alicloud_hbr_ots_backup_plan" "example" {
  ots_backup_plan_name = var.name
  vault_id             = alicloud_hbr_vault.default.id
  backup_type          = "COMPLETE"
  schedule             = "I|1602673264|PT2H"
  retention            = "2"
  instance_name        = alicloud_ots_instance.foo.name
  ots_detail {
    table_names = [alicloud_ots_table.basic.table_name]
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_type` - (Required) Backup type. Valid values: `COMPLETE`.
* `disabled` - (Optional) Whether to disable the backup task. Valid values: `true`, `false`. Default values: `false`.
* `ots_backup_plan_name` - (Required) The name of the backup plan. 1~64 characters, the backup plan name of each data source type in a single warehouse required to be unique.
* `retention` - (Required) Backup retention days, the minimum is 1.
* `schedule` - (Optional, Deprecated) Backup strategy. Optional format: `I|{startTime}|{interval}`. It means to execute a backup task every `{interval}` starting from `{startTime}`. The backup task for the elapsed time will not be compensated. If the last backup task has not completed yet, the next backup task will not be triggered.
  - `startTime` Backup start time, UNIX time seconds.
  - `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `P1D` means one day apart.
* `vault_id` - (Required) The ID of backup vault.
* `instance_name` - (Optional) The name of the Table store instance. **Note:** Required while source_type equals `OTS_TABLE`.
* `ots_detail` - (Optional) The details about the Table store instance. See the following `Block ots_detail`. **Note:** Required while source_type equals `OTS_TABLE`.
* `rules` - (Optional,Available in v1.164.0+) The backup plan rule. See the following `Block rules`. **Note:** Required while source_type equals `OTS_TABLE`.


### Block ots_detail

The ots_detail supports the following:

* `table_names` - (Optional) The names of the destination tables in the Tablestore instance. **Note:** Required while source_type equals `OTS_TABLE`.

### Block rules

The rules support the following:

* `schedule` - (Optional) Backup strategy. Optional format: `I|{startTime}|{interval}`. It means to execute a backup task every `{interval}` starting from `{startTime}`. The backup task for the elapsed time will not be compensated. If the last backup task has not completed yet, the next backup task will not be triggered. **Note:** Required while source_type equals `OTS_TABLE`.
  - `startTime` Backup start time, UNIX time seconds.
  - `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `P1D` means one day apart.
* `retention` - (Optional) Backup retention days, the minimum is 1. **Note:** Required while source_type equals `OTS_TABLE`.
* `rule_name` - (Optional)  The name of the backup rule.**Note:** Required while source_type equals `OTS_TABLE`. `rule_name` should be unique for the specific user.
* `backup_type` - (Optional) The name of the tableStore instance. Valid values: `COMPLETE`, `INCREMENTAL`. **Note:** Required while source_type equals `OTS_TABLE`.
* `disabled` - (Optional) Whether to disable the backup task. Valid values: true, false.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ots Backup Plan.

## Import

HBR Ots Backup Plan can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_ots_backup_plan.example <id>
```
