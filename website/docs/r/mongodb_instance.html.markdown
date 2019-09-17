---
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_instance"
sidebar_current: "docs-alicloud-resource-mongodb-instance"
description: |-
  Provides a MongoDB instance resource.
---

# alicloud\_mongodb\_instance

Provides a MongoDB instance resource supports replica set instances only. the MongoDB provides stable, reliable, and automatic scalable database services. 
It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.
You can see detail product introduction [here](https://www.alibabacloud.com/help/doc-detail/26558.htm)

-> **NOTE:**  Available in 1.37.0+

-> **NOTE:**  The following regions don't support create Classic network MongoDB instance.
[`cn-zhangjiakou`,`cn-huhehaote`,`ap-southeast-2`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`me-east-1`,`ap-northeast-1`,`eu-west-1`] 

-> **NOTE:**  Create MongoDB instance or change instance type and storage would cost 5~10 minutes. Please make full preparation

## Example Usage

### Create a Mongodb instance

```
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}

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

resource "alicloud_mongodb_instance" "example" {
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  vswitch_id          = "${alicloud_vswitch.default.id}"
  security_ip_list    = ["10.168.1.12", "100.69.7.112"]
}

```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/61763.htm) `EngineVersion`.
* `db_instance_class` - (Required) Instance specification. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
* `db_instance_storage` - (Required) User-defined DB instance storage space.Unit: GB. Value range:
  - Custom storage space; value range: [10,2000]
  - 10-GB increments. 
* `replication_factor` - (Optional) Number of replica set nodes. Valid values: [3,5,7]
* `storage_engine` (Optional, ForceNew) Storage engine: WiredTiger or RocksDB. System Default value: WiredTiger.
* `name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `instance_charge_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`,System default to `PostPaid`.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. System default to 1.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance. it supports multiple zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in one of them.
The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alicloud_zones`.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `account_password` -  (Optional, Sensitive) Password of the root account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `security_ip_list` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `backup_period` - (Optional, Available in 1.42.0+) MongoDB Instance backup period. It is required when `backup_time` was existed. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]
* `backup_time` - (Optional, Available in 1.42.0+) MongoDB instance backup time. It is required when `backup_period` was existed. In the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. Default to a random time, like "23:00Z-24:00Z".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MongoDB.
* `retention_period` - Instance log backup retention days. Available in 1.42.0+.

### Timeouts

-> **NOTE:** Available in 1.53.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the MongoDB instance (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the MongoDB instance (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 30 mins) Used when terminating the MongoDB instance. 

## Import

MongoDB can be imported using the id, e.g.

```
$ terraform import alicloud_mongodb_instance.example dds-bp1291daeda44194
```
