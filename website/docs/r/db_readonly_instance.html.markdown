---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_readonly_instance"
sidebar_current: "docs-alicloud-resource-db-readonly-instance"
description: |-
  Provides an RDS readonly instance resource.
---

# alicloud\_db\_readonly\_instance

Provides an RDS readonly instance resource. 

## Example Usage

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

resource "alicloud_db_readonly_instance" "foo" {
	master_db_instance_id = "${alicloud_db_instance.default.id}"
	engine_version = "${alicloud_db_instance.default.engine_version}"
	instance_type = "${alicloud_db_instance.default.instance_type}"
	instance_storage = "30"
	instance_name = "rm-123456_ro"
	vswitch_id = "${alicloud_vswitch.default.id}"
	parameters = [{
		name = "innodb_large_prefix"
		value = "ON"
	},{
		name = "connect_timeout"
		value = "50"
	}]
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `master_db_instance_id` - (Required) ID of the master instance.
* `instance_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_storage` - (Required) User-defined DB instance storage space. Value range: [5, 2000] for MySQL/SQL Server HA dual node edition. Increase progressively at a rate of 5 GB. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm).
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.

-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `engine` - Database type.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.

## Import

RDS readonly instance can be imported using the id, e.g.

```
$ terraform import alicloud_db_readonly_instance.example rm-abc12345678
```
