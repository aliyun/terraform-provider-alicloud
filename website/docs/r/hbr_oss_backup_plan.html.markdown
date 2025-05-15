---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_oss_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-oss-backup-plan"
description: |-
  Provides a Alicloud HBR Oss Backup Plan resource.
---

# alicloud_hbr_oss_backup_plan

Provides a HBR Oss Backup Plan resource.

For information about HBR Oss Backup Plan and how to use it, see [What is Oss Backup Plan](https://www.alibabacloud.com/help/doc-detail/130040.htm).

-> **NOTE:** Available since v1.131.0.

-> **NOTE:** Deprecated since v1.249.0.

-> **DEPRECATED:** This resource has been deprecated from version `1.249.0`. Please use new resource [alicloud_hbr_policy](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/hbr_policy) and [alicloud_hbr_policy_binding](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/hbr_policy_binding).

## Example Usage

Basic Usage

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_hbr_vault" "default" {
  vault_name = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_oss_bucket" "default" {
  bucket = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_hbr_oss_backup_plan" "default" {
  oss_backup_plan_name = "terraform-example"
  # the prefix of object you want to back up
  prefix      = "/example"
  bucket      = alicloud_oss_bucket.default.bucket
  vault_id    = alicloud_hbr_vault.default.id
  schedule    = "I|1602673264|PT2H"
  backup_type = "COMPLETE"
  retention   = "2"
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
* `cross_account_type` - (Optional, ForceNew, Computed, Available since v1.189.0) The type of the cross account backup. Valid values: `SELF_ACCOUNT`, `CROSS_ACCOUNT`.
* `cross_account_user_id` - (Optional, ForceNew, Available since v1.189.0) The original account ID of the cross account backup managed by the current account.
* `cross_account_role_name` - (Optional, ForceNew, Available since v1.189.0) The role name created in the original account RAM backup by the cross account managed by the current account.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Oss Backup Plan.

## Import

HBR Oss Backup Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_oss_backup_plan.example <id>
```
