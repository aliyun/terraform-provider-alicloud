---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ots_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-ots-backup-plan"
description: |-
  Provides a Alicloud HBR Ots Backup Plan resource.
---

# alicloud_hbr_ots_backup_plan

Provides a HBR Ots Backup Plan resource.

For information about HBR Ots Backup Plan and how to use it, see [What is Ots Backup Plan](https://www.alibabacloud.com/help/en/hybrid-backup-recovery/latest/overview).

-> **NOTE:** Available since v1.163.0.

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
  vault_type = "STANDARD"
}

resource "alicloud_ots_instance" "default" {
  name        = "Example-${random_integer.default.result}"
  description = "terraform-example"
  accessed_by = "Any"
  tags = {
    Created = "TF"
    For     = "example"
  }
}

resource "alicloud_ots_table" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = "terraform_example"
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  time_to_live                  = -1
  max_version                   = 1
  deviation_cell_version_in_sec = 1
}

resource "alicloud_ram_role" "default" {
  name     = "hbrexamplerole"
  document = <<EOF
		{
			"Statement": [
			{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"Service": [
						"crossbackup.hbr.aliyuncs.com"
					]
				}
			}
			],
  			"Version": "1"
		}
  		EOF
  force    = true
}

data "alicloud_account" "default" {}
resource "alicloud_hbr_ots_backup_plan" "example" {
  ots_backup_plan_name    = "terraform-example-${random_integer.default.result}"
  vault_id                = alicloud_hbr_vault.default.id
  backup_type             = "COMPLETE"
  retention               = "1"
  instance_name           = alicloud_ots_instance.default.name
  cross_account_type      = "SELF_ACCOUNT"
  cross_account_user_id   = data.alicloud_account.default.id
  cross_account_role_name = alicloud_ram_role.default.id

  ots_detail {
    table_names = [alicloud_ots_table.default.table_name]
  }
  rules {
    schedule    = "I|1602673264|PT2H"
    retention   = "1"
    disabled    = "false"
    rule_name   = "terraform-example"
    backup_type = "COMPLETE"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_hbr_ots_backup_plan&spm=docs.r.hbr_ots_backup_plan.example&intl_lang=EN_US)

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
* `cross_account_type` - (Optional, ForceNew, Computed, Available since v1.189.0) The type of the cross account backup. Valid values: `SELF_ACCOUNT`, `CROSS_ACCOUNT`.
* `cross_account_user_id` - (Optional, ForceNew, Available since v1.189.0) The original account ID of the cross account backup managed by the current account.
* `cross_account_role_name` - (Optional, ForceNew, Available since v1.189.0) The role name created in the original account RAM backup by the cross account managed by the current account.
* `ots_detail` - (Optional) The details about the Table store instance. See the following `Block ots_detail`. **Note:** Required while source_type equals `OTS_TABLE`.
* `rules` - (Optional,Available since v1.164.0) The backup plan rule. See the following `Block rules`. **Note:** Required while source_type equals `OTS_TABLE`.

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

```shell
$ terraform import alicloud_hbr_ots_backup_plan.example <id>
```
