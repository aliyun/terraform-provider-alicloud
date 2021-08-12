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

For information about HBR Oss Backup Plan and how to use it, see [What is Oss Backup Plan](https://www.alibabacloud.com/product/hybrid-backup-recovery).

-> **NOTE:** Available in v1.131.0+.

## Example Usage

Basic Usage

```terraform

variable "name" {
  default = "%s"
}
resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

data "alicloud_oss_buckets" "default" {
  name_regex = "bosh-cf-blobstore-hz"
}

resource "alicloud_hbr_oss_backup_plan" "example" {
  oss_backup_plan_name = "example_value"
  instance_id          = data.alicloud_oss_buckets.default.buckets.0.name
  vault_id             = alicloud_hbr_vault.default.id
  retention            = "1"
  schedule             = "I|1602673264|PT2H"
  backup_type          = "COMPLETE"
}
```

## Argument Reference

The following arguments are supported:

* `oss_backup_plan_name` - (Required) The Configuration Page of a Backup Plan Name. 1-64 Characters, requiring a Single Warehouse under Each of the Data Source Type Drop-down List of the Configuration Page of a Backup Plan Name Is Unique.
* `vault_id` - (Required, ForceNew) Vault ID.
* `bucket` - (Required, ForceNew) The OSS Bucket Name.
* `retention` - (Required) Backup Retention Period, the Minimum Value of 1.
* `schedule` - (Required) Backup strategy. Optional format: I|{startTime}|{interval} * startTime Backup start time, UNIX time, in seconds. * interval ISO8601 time interval. E.g: ** PT1H, one hour apart. ** P1D, one day apart. It means to execute a backup task every {interval} starting from {startTime}. The backup task for the elapsed time will not be compensated. If the last backup task is not completed, the next backup task will not be triggered.
* `disabled` - (Optional) Whether to Disable the Backup Task. Valid Values: true, false.
* `backup_type` - (Optional, Computed, ForceNew) Backup Type. Valid Values: * Complete. Valid values: `COMPLETE`.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Oss Backup Plan.

## Import

HBR Oss Backup Plan can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_oss_backup_plan.example <id>
```
