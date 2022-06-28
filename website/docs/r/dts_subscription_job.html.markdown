---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_subscription_job"
sidebar_current: "docs-alicloud-resource-dts-subscription-job"
description: |-
  Provides a Alicloud DTS Subscription Job resource.
---

# alicloud\_dts\_subscription\_job

Provides a DTS Subscription Job resource.

For information about DTS Subscription Job and how to use it, see [What is Subscription Job](https://help.aliyun.com/document_detail/254791.html).

-> **NOTE:** Available in v1.138.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "dtsSubscriptionJob"
}

variable "creation" {
  default = "Rds"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.instance.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  instance_id = alicloud_db_instance.instance.id
  name        = "tftestprivilege"
  password    = "Test12345"
  description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_instance.instance.id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db.*.name
}


data "alicloud_vpcs" "default1" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default_1" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
}

resource "alicloud_dts_subscription_job" "default" {
  dts_job_name                       = var.name
  payment_type                       = "PostPaid"
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = "cn-hangzhou"
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_instance.instance.id
  source_endpoint_database_name      = "tfaccountpri_0"
  source_endpoint_user_name          = "tftestprivilege"
  source_endpoint_password           = "Test12345"
  db_list                            = <<EOF
        {"dtstestdata": {"name": "tfaccountpri_0", "all": true}}
    EOF
  subscription_instance_network_type = "vpc"
  subscription_instance_vpc_id       = data.alicloud_vpcs.default1.ids[0]
  subscription_instance_vswitch_id   = data.alicloud_vswitches.default_1.ids[0]
  status                             = "Normal"
}
```

## Argument Reference

The following arguments were support:

* `dts_instance_id` - (Computed, ForceNew) The ID of subscription instance.
* `dts_job_name` - (Required) The name of subscription task.
* `checkpoint` - (Optional, OtherParam) Subscription start time in Unix timestamp format.
* `compute_unit` - (Optional, OtherParam) [ETL specifications](https://help.aliyun.com/document_detail/212324.html). The unit is the computing unit ComputeUnit (CU), 1CU=1vCPU+4 GB memory. The value range is an integer greater than or equal to 2.
* `database_count` - (Optional, OtherParam) The number of private customized RDS instances under PolarDB-X. The default value is 1. This parameter needs to be passed only when `source_endpoint_engine_name` equals `drds`.
* `db_list` - (Required) Subscription object, in the format of JSON strings. For detailed definitions, please refer to the description of migration, synchronization or subscription objects [document](https://help.aliyun.com/document_detail/209545.html).
* `delay_notice` - (Optional) This parameter decides whether to monitor the delay status. Valid values: `true`, `false`.
* `delay_phone` - (Optional) The mobile phone number of the contact who delayed the alarm. Multiple mobile phone numbers separated by English commas `,`. This parameter currently only supports China stations, and only supports mainland mobile phone numbers, and up to 10 mobile phone numbers can be passed in.
* `delay_rule_time` - (Optional) When `delay_notice` is set to `true`, this parameter must be passed in. The threshold for triggering the delay alarm. The unit is second and needs to be an integer. The threshold can be set according to business needs. It is recommended to set it above 10 seconds to avoid delay fluctuations caused by network and database load.
* `destination_endpoint_engine_name` - (Optional) The destination endpoint engine name. Valid values: `ADS`, `DB2`, `DRDS`, `DataHub`, `Greenplum`, `MSSQL`, `MySQL`, `PolarDB`, `PostgreSQL`, `Redis`, `Tablestore`, `as400`, `clickhouse`, `kafka`, `mongodb`, `odps`, `oracle`, `polardb_o`, `polardb_pg`, `tidb`.
* `destination_region` - (Optional) The destination region. List of [supported regions](https://help.aliyun.com/document_detail/141033.html).
* `error_notice` - (Optional) This parameter decides whether to monitor abnormal status. Valid values: `true`, `false`.
* `error_phone` - (Optional) The mobile phone number of the contact for abnormal alarm. Multiple mobile phone numbers separated by English commas `,`. This parameter currently only supports China stations, and only supports mainland mobile phone numbers, and up to 10 mobile phone numbers can be passed in.
* `instance_class` - (Required) The instance class. Valid values: `large`, `medium`, `micro`, `small`, `xlarge`, `xxlarge`.
* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `Subscription`, `PayAsYouGo`.
* `payment_duration_unit` - (Optional) The payment duration unit. Valid values: `Month`, `Year`. When `payment_type` is `Subscription`, this parameter is valid and must be passed in.
* `payment_duration` - (Required when `payment_type` equals `Subscription`) The duration of prepaid instance purchase. When `payment_type` is `Subscription`, this parameter is valid and must be passed in.
* `reserve` - (Optional) DTS reserves parameters, the format is a JSON string, you can pass in this parameter to complete the source and target database information (such as the data storage format of the target Kafka database, the instance ID of the cloud enterprise network CEN). For more information, please refer to the parameter description of the [Reserve parameter](https://help.aliyun.com/document_detail/176470.html).
* `source_endpoint_database_name` - (Optional) To subscribe to the name of the database.
* `source_endpoint_engine_name` - (Required) The source database type value is MySQL or Oracle. Valid values: `MySQL`, `Oracle`.
* `source_endpoint_instance_id` - (Optional) The ID of source instance. Only when the type of source database instance was RDS MySQL, PolarDB-X 1.0, PolarDB MySQL, this parameter can be available and must be set.
* `source_endpoint_instance_type` - (Required, ForceNew) The type of source instance. Valid values: `RDS`, `PolarDB`, `DRDS`, `LocalInstance`, `ECS`, `Express`, `CEN`, `dg`.
* `source_endpoint_ip` - (Optional) The IP of source endpoint.
* `source_endpoint_oracle_sid` - (Optional) The SID of Oracle Database. When the source database is self-built Oracle and the Oracle database is a non-RAC instance, this parameter is available and must be passed in.
* `source_endpoint_owner_id` - (Optional) The Alibaba Cloud account ID to which the source instance belongs. This parameter is only available when configuring data subscriptions across Alibaba Cloud accounts and must be passed in.
* `source_endpoint_user_name` - (Optional) The username of source database instance account.
* `source_endpoint_password` - (Optional) The password of source database instance account.
* `source_endpoint_port` - (Optional) The port of source database.
* `source_endpoint_region` - (Required) The region of source database.
* `source_endpoint_role` - (Optional) Both the authorization roles. When the source instance and configure subscriptions task of the Alibaba Cloud account is not the same as the need to pass the parameter, to specify the source of the authorization roles, to allow configuration subscription task of the Alibaba Cloud account to access the source of the source instance information.
* `subscription_data_type_ddl` - (Optional, Computed) Whether to subscribe the DDL type of data. Valid values: `true`, `false`.
* `subscription_data_type_dml` - (Optional, Computed) Whether to subscribe the DML type of data. Valid values: `true`, `false`.
* `subscription_instance_network_type` - (Required) Subscription task type of network value: classic: classic Network. Virtual Private Cloud (vpc): a vpc. Valid values: `classic`, `vpc`.
* `subscription_instance_vpc_id` - (Optional) The ID of subscription vpc instance. When the value of `subscription_instance_network_type` is vpc, this parameter is available and must be passed in.
* `subscription_instance_vswitch_id` - (Optional) The ID of subscription VSwitch instance. When the value of `subscription_instance_network_type` is vpc, this parameter is available and must be passed in.
* `sync_architecture` - (Optional) The sync architecture. Valid values: `bidirectional`, `oneway`.
* `synchronization_direction` - (Optional) The synchronization direction. Valid values: `Forward`, `Reverse`. When the topology type of the data synchronization instance is bidirectional, it can be passed in to reverse to start the reverse synchronization link.
* `status` - (Optional, Computed) The status of the task. Valid values: `Normal`, `Abnormal`. When a task created, it is in this state of `NotStarted`. You can specify this state to `Normal` to start the job, and specify this state of `Abnormal` to stop the job. **Note: We treat the state `Starting` as the state of `Normal`, and consider the two states to be consistent on the user side.**
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Subscription Job.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `update` - (Defaults to 1 mins) Used when update the Subscription Job.

## Import

DTS Subscription Job can be imported using the id, e.g.

```
$ terraform import alicloud_dts_subscription_job.example <id>
```
