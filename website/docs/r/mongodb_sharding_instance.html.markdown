---
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_sharding_instance"
sidebar_current: "docs-alicloud-resource-mongodb-instance"
description: |-
  Provides a MongoDB sharding instance resource.
---

# alicloud\_mongodb\_sharding_instance

Provides a MongoDB sharding instance resource supports replica set instances only. the MongoDB provides stable, reliable, and automatic scalable database services. 
It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.
You can see detail product introduction [here](https://www.alibabacloud.com/help/doc-detail/26558.htm)

-> **NOTE:**  Available in 1.40.0+

-> **NOTE:**  The following regions don't support create Classic network MongoDB sharding instance.
[`cn-zhangjiakou`,`cn-huhehaote`,`ap-southeast-2`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`me-east-1`,`ap-northeast-1`,`eu-west-1`] 

-> **NOTE:**  Create MongoDB Sharding instance or change instance type and storage would cost 10~20 minutes. Please make full preparation

## Example Usage

### Create a Mongodb Sharding instance

```tf
variable "name" {
  default = "tf-example"
}

variable "shard" {
  default = {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
}

variable "mongo" {
  default = {
    node_class = "dds.mongos.mid"
  }
}

data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
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

resource "alicloud_mongodb_sharding_instance" "foo" {
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_id     = "${alicloud_vswitch.default.id}"
  engine_version = "3.4"
  name           = "${var.name}"
  shard_list     = ["${var.shard}", "${var.shard}"]
  mongo_list     = ["${var.mongo}", "${var.mongo}"]
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/zh/doc-detail/61884.htm) `EngineVersion`. 
* `storage_engine` (Optional, ForceNew) Storage engine: WiredTiger or RocksDB. System Default value: WiredTiger.
* `name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `instance_charge_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`,System default to `PostPaid`.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. System default to 1.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance. MongoDB sharding instance does not support multiple-zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in one of them.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `account_password` -  (Optional, Sensitive) Password of the root account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `kms_encrypted_password` - (Optional, Available in 1.57.1+) An KMS encrypts password used to a instance. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `security_ip_list` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]). System default to `["127.0.0.1"]`.
* `mongo_list` - (Required) The mongo-node count can be purchased is in range of [2, 32].
    * `node_class` -(Required) Node specification. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
* `shard_list` - (Required) the shard-node count can be purchased is in range of [2, 32].
    * `node_class` -(Required) Node specification. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
    * `node_storage` - (Required)
        - Custom storage space; value range: [10, 1,000]
        - 10-GB increments. Unit: GB.
* `backup_period` - (Optional, Available in 1.42.0+) MongoDB Instance backup period. It is required when `backup_time` was existed. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]
* `backup_time` - (Optional, Available in 1.42.0+) MongoDB instance backup time. It is required when `backup_period` was existed. In the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. Default to a random time, like "23:00Z-24:00Z".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MongoDB.
* `mongo_list`
    * `node_id` - The ID of the mongo-node.
    * `connect_string` - Mongo node connection string
    * `port` - Mongo node port
* `shard_list`
    * `node_id` - The ID of the shard-node.
* `retention_period` - Instance log backup retention days. Available in 1.42.0+.

## Import

MongoDB can be imported using the id, e.g.

```
$ terraform import alicloud_mongodb_sharding_instance.example dds-bp1291daeda44195
```
