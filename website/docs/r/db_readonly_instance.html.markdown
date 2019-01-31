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
resource "alicloud_db_instance" "default" {
	engine = "MySQL"
	engine_version = "5.6"
	db_instance_class = "rds.mysql.t1.small"
	db_instance_storage = "10"
}

resource "alicloud_db_readonly_instance" "foo" {
	master_db_instance_id = "${alicloud_db_instance.default.id}"
	zone_id = "${alicloud_db_instance.default.zone_id}"
	engine_version = "${alicloud_db_instance.default.engine_version}"
	instance_type = "${alicloud_db_instance.default.instance_type}"
	instance_storage = "30"
	instance_name = "${var.name}ro"
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `zone_id` - (Required) The Zone to launch the DB instance.
* `master_db_instance_id` - (Required) ID of the master instance.
* `instance_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_storage` - (Required) User-defined DB instance storage space. Value range:
    - [5, 2000] for MySQL/PostgreSQL/PPAS HA dual node edition;
    Increase progressively at a rate of 5 GB. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).

* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `vswitch_id` - (Optional) The virtual switch ID to launch DB instances in one VPC.

~> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `engine` - Database type.
* `engine_version` - The database engine version.
* `instance_type` - The RDS instance type.
* `instance_storage` - The RDS instance storage space.
* `instance_name` - The name of DB instance.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.
* `zone_id` - The zone ID of the RDS instance.
* `vswitch_id` - If the rds instance created in VPC, then this value is virtual switch ID.

## Import

RDS instance can be imported using the id, e.g.

```
$ terraform import alicloud_db_readonly_instance.example rm-abc12345678
```