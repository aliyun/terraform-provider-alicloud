---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_instance_cross_backup_policy"
sidebar_current: "docs-alicloud-resource-rds-instance-cross-backup-policy"
description: |-
  Provide RDS disaster recovery backup policy resources.
---

# alicloud_rds_instance_cross_backup_policy

Provides an RDS instance emote disaster recovery strategy policy resource and used to configure instance emote disaster recovery strategy policy.

For information about RDS cross region backup settings and how to use them, see [What is cross region backup](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/modify-cross-region-backup-settings).

-> **NOTE:** Available since v1.195.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_instance_cross_backup_policy&exampleId=c58314df-3c7f-b868-d9ac-26adbc946f3a1d7fe94e&activeTab=example&spm=docs.r.rds_instance_cross_backup_policy.0.c58314df3c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "local_ssd"
  category                 = "HighAvailability"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.ids.0
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "local_ssd"
  category                 = "HighAvailability"
}
data "alicloud_rds_cross_regions" "regions" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.ids.0
  vswitch_name = var.name
}


resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  instance_charge_type     = "Postpaid"
  category                 = "HighAvailability"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.default.id
  db_instance_storage_type = "local_ssd"
}

resource "alicloud_rds_instance_cross_backup_policy" "default" {
  instance_id         = alicloud_db_instance.default.id
  cross_backup_region = data.alicloud_rds_cross_regions.regions.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the instance.
* `log_backup_enabled` - (Optional)The status of the cross-region log backup feature on the instance. Valid values:
  - Enable: Enables the feature.
  - Disabled: Disables the feature.
* `cross_backup_region` - (Required) The ID of the destination region where the cross-region backup files of the instance are stored.
* `retention` - (Optional) The number of days for which the cross-region backup files of the instance are retained. Valid values: 7 to 1825. Default value: 7.

## Attributes Reference

The following attributes are exported:

* `id` - The Id of DB instance.
* `backup_enabled` - The status of the overall cross-region backup switch on the instance. Valid values:
  - Disabled
  - Enable
* `backup_enabled_time` - The time when cross-region backup was enabled on the instance. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `log_backup_enabled_time` - The time when cross-region log backup was enabled on the instance. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `db_instance_status` - The state of the instance. For more information, see Instance status.
* `lock_mode` - The lock status of the instance. Valid values:
  - Unlock: The instance is not locked.
  - ManualLock: The instance is manually locked.
  - LockByExpiration: The instance is locked upon expiration.
  - LockByRestoration: The instance is automatically locked before a rollback.
  - LockByDiskQuota: The instance is automatically locked because its storage space is exhausted. In this situation, the instance is inaccessible.
* `retent_type` - The policy that is used to retain cross-region backups of the instance. Default value: 1. The default value 1 indicate that cross-region backups are retained based on the specified retention period.
* `cross_backup_type` - The policy that is used to save cross-region backups of the instance. Default value: 1. The default value 1 indicates that all cross-region backups are saved.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when the database instance is set for remote disaster recovery.
* `update` - (Defaults to 10 mins) Used when the database instance modifies the remote disaster recovery settings.
* `delete` - (Defaults to 10 mins) Used when the database instance shuts down remote disaster recovery.

## Import

RDS remote disaster recovery policies can be imported using id or instance id, e.g.

```shell
$ terraform import alicloud_rds_instance_cross_backup_policy.example "rm-12345678"
```
