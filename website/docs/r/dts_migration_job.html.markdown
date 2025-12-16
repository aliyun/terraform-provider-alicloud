---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_migration_job"
sidebar_current: "docs-alicloud-resource-dts-migration-job"
description: |-
  Provides a Alicloud DTS Migration Job resource.
---

# alicloud_dts_migration_job

Provides a DTS Migration Job resource.

For information about DTS Migration Job and how to use it, see [What is Migration Job](https://www.alibabacloud.com/help/en/doc-detail/208399.html).

-> **NOTE:** Available since v1.157.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dts_migration_job&exampleId=ea959101-7442-6daa-8188-6eaa970d360efac93b4e&activeTab=example&spm=docs.r.dts_migration_job.0.ea95910174&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}
data "alicloud_regions" "example" {
  current = true
}
data "alicloud_db_zones" "example" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "Basic"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "example" {
  zone_id                  = data.alicloud_db_zones.example.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "Basic"
  db_instance_storage_type = "cloud_essd"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.example.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_db_instance" "example" {
  count                    = 2
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.example.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.example.instance_classes.0.storage_range.min
  instance_charge_type     = "Postpaid"
  instance_name            = format("${var.name}_%d", count.index + 1)
  vswitch_id               = alicloud_vswitch.example.id
  monitoring_period        = "60"
  db_instance_storage_type = "cloud_essd"
  security_group_ids       = [alicloud_security_group.example.id]
}

resource "alicloud_rds_account" "example" {
  count            = 2
  db_instance_id   = alicloud_db_instance.example[count.index].id
  account_name     = format("example_name_%d", count.index + 1)
  account_password = format("example_password_%d", count.index + 1)
}

resource "alicloud_db_database" "example" {
  count       = 2
  instance_id = alicloud_db_instance.example[count.index].id
  name        = format("${var.name}_%d", count.index + 1)
}

resource "alicloud_db_account_privilege" "example" {
  count        = 2
  instance_id  = alicloud_db_instance.example[count.index].id
  account_name = alicloud_rds_account.example[count.index].name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.example[count.index].name]
}

resource "alicloud_dts_migration_instance" "example" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = data.alicloud_regions.example.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = data.alicloud_regions.example.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

resource "alicloud_dts_migration_job" "example" {
  dts_instance_id                    = alicloud_dts_migration_instance.example.id
  dts_job_name                       = var.name
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_account_privilege.example.0.instance_id
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = data.alicloud_regions.example.regions.0.id
  source_endpoint_user_name          = alicloud_rds_account.example.0.account_name
  source_endpoint_password           = alicloud_rds_account.example.0.account_password
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = alicloud_db_account_privilege.example.1.instance_id
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = data.alicloud_regions.example.regions.0.id
  destination_endpoint_user_name     = alicloud_rds_account.example.1.account_name
  destination_endpoint_password      = alicloud_rds_account.example.1.account_password
  db_list = jsonencode(
    {
      "${alicloud_db_database.example.0.name}" = { name = alicloud_db_database.example.1.name, all = true }
    }
  )
  structure_initialization = true
  data_initialization      = true
  data_synchronization     = true
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_dts_migration_job&spm=docs.r.dts_migration_job.example&intl_lang=EN_US)

## Argument Reference

The following arguments supported:

* `dts_instance_id` - (Required, ForceNew) The Migration instance ID. The ID of `alicloud_dts_migration_instance`.
* `dts_job_name` - (Optional, ForceNew) The name of migration job.
* `instance_class` - (Optional) The instance class. Valid values: `large`, `medium`, `micro`, `small`, `xlarge`, `xxlarge`. 
* `checkpoint` - (Optional, ForceNew) Start time in Unix timestamp format.
* `data_initialization` - (Required, ForceNew) Whether to execute DTS supports schema migration.
* `structure_initialization` - (Required, ForceNew) Whether to perform a database table structure to migrate.
* `data_synchronization` - (Required, ForceNew) Whether to perform incremental data migration.
* `db_list` - (Required, ForceNew) Migration object, in the format of JSON strings. For detailed definition instructions, please refer to [the description of migration, migration or subscription objects](https://help.aliyun.com/document_detail/209545.html).
* `source_endpoint_instance_type` - (Required, ForceNew) The type of source instance. Valid values: `CEN`, `DG`, `DISTRIBUTED_DMSLOGICDB`, `ECS`, `EXPRESS`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`.
* `source_endpoint_engine_name` - (Required, ForceNew) The type of source database. Valid values: `AS400`, `DB2`, `DMSPOLARDB`, `HBASE`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `POSTGRESQL`, `TERADATA`.
* `source_endpoint_instance_id` - (Optional, ForceNew) The ID of source instance.
* `source_endpoint_region` - (Optional, ForceNew) The region of source instance.
* `source_endpoint_ip` - (Optional, ForceNew) The ip of source endpoint.
* `source_endpoint_port` - (Optional, ForceNew) The port of source endpoint.
* `source_endpoint_oracle_sid` - (Optional, ForceNew) The SID of Oracle database.
* `source_endpoint_database_name` - (Optional, ForceNew) The name of migrate the database.
* `source_endpoint_user_name` - (Optional, ForceNew) The username of database account.
* `source_endpoint_password` - (Optional) The password of database account.
* `source_endpoint_owner_id` - (Optional, ForceNew) The Alibaba Cloud account ID to which the source instance belongs.
* `source_endpoint_role` - (Optional, ForceNew) The name of the role configured for the cloud account to which the source instance belongs.
* `destination_endpoint_instance_type` - (Required, ForceNew) The type of destination instance. Valid values: `ADS`, `CEN`, `DATAHUB`, `DG`, `ECS`, `EXPRESS`, `GREENPLUM`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`.
* `destination_endpoint_engine_name` - (Required, ForceNew) The type of destination database. Valid values: `ADS`, `ADB30`, `AS400`, `DATAHUB`, `DB2`, `GREENPLUM`, `KAFKA`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `PostgreSQL`.
* `destination_endpoint_instance_id` - (Optional, ForceNew) The ID of destination instance.
* `destination_endpoint_region` - (Optional, ForceNew) The region of destination instance.
* `destination_endpoint_ip` - (Optional, ForceNew) The ip of source endpoint.
* `destination_endpoint_port` - (Optional, ForceNew) The port of source endpoint.
* `destination_endpoint_database_name` - (Optional, ForceNew) The name of migrate the database.
* `destination_endpoint_user_name` - (Optional, ForceNew) The username of database account.
* `destination_endpoint_password` - (Optional) The password of database account.
* `destination_endpoint_oracle_sid` - (Optional, ForceNew) The SID of Oracle database.
* `status` - (Optional) The status of the resource. Valid values: `Migrating`, `Suspending`. You can suspend the task by specifying `Suspending` and start the task by specifying `Migrating`.

## Notice

1. The expiration time cannot be changed after the work of the annual and monthly subscription suspended;
2. After the pay-as-you-go type job suspended, your job configuration fee will still be charged;
3. If the task suspended for more than 6 hours, the task will not start successfully.
4. Suspending the task will only stop writing to the target library, but will still continue to obtain the incremental log of the source, so that the task can be quickly resumed after the suspension is canceled. Therefore, some resources of the source library, such as bandwidth resources, will continue to be occupied during the period.
5. Charges will continue during the task suspension period. If you need to stop charging, please release the instance
6. When a DTS instance suspended for more than 7 days, the instance cannot be resumed, and the status will change from suspended to failed.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Migration Job.

## Import

DTS Migration Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_dts_migration_job.example <id>
```