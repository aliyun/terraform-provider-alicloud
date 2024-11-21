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

For information about Hybrid Backup Recovery (HBR) Hana Backup Plan and how to use it, see [What is Hana Backup Plan](https://www.alibabacloud.com/help/en/hybrid-backup-recovery/latest/api-hbr-2017-09-08-createhanabackupplan).

-> **NOTE:** Available in v1.179.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_hbr_hana_backup_plan&exampleId=72b0888a-4bae-9125-31a8-2f9fe2ddc35aadf0a6a9&activeTab=example&spm=docs.r.hbr_hana_backup_plan.0.72b0888a4b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_resource_manager_resource_groups" "example" {
  status = "OK"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_hbr_vault" "example" {
  vault_name = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_hbr_hana_instance" "example" {
  alert_setting        = "INHERITED"
  hana_name            = "terraform-example-${random_integer.default.result}"
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