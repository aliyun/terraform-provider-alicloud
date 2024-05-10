---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_read_write_splitting_connection"
sidebar_current: "docs-alicloud-resource-db-read-write-splitting-connection"
description: |-
  Provides an RDS instance read write splitting connection resource.
---

# alicloud_db_read_write_splitting_connection

Provides an RDS read write splitting connection resource to allocate an Intranet connection string for RDS instance, see [What is DB Read Write Splitting Connection](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/api-rds-2014-08-15-allocatereadwritesplittingconnection).

-> **NOTE:** Available since v1.48.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_db_zones" "example" {
  engine         = "MySQL"
  engine_version = "5.6"
}
data "alicloud_db_instance_classes" "example" {
  zone_id        = data.alicloud_db_zones.example.ids.0
  engine         = "MySQL"
  engine_version = "5.6"
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
  engine               = "MySQL"
  engine_version       = "5.6"
  category             = "HighAvailability"
  instance_type        = data.alicloud_db_instance_classes.example.instance_classes.1.instance_class
  instance_storage     = "20"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  vswitch_id           = alicloud_vswitch.example.id
  security_ips         = ["10.168.1.12", "100.69.7.112"]
}

resource "alicloud_db_readonly_instance" "example" {
  zone_id               = alicloud_db_instance.example.zone_id
  master_db_instance_id = alicloud_db_instance.example.id
  engine_version        = alicloud_db_instance.example.engine_version
  instance_storage      = alicloud_db_instance.example.instance_storage
  instance_type         = alicloud_db_instance.example.instance_type
  instance_name         = "${var.name}readonly"
  vswitch_id            = alicloud_vswitch.example.id
}

resource "alicloud_db_read_write_splitting_connection" "example" {
  instance_id       = alicloud_db_readonly_instance.example.master_db_instance_id
  connection_prefix = "example-con-123"
  distribution_type = "Standard"
}
```

-> **NOTE:** Resource `alicloud_db_read_write_splitting_connection` should be created after `alicloud_db_readonly_instance`, so the `depends_on` statement is necessary.

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

```shell
$ terraform import alicloud_db_read_write_splitting_connection.example abc12345678
```
