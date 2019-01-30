---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instance"
sidebar_current: "docs-alicloud-resource-db-instance"
description: |-
  Provides an RDS instance resource.
---

# alicloud\_db\_instance

Provides an RDS instance resource. A DB instance is an isolated database
environment in the cloud. A DB instance can contain multiple user-created
databases.

## Example Usage

```
resource "alicloud_db_instance" "default" {
	engine = "MySQL"
	engine_version = "5.6"
	db_instance_class = "rds.mysql.t1.small"
	db_instance_storage = "10"
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required) Database type. Value options: MySQL, SQLServer, PostgreSQL, and PPAS.
* `engine_version` - (Required) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `db_instance_class` - (Deprecated) It has been deprecated from version 1.5.0 and use 'instance_type' to replace.
* `instance_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `db_instance_storage` - (Deprecated) It has been deprecated from version 1.5.0 and use 'instance_storage' to replace.
* `instance_storage` - (Required) User-defined DB instance storage space. Value range:
    - [5, 2000] for MySQL/PostgreSQL/PPAS HA dual node edition;
    - [20,1000] for MySQL 5.7 basic single node edition;
    - [10, 2000] for SQL Server 2008R2;
    - [20,2000] for SQL Server 2012 basic single node edition
    Increase progressively at a rate of 5 GB. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).

* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `instance_charge_type` - (Optional) Valid values are `Prepaid`, `Postpaid`, Default to `Postpaid`.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. Default to 1.
* `zone_id` - (Optional) The Zone to launch the DB instance. From version 1.8.1, it supports multiple zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in the one of them.
The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alicloud_zones`.
* `multi_az` - (Optional) It has been deprecated from version 1.8.1, and `zone_id` can support multiple zone.
* `db_instance_net_type` - (Deprecated) It has been deprecated from version 1.5.0. If you want to set public connection, please use new resource `alicloud_db_connection`. Default to Intranet.
* `allocate_public_connection` - (Deprecated) It has been deprecated from version 1.5.0. If you want to allocate public connection string, please use new resource `alicloud_db_connection`.
* `instance_network_type` - (Deprecated) It has been deprecated from version 1.5.0. If you want to create instances in VPC network, this parameter must be set.
* `vswitch_id` - (Optional) The virtual switch ID to launch DB instances in one VPC.
* `master_user_name` - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_account` field 'name' replaces it.
* `master_user_password`  - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_account` field 'password' replaces it.
* `preferred_backup_period`  - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_backup_policy` field 'backup_period' replaces it.
* `preferred_backup_time` - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_backup_policy` field 'backup_time' replaces it.
* `backup_retention_period` - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_backup_policy` field 'retention_period' replaces it.
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `db_mappings` - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_database` replaces it.
* `parameters` - (Optional) The parameters needs to be set after DB instance was launched.

~> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `instance_charge_type` - The instance charge type.
* `period` - The DB instance using duration.
* `engine` - Database type.
* `engine_version` - The database engine version.
* `db_instance_class` - (Deprecated from version 1.5.0)
* `instance_type` - The RDS instance type.
* `db_instance_storage` - (Deprecated from version 1.5.0)
* `instance_storage` - The RDS instance storage space.
* `instance_name` - The name of DB instance.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.
* `zone_id` - The zone ID of the RDS instance.
* `db_instance_net_type` - (Deprecated from version 1.5.0).
* `instance_network_type` - (Deprecated from version 1.5.0).
* `db_mappings` - - (Deprecated from version 1.5.0).
* `preferred_backup_period` - (Deprecated from version 1.5.0).
* `preferred_backup_time` - (Deprecated from version 1.5.0).
* `backup_retention_period` - (Deprecated from version 1.5.0).
* `security_ips` - Security ips of instance whitelist.
* `connections` - (Deprecated from version 1.5.0).
* `vswitch_id` - If the rds instance created in VPC, then this value is virtual switch ID.
* `master_user_name` - (Deprecated from version 1.5.0).
* `preferred_backup_period` - (Deprecated from version 1.5.0).
* `preferred_backup_time` - (Deprecated from version 1.5.0).
* `backup_retention_period` - (Deprecated from version 1.5.0).
* `parameters` - Database parameters modified.

## Import

RDS instance can be imported using the id, e.g.

```
$ terraform import alicloud_db_instance.example rm-abc12345678
```