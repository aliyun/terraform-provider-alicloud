---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_read_write_splitting_connection"
sidebar_current: "docs-alicloud-resource-db-read-write-splitting-connection"
description: |-
  Provides an RDS instance read write splitting connection resource.
---

# alicloud\_db\_read\_write\_splitting\_connection

Provides an RDS read write splitting connection resource to allocate an Intranet connection string for RDS instance.

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
}

resource "alicloud_db_read_write_splitting_connection" "foo" {
	instance_id = "${alicloud_db_instance.default.id}"
	connection_prefix = "t-con-${alicloud_db_instance.default.id}"
	distribution_type = "Custom"
	max_delay_time = 300
	weight = "${map(
		"${alicloud_db_instance.default.id}", "0",
		"${alicloud_db_readonly_instance.foo.id}", "500"
	)}"

	depends_on = ["alicloud_db_readonly_instance.foo"]
}
```

~> **NOTE:** Resource `alicloud_db_read_write_splitting_connection` should be created after `alicloud_db_readonly_instance`, so the `depends_on` statement is necessary.

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `distribution_type` - (Required) Read weight distribution mode. Values are as follows: `Standard` indicates automatic weight distribution based on types, `Custom` indicates custom weight distribution. 
* `connection_prefix` - (Optional, ForceNew) Prefix of an Internet connection string. It must be checked for uniqueness. It may consist of lowercase letters, numbers, and underlines, and must start with a letter and have no more than 30 characters. Default to <instance_id> + 'rw'.
* `port` - (Optional) Intranet connection port. Valid value: [3001-3999]. Default to 3306.
* `max_delay_time` - (Optional) Delay threshold, in seconds. The value range is 0 to 7200. Default to 30. Read requests are not routed to the read-only instances with a delay greater than the threshold.  
* `weight` - (Optional) Read weight distribution. Read weights increase at a step of 100 up to 10,000. Enter weights in the following format: {"Instanceid":"Weight","Instanceid":"Weight"}. This parameter must be set when distribution_type is set to Custom. 

## Attributes Reference

The following attributes are exported:

* `id` - The Id of DB instance.
* `connection_string` - Connection instance string.

## Import

RDS read write splitting connection can be imported using the id, e.g.

```
$ terraform import alicloud_db_read_write_splitting_connection.example abc12345678
```
