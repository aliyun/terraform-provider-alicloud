---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_policy"
description: |-
  Provides a Alicloud HBR Policy resource.
---

# alicloud_hbr_policy

Provides a HBR Policy resource.

For information about HBR Policy and how to use it, see [What is Policy](https://www.alibabacloud.com/help/en/cloud-backup/developer-reference/api-hbr-2017-09-08-createpolicyv2).

-> **NOTE:** Available since v1.221.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_hbr_policy&exampleId=1815e3cf-c5d0-cd30-ebab-0d1865f44adaa392e053&activeTab=example&spm=docs.r.hbr_policy.0.1815e3cfc5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_hbr_vault" "defaultyk84Hc" {
  vault_type = "STANDARD"
  vault_name = "example-value-${random_integer.default.result}"
}

resource "alicloud_hbr_policy" "defaultoqWvHQ" {
  policy_name = "example-value-${random_integer.default.result}"
  rules {
    rule_type    = "BACKUP"
    backup_type  = "COMPLETE"
    schedule     = "I|1631685600|P1D"
    retention    = "7"
    archive_days = "0"
    vault_id     = alicloud_hbr_vault.defaultyk84Hc.id
  }
  policy_description = "policy example"
}
```

## Argument Reference

The following arguments are supported:
* `policy_description` - (Optional) The policy description.
* `policy_name` - (Optional) Policy Name
* `rules` - (Optional, List) A list of policy rules See [`rules`](#rules) below.

### `rules`

The rules supports the following:
* `archive_days` - (Optional, Computed, Int) This parameter is required only when the value of `RuleType` is **TRANSITION. The minimum value is 30, and the Retention-ArchiveDays needs to be greater than or equal to 60.
* `backup_type` - (Optional) This parameter is required only when the `RuleType` value is **BACKUP. Backup Type.
* `keep_latest_snapshots` - (Optional, Int) This parameter is required only when `RuleType` is set to `BACKUP`.
* `replication_region_id` - (Optional) Only when the `RuleType` value is.
* `retention` - (Optional, Int) This parameter is required only when the value of `RuleType` is `TRANSITION` or **REPLICATION.
  - `RuleType`: `TRANSITION`: the backup retention time. The minimum value is 1 and the maximum value is 364635, in days.
  - `RuleType`: `REPLICATION`: The minimum value is 1 and the maximum value is 364635. The unit is days.
* `retention_rules` - (Optional, List) This parameter is required only when the value of `RuleType` is `TRANSITION`. See [`retention_rules`](#rules-retention_rules) below.
* `rule_type` - (Required) Rule Type.
* `schedule` - (Optional) This parameter is required only if you set the `RuleType` parameter to `BACKUP`. This parameter specifies the backup schedule settings. Format: `I|{startTime}|{interval}`. The system runs the first backup job at a point in time that is specified in the {startTime} parameter and the subsequent backup jobs at an interval that is specified in the {interval} parameter. The system does not run a backup job before the specified point in time. Each backup job, except the first one, starts only after the previous backup job is complete. For example, `I|1631685600|P1D` specifies that the system runs the first backup job at 14:00:00 on September 15, 2021 and the subsequent backup jobs once a day.  *   startTime: the time at which the system starts to run a backup job. The time must follow the UNIX time format. Unit: seconds. *   interval: the interval at which the system runs a backup job. The interval must follow the ISO 8601 standard. For example, PT1H specifies an interval of one hour. P1D specifies an interval of one day.
* `vault_id` - (Optional) Vault ID.

### `rules-retention_rules`

The rules-retention_rules supports the following:
* `advanced_retention_type` - (Optional) Valid values: `annually`, `MONTHLY`, and `WEEKLY`:- `annually`: the first backup of each year. - `MONTHLY`: The first backup of the month. - `WEEKLY`: The first backup of the week. - `DAILY`: The first backup of the day.
* `retention` - (Optional, Int) Retention time, in days.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Policy creation time
* `rules` - A list of policy rules
  * `rule_id` - Rule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Policy.
* `update` - (Defaults to 5 mins) Used when update the Policy.

## Import

HBR Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_policy.example <id>
```