---
subcategory: "Database Backup(DBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbs_backup_plan"
sidebar_current: "docs-alicloud-resource-dbs-backup-plan"
description: |-
  Provides a Alicloud DBS Backup Plan resource.
---

# alicloud_dbs_backup_plan

Provides a DBS Backup Plan resource.

For information about DBS Backup Plan and how to use it, see [What is Backup Plan](https://www.alibabacloud.com/help/en/dms/developer-reference/api-dbs-2019-03-06-createandstartbackupplan).

-> **NOTE:** Available since v1.185.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dbs_backup_plan&exampleId=dedb99ed-bba5-31d5-a52c-dd862a43bd01c308b20d&activeTab=example&spm=docs.r.dbs_backup_plan.0.dedb99edbb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.1.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.1.id
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids.0
  zone_id    = data.alicloud_db_zones.default.ids.1
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id               = local.vswitch_id
  instance_name            = var.name
}

resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.default.id
  name        = "tfdatabase"
}

resource "alicloud_rds_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = "tfnormal000"
  account_password = "Test12345"
}

resource "alicloud_db_account_privilege" "default" {
  instance_id  = alicloud_db_instance.default.id
  account_name = alicloud_rds_account.default.account_name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default.name]
}

resource "alicloud_dbs_backup_plan" "default" {
  backup_plan_name              = var.name
  payment_type                  = "PayAsYouGo"
  instance_class                = "xlarge"
  backup_method                 = "logical"
  database_type                 = "MySQL"
  database_region               = "cn-hangzhou"
  storage_region                = "cn-hangzhou"
  instance_type                 = "RDS"
  source_endpoint_instance_type = "RDS"
  resource_group_id             = data.alicloud_resource_manager_resource_groups.default.ids.0
  source_endpoint_region        = "cn-hangzhou"
  source_endpoint_instance_id   = alicloud_db_instance.default.id
  source_endpoint_user_name     = alicloud_db_account_privilege.default.account_name
  source_endpoint_password      = alicloud_rds_account.default.account_password
  backup_objects                = "[{\"DBName\":\"${alicloud_db_database.default.name}\"}]"
  backup_period                 = "Monday"
  backup_start_time             = "14:22"
  backup_storage_type           = "system"
  backup_retention_period       = 740
}
```

## Argument Reference

The following arguments are supported:

* `backup_gateway_id` - (Optional, ForceNew) The ID of the backup gateway. This parameter is required when the `source_endpoint_instance_type` is `Agent`.
* `backup_log_interval_seconds` - (Optional) The backup log interval seconds.
* `backup_method` - (Required, ForceNew) Backup method. Valid values: `duplication`, `logical`, `physical`.
* `backup_objects` - (Optional, ForceNew) The backup object.
* `backup_period` - (Optional, ForceNew) Full backup cycle, Valid values: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`. supports the selection of multiple fetch values, separated by English commas (,).
* `backup_plan_name` - (Required, ForceNew) The name of the resource.
* `backup_rate_limit` - (Optional) The backup rate limit.
* `backup_retention_period` - (Optional, ForceNew) The retention time of backup data. Valid values: 0 to 1825. Default value: 730 days.
* `backup_speed_limit` - (Optional) The backup speed limit.
* `backup_start_time` - (Optional, ForceNew) The start time of full Backup. The format is `<I> HH:mm</I>` Z(UTC time). 
* `backup_storage_type` - (Optional, ForceNew) Built-in storage type, Valid values: `system`.
* `backup_strategy_type` - (Optional) The backup strategy type. Valid values: `simple`, `manual`.
* `cross_aliyun_id` - (Optional, ForceNew) The UID that is backed up across Alibaba cloud accounts. 
* `cross_role_name` - (Optional, ForceNew) The name of the RAM role that is backed up across Alibaba cloud accounts.
* `database_region` - (Optional) The database region.
* `database_type` - (Required, ForceNew) Database type. Valid values: `DRDS`, `FIle`, `MSSQL`, `MariaDB`, `MongoDB`, `MySQL`, `Oracle`, `PPAS`, `PostgreSQL`, `Redis`.
* `duplication_archive_period` - (Optional, ForceNew) The storage time for conversion to archive cold standby is 365 days by default.
* `duplication_infrequent_access_period` - (Optional, ForceNew) The storage time is converted to low-frequency access. The default time is 180 days.
* `enable_backup_log` - (Optional, ForceNew) Whether to enable incremental log Backup.
* `instance_class` - (Required, ForceNew) The instance class. Valid values: `large`, `medium`, `micro`, `small`, `xlarge`.
* `instance_type` - (Optional) The instance type. Valid values: `RDS`, `PolarDB`, `DDS`, `Kvstore`, `Other`.
* `oss_bucket_name` - (Optional, ForceNew) The OSS Bucket name. The system automatically generates a new name by default.
* `payment_type` - (Optional, ForceNew) The payment type of the resource. Valid values: `PayAsYouGo`, `Subscription`.
* `period` - (Optional) Specify that the prepaid instance is of the package year or monthly type. Valid values: `Month`, `Year`.
* `resource_group_id` - (Optional) The ID of the resource group.
* `source_endpoint_database_name` - (Optional, ForceNew) The name of the database. This parameter is required when the `database_type` is `PostgreSQL` or `MongoDB`.
* `source_endpoint_instance_id` - (Optional, ForceNew) The ID of the database instance. This parameter is required when the `source_endpoint_instance_type` is `RDS`, `ECS`, `DDS`, or `Express`.
* `source_endpoint_instance_type` - (Required, ForceNew) The location of the database. Valid values: `RDS`, `ECS`, `Express`, `Agent`, `DDS`, `Other`.
* `source_endpoint_ip` - (Optional) The source endpoint ip.
* `source_endpoint_password` - (Optional, ForceNew) The source endpoint password.  This parameter is not required when the `database_type` is `Redis`, or when the `source_endpoint_instance_type` is `Agent` and the `database_type` is `MSSQL`. This parameter is required in other scenarios.
* `source_endpoint_port` - (Optional) The source endpoint port.
* `source_endpoint_region` - (Optional, ForceNew) The region of the database. This parameter is required when the `source_endpoint_instance_type` is `RDS`, `ECS`, `DDS`, `Express`, or `Agent`.
* `source_endpoint_sid` - (Optional, ForceNew) Oracle SID name. This parameter is required when the `database_type` is `Oracle`.
* `source_endpoint_user_name` - (Optional, ForceNew) The source endpoint username. This parameter is not required when the `database_type` is `Redis`, or when the `source_endpoint_instance_type` is `Agent` and the `database_type` is `MSSQL`. This parameter is required in other scenarios.
* `status` - (Optional) The status of the resource. Valid values: `pause`, `running`.
* `storage_region` - (Optional) The storage region.
* `used_time` - (Optional) Specify purchase duration. When the parameter `period` is `Year`, the `used_time` value is 1 to 9. When the parameter `period` is `Month`, the `used_time` value is 1 to 11.
* `source_endpoint_oracle_sid` - (Optional) Oracle SID name. This parameter is required when the `database_type` is `Oracle`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Backup Plan.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Backup Plan.
* `update` - (Defaults to 10 mins) Used when update the Backup Plan.
* `delete` - (Defaults to 3 mins) Used when delete the Backup Plan.

## Import

DBS Backup Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_dbs_backup_plan.example <id>
```
