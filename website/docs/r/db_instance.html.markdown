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

### Create a RDS MySQL instance

```
variable "name" {
	default = "dbInstanceconfig"
}
variable "creation" {
		default = "Rds"
}
data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.s2.large"
	instance_storage = "30"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}
```

### Create a RDS MySQL instance with specific parameters

```
resource "alicloud_vpc" "default" {
	name       = "vpc-123456"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
	vpc_id            = "${alicloud_vpc.default.id}"
	cidr_block        = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name              = "vpc-123456"
}

resource "alicloud_db_instance" "default" {
	engine = "MySQL"
	engine_version = "5.6"
	db_instance_class = "rds.mysql.t1.small"
	db_instance_storage = "10"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_db_instance" "default" {
	engine = "MySQL"
	engine_version = "5.6"
	db_instance_class = "rds.mysql.t1.small"
	db_instance_storage = "10"
	parameters {
		name = "innodb_large_prefix"
		value = "ON"
	}
	parameters {
		name = "connect_timeout"
		value = "50"
	}
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required,ForceNew) Database type. Value options: MySQL, SQLServer, PostgreSQL, and PPAS.
* `engine_version` - (Required,ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `db_instance_class` - (Deprecated) It has been deprecated from version 1.5.0 and use 'instance_type' to replace.
* `instance_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `db_instance_storage` - (Deprecated) It has been deprecated from version 1.5.0 and use 'instance_storage' to replace.
* `instance_storage` - (Required) User-defined DB instance storage space. Value range:
    - [5, 2000] for MySQL/PostgreSQL/PPAS HA dual node edition;
    - [20,1000] for MySQL 5.7 basic single node edition;
    - [10, 2000] for SQL Server 2008R2;
    - [20,2000] for SQL Server 2012 basic single node edition
    Increase progressively at a rate of 5 GB. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
    Note: There is extra 5 GB storage for SQL Server Instance and it is not in specified `instance_storage`.

* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `instance_charge_type` - (Optional) Valid values are `Prepaid`, `Postpaid`, Default to `Postpaid`. Currently, the resource only supports PostPaid to PrePaid.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. Default to 1.
* `monitoring_period` - (Optional) The monitoring frequency in seconds. Valid values are 5, 60, 300. Defaults to 300. 
* `auto_renew` - (Optional, Available in 1.34.0+) Whether to renewal a DB instance automatically or not. It is valid when instance_charge_type is `PrePaid`. Default to `false`.
* `auto_renew_period` - (Optional, Available in 1.34.0+) Auto-renewal period of an instance, in the unit of the month. It is valid when instance_charge_type is `PrePaid`. Valid value:[1~12], Default to 1.
* `zone_id` - (ForceNew) The Zone to launch the DB instance. From version 1.8.1, it supports multiple zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in the one of them.
The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alicloud_zones`.
* `multi_az` - (Optional) It has been deprecated from version 1.8.1, and `zone_id` can support multiple zone.
* `db_instance_net_type` - (Deprecated) It has been deprecated from version 1.5.0. If you want to set public connection, please use new resource `alicloud_db_connection`. Default to Intranet.
* `allocate_public_connection` - (Deprecated) It has been deprecated from version 1.5.0. If you want to allocate public connection string, please use new resource `alicloud_db_connection`.
* `instance_network_type` - (Deprecated) It has been deprecated from version 1.5.0. If you want to create instances in VPC network, this parameter must be set.
* `vswitch_id` - (ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `master_user_name` - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_account` field 'name' replaces it.
* `master_user_password`  - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_account` field 'password' replaces it.
* `preferred_backup_period`  - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_backup_policy` field 'backup_period' replaces it.
* `preferred_backup_time` - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_backup_policy` field 'backup_time' replaces it.
* `backup_retention_period` - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_backup_policy` field 'retention_period' replaces it.
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `db_mappings` - (Deprecated) It has been deprecated from version 1.5.0. New resource `alicloud_db_database` replaces it.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm) .
* `tags` - (Optional) the instance bound to the tag. The format of the incoming value is `json` string, including `TagKey` and `TagValue`. `TagKey` cannot be null, and `TagValue` can be empty, and both cannot begin with `aliyun`. Format example `{"key1":"value1"}`.
* `security_group_id` - (Optional) Input the ECS Security Group ID to join ECS Security Group. Only support mysql 5.5, mysql 5.6

-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.

## Import

RDS instance can be imported using the id, e.g.

```
$ terraform import alicloud_db_instance.example rm-abc12345678
```