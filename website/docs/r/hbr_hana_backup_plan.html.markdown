---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_hana_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-hana-backup-plan"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Hana Backup Plan resource.
---

# alicloud\_hbr\_hana\_backup\_plan

Provides a Hybrid Backup Recovery (HBR) Hana Backup Plan resource.

For information about Hybrid Backup Recovery (HBR) Hana Backup Plan and how to use it, see [What is Hana Backup Plan](https://www.alibabacloud.com/help/en/hybrid-backup-recovery/latest/api-doc-hbr-2017-09-08-api-doc-createhanabackupplan).

-> **NOTE:** Available in v1.179.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "example" {
  status = "OK"
}

resource "alicloud_hbr_vault" "example" {
  vault_name = "terraform-example"
}

resource "alicloud_hbr_hana_instance" "example" {
  alert_setting        = "INHERITED"
  hana_name            = "terraform-example"
  host                 = "1.1.1.1"
  instance_number      = 1
  password             = "YouPassword123"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.example.groups.0.id
  sid                  = "HXE"
  use_ssl              = false
  user_name            = "admin"
  validate_certificate = false
  vault_id             = alicloud_hbr_vault.example.id
}

resource "alicloud_hbr_hana_backup_plan" "example" {
  backup_prefix     = "DIFF_DATA_BACKUP"
  backup_type       = "COMPLETE"
  cluster_id        = alicloud_hbr_hana_instance.example.hana_instance_id
  database_name     = "SYSTEMDB"
  plan_name         = "terraform-example"
  resource_group_id = data.alicloud_resource_manager_resource_groups.example.groups.0.id
  schedule          = "I|1602673264|P1D"
  vault_id          = alicloud_hbr_hana_instance.example.vault_id
}
```

## Argument Reference

The following arguments are supported:

* `backup_prefix` - (Optional) The backup prefix.
* `backup_type` - (Required, ForceNew) The backup type. Valid values:
  - `COMPLETE`: full backup.
  - `INCREMENTAL`: incremental backup.
  - `DIFFERENTIAL`: differential backup.
* `cluster_id` - (Required, ForceNew) The ID of the SAP HANA instance.
* `database_name` - (Required, ForceNew) The name of the database.
* `plan_name` - (Required) The name of the backup plan.
* `resource_group_id` - (Optional) The resource attribute field that represents the resource group ID.
* `schedule` - (Required) The backup policy. Format: `I|{startTime}|{interval}`. The system runs the first backup job at a point in time that is specified in the {startTime} parameter and the subsequent backup jobs at an interval that is specified in the {interval} parameter. The system does not run a backup job before the specified point in time. Each backup job, except the first one, starts only after the previous backup job is completed. For example, I|1631685600|P1D specifies that the system runs the first backup job at 14:00:00 on September 15, 2021 and the subsequent backup jobs once a day.
* `vault_id` - (Required, ForceNew) The ID of the backup vault.
* `status` - (Optional, Computed) The status of the resource. Valid values: `Enabled`, `Disabled`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Hana Backup Plan. The value formats as `<plan_id>:<vault_id>:<cluster_id>`.
* `plan_id` - The id of the plan.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Hana Backup Plan.
* `update` - (Defaults to 1 mins) Used when updating the Hana Backup Plan.
* `delete` - (Defaults to 1 mins) Used when deleting the Hana Backup Plan.


## Import

Hybrid Backup Recovery (HBR) Hana Backup Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_hana_backup_plan.example <plan_id>:<vault_id>:<cluster_id>
```