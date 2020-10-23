---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instance"
sidebar_current: "docs-alicloud-resource-db-instance"
description: |-
  Provides an RDS instance resource.
---

# alicloud\_db\_instance

Provides an RDS instance resource. A DB instance is an isolated database environment in the cloud. A DB instance can contain multiple user-created databases.

For information about RDS and how to use it, see [What is ApsaraDB for RDS](https://www.alibabacloud.com/help/en/doc-detail/26092.htm).

## Example Usage

### Create a RDS MySQL instance

```terraform
variable "name" {
  default = "tf-testaccdbinstance"
}

variable "creation" {
  default = "Rds"
}

data "alicloud_zones" "example" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "example" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id            = alicloud_vpc.example.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alicloud_zones.example.zones[0].id
  name              = var.name
}

resource "alicloud_db_instance" "example" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  vswitch_id           = alicloud_vswitch.example.id
  monitoring_period    = "60"
}
```

### Create a RDS MySQL instance with specific parameters

```terraform
resource "alicloud_vpc" "example" {
  name       = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id            = alicloud_vpc.example.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alicloud_zones.example.zones.0.id
  name              = "vpc-123456"
}

resource "alicloud_db_instance" "default" {
  engine              = "MySQL"
  engine_version      = "5.6"
  db_instance_class   = "rds.mysql.t1.small"
  db_instance_storage = "10"
  vswitch_id          = alicloud_vswitch.example.id
}

resource "alicloud_db_instance" "example" {
  engine              = "MySQL"
  engine_version      = "5.6"
  db_instance_class   = "rds.mysql.t1.small"
  db_instance_storage = "10"
  parameters {
    name  = "innodb_large_prefix"
    value = "ON"
  }
  parameters {
    name  = "connect_timeout"
    value = "50"
  }
}
```
### Create a High Availability RDS MySQL Instance

```terraform
variable "name" {
  default = "tf-testaccdbinstance"
}

data "alicloud_zones" "example" {
  available_resource_creation = "Rds"
}

resource "alicloud_vpc" "example" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  count             = 2
  vpc_id            = alicloud_vpc.example.id
  cidr_block        = format("172.16.%d.0/24", count.index+1)
  availability_zone = data.alicloud_zones.example.zones[count.index].id
  name              = format("vswich_%d", var.name, count.index)
}

resource "alicloud_db_instance" "example" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_storage     = "30"
  instance_type        = "rds.mysql.t1.small"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  zone_id              = data.alicloud_zones.example.zones.0.id
  zone_id_slave_a      = data.alicloud_zones.example.zones.1.id
  vswitch_id           = join(",", alicloud_vswitch.example.*.id)
  monitoring_period    = "60"
}

```

### Create a High Availability RDS MySQL Instance with multi zones

```terraform
variable "name" {
  default = "tf-testaccdbinstance"
}

data "alicloud_zones" "example" {
  available_resource_creation = "Rds"
  multi                       = true
}

resource "alicloud_vpc" "example" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  count             = length(data.alicloud_zones.example.zones.0.multi_zone_ids)
  vpc_id            = alicloud_vpc.example.id
  cidr_block        = format("172.16.%d.0/24", count.index+1)
  availability_zone = data.alicloud_zones.example.zones.0.multi_zone_ids[count.index]
  name              = format("vswitch_%d", count.index)
}

resource "alicloud_db_instance" "this" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_storage     = "30"
  instance_type        = "rds.mysql.t1.small"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  zone_id              = data.alicloud_zones.example.zones.0.id
  vswitch_id           = join(",", alicloud_vswitch.example.*.id)
  monitoring_period    = "60"
}
```

### Create a Enterprise Edition RDS MySQL Instance 
```terraform
variable "name" {
  default = "tf-testaccdbinstance"
}

data "alicloud_zones" "example" {
  available_resource_creation = "Rds"
}

resource "alicloud_vpc" "example" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  count             = 3
  vpc_id            = alicloud_vpc.example.id
  cidr_block        = format("172.16.%d.0/24", count.index+1)
  availability_zone = data.alicloud_zones.example.zones[count.index].id
  name              = format("vswich_%d", var.name, count.index)
}

resource "alicloud_db_instance" "example" {
  engine               = "MySQL"
  engine_version       = "8.0"
  instance_storage     = "30"
  instance_type        = "mysql.n2.small.25"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  zone_id              = data.alicloud_zones.example.zones.0.id
  zone_id_slave_a      = data.alicloud_zones.example.zones.1.id
  zone_id_slave_b      = data.alicloud_zones.example.zones.2.id
  vswitch_id           = join(",", alicloud_vswitch.example.*.id)
  monitoring_period    = "60"
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required,ForceNew) Database type. Value options: MySQL, SQLServer, PostgreSQL, and PPAS.
* `engine_version` - (Required,ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `instance_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_storage` - (Required) User-defined DB instance storage space. Value range:
    - [5, 2000] for MySQL/PostgreSQL/PPAS HA dual node edition;
    - [20,1000] for MySQL 5.7 basic single node edition;
    - [10, 2000] for SQL Server 2008R2;
    - [20,2000] for SQL Server 2012 basic single node edition
    Increase progressively at a rate of 5 GB. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
    Note: There is extra 5 GB storage for SQL Server Instance and it is not in specified `instance_storage`.

* `db_instance_storage_type` - (Optional, Available in 1.68.0+) The storage type of the instance. Valid values:
    - local_ssd: specifies to use local SSDs. This value is recommended.
    - cloud_ssd: specifies to use standard SSDs.
    - cloud_essd: specifies to use enhanced SSDs (ESSDs).
    - cloud_essd2: specifies to use enhanced SSDs (ESSDs).
    - cloud_essd3: specifies to use enhanced SSDs (ESSDs).

* `sql_collector_status` - (Optional, Available in 1.70.0+) The sql collector status of the instance. Valid values are `Enabled`, `Disabled`, Default to `Disabled`.
* `sql_collector_config_value` - (Optional, Available in 1.70.0+) The sql collector keep time of the instance. Valid values are `30`, `180`, `365`, `1095`, `1825`, Default to `30`.
    
* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `instance_charge_type` - (Optional) Valid values are `Prepaid`, `Postpaid`, Default to `Postpaid`. Currently, the resource only supports PostPaid to PrePaid.
* `resource_group_id` (Optional, ForceNew, Computed, Available in 1.86.0+) The ID of resource group which the DB instance belongs.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. Default to 1.
* `monitoring_period` - (Optional) The monitoring frequency in seconds. Valid values are 5, 60, 300. Defaults to 300. 
* `auto_renew` - (Optional, Available in 1.34.0+) Whether to renewal a DB instance automatically or not. It is valid when instance_charge_type is `PrePaid`. Default to `false`.
* `auto_renew_period` - (Optional, Available in 1.34.0+) Auto-renewal period of an instance, in the unit of the month. It is valid when instance_charge_type is `PrePaid`. Valid value:[1~12], Default to 1.
* `zone_id` - (ForceNew) The Zone to launch the DB instance. From version 1.8.1, it supports multiple zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in the one of them.
The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alicloud_zones`.
* `vswitch_id` - (ForceNew) The virtual switch ID to launch DB instances in one VPC. If there are multiple vswitches, separate them with commas.
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `security_ip_mode` - (Optional, Available in 1.62.1+)  Valid values are `normal`, `safety`, Default to `normal`. support `safety` switch to high security access mode 
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm) .
* `force_restart` - (Optional, Available in 1.75.0+) Set it to true to make some parameter efficient when modifying them. Default to false.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

   Note: From 1.63.0, the tag key and value are case sensitive. Before that, they are not case sensitive.

* `security_group_id` - (Deprecated) It has been deprecated from 1.69.0 and use `security_group_ids` instead.
* `security_group_ids` - (Optional, List(string), Available in 1.69.0+) The list IDs to join ECS Security Group. At most supports three security groups.
* `maintain_time` - (Optional, Available in 1.56.0+) Maintainable time period format of the instance: HH:MMZ-HH:MMZ (UTC time)
* `auto_upgrade_minor_version` - (Optional, Available in 1.62.1+) The upgrade method to use. Valid values:
   - Auto: Instances are automatically upgraded to a higher minor version.
   - Manual: Instances are forcibly upgraded to a higher minor version when the current version is unpublished.
   
   Default to "Manual". See more [details and limitation](https://www.alibabacloud.com/help/doc-detail/123605.htm).
   
* `zone_id_slave_a` - (Optional, ForceNew, Available in 1.101.0+) The region ID of the secondary instance if you create a secondary instance. If you set this parameter to the same value as the ZoneId parameter, the instance is deployed in a single zone. Otherwise, the instance is deployed in multiple zones.
* `zone_id_slave_b`- (Optional, ForceNew, Available in 1.101.0+) The region ID of the log instance if you create a log instance. If you set this parameter to the same value as the ZoneId parameter, the instance is deployed in a single zone. Otherwise, the instance is deployed in multiple zones.
* `ssl_action` - (Optional, Available in v1.90.0+) Actions performed on SSL functions, Valid values: `Open`: turn on SSL encryption; `Close`: turn off SSL encryption; `Update`: update SSL certificate. See more [engine and engineVersion limitation](https://www.alibabacloud.com/help/zh/doc-detail/26254.htm).
* `tde_status` - (Optional, ForceNew, Available in 1.90.0+) The TDE(Transparent Data Encryption) status. See more [engine and engineVersion limitation](https://www.alibabacloud.com/help/zh/doc-detail/26256.htm).

-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

-> **NOTE:** `zone_id_slave_a` and `zone_id_slave_b` can specify slave zone ids when creating the high-availability or enterprise edition instances. Meanwhile, `vswitch_id` needs to pass in the corresponding vswitch id to the slave zone by order (If the `vswitch_id` is not specified, the classic network version will be created). For example, `zone_id` = "zone-a" and `zone_id_slave_a` = "zone-c", `zone_id_slave_b` = "zone-b", then the `vswitch_id` must be "vsw-zone-a,vsw-zone-c,vsw-zone-b". Of course, you can also choose automatic allocation , for example, `zone_id` = "zone-a" and `zone_id_slave_a` = "Auto",`zone_id_slave_b` = "Auto", then the `vswitch_id` must be "vsw-zone-a,Auto,Auto". The list contains up to 2 slave zone ids , separated by commas.


## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.
* `ssl_status` - Status of the SSL feature. `Yes`: SSL is turned on; `No`: SSL is turned off.


### Timeouts

-> **NOTE:** Available in 1.52.1+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the db instance (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the db instance (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 20 mins) Used when terminating the db instance. 

## Import

RDS instance can be imported using the id, e.g.

```
$ terraform import alicloud_db_instance.example rm-abc12345678
```
