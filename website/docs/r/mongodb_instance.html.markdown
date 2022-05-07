---
subcategory: "MongoDB"
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
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  name              = "vpc-123456"
}

resource "alicloud_mongodb_instance" "example" {
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  vswitch_id          = alicloud_vswitch.default.id
  security_ip_list    = ["10.168.1.12", "100.69.7.112"]
}
```

## Module Support

You can use to the existing [mongodb module](https://registry.terraform.io/modules/terraform-alicloud-modules/mongodb/alicloud) 
to create a MongoDB instance resource one-click.


## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/61763.htm) `EngineVersion`.
* `db_instance_class` - (Required) Instance specification. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
* `db_instance_storage` - (Required) User-defined DB instance storage space.Unit: GB. Value range:
  - Custom storage space.
  - 10-GB increments. 
* `replication_factor` - (Optional) Number of replica set nodes. Valid values: [1, 3, 5, 7]
* `storage_engine` (Optional, ForceNew) Storage engine: WiredTiger or RocksDB. System Default value: WiredTiger.
* `name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `instance_charge_type` - (Optional) Valid values are `PrePaid`, `PostPaid`, System default to `PostPaid`. **NOTE:** It can be modified from `PostPaid` to `PrePaid` after version 1.63.0.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. System default to 1.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance. it supports multiple zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in one of them.
The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alicloud_zones`.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `account_password` - (Optional, Sensitive) Password of the root account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `kms_encrypted_password` - (Optional, Available in 1.57.1+) An KMS encrypts password used to a instance. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `security_ip_list` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `security_group_id` - (Optional, Available in 1.73.0+) The Security Group ID of ECS.
* `backup_period` - (Optional, Available in 1.42.0+) MongoDB Instance backup period. It is required when `backup_time` was existed. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]
* `backup_time` - (Optional, Available in 1.42.0+) MongoDB instance backup time. It is required when `backup_period` was existed. In the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. If not set, the system will return a default, like "23:00Z-24:00Z".
* `tde_status` - (Optional, Available in 1.73.0+) The TDE(Transparent Data Encryption) status.
* `maintain_start_time` - (Optional, Available in v1.56.0+) The start time of the operation and maintenance time period of the instance, in the format of HH:mmZ (UTC time).
* `maintain_end_time` - (Optional, Available in v1.56.0+) The end time of the operation and maintenance time period of the instance, in the format of HH:mmZ (UTC time).
* `order_type` - (Optional, Available in v1.134.0+) The type of configuration changes performed. Default value: DOWNGRADE. Valid values:
  * UPGRADE: The specifications are upgraded.
  * DOWNGRADE: The specifications are downgraded.
    Note: This parameter is only applicable to instances when `instance_charge_type` is PrePaid.
    
* `ssl_action` - (Optional, Available in v1.78.0+) Actions performed on SSL functions, Valid values: `Open`: turn on SSL encryption; `Close`: turn off SSL encryption; `Update`: update SSL certificate.
* `tags` - (Optional, Available in v1.66.0+) A mapping of tags to assign to the resource.
* `auto_renew` - (Optional, Available in v1.141.0+) Auto renew for prepaid, true of false. Default is false.
-> **NOTE:** The start time to the end time must be 1 hour. For example, the MaintainStartTime is 01:00Z, then the MaintainEndTime must be 02:00Z.
* `network_type` - (Optional, ForceNew, Computed, Available in v1.161.0+) The network type of the instance. Valid values:`Classic` or `VPC`. Default value: `Classic`.
* `vpc_id` - (Optional, ForceNew, Computed, Available in v1.161.0+) The ID of the VPC. -> **NOTE:** This parameter is valid only when NetworkType is set to VPC.
* `resource_group_id` - (Optional, Computed, Available in v1.161.0+) The ID of the Resource Group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MongoDB.
* `retention_period` - Instance log backup retention days. Available in 1.42.0+.
* `replica_set_name` - The name of the mongo replica set
* `ssl_status` - Status of the SSL feature. `Open`: SSL is turned on; `Closed`: SSL is turned off.
* `replica_sets` - Replica set instance information. The details see Block replica_sets. **NOTE:** Available in v1.140+.

#### replica_sets
The replica_sets supports the following:
* `vswitch_id` - The virtual switch ID to launch DB instances in one VPC.
* `connection_port` - The connection port of the node.
* `replica_set_role` - The role of the node. Valid values: `Primary`,`Secondary`.
* `connection_domain` - The connection address of the node.
* `vpc_cloud_instance_id` - VPC instance ID.
* `network_type` - The network type of the node. Valid values: `Classic`,`VPC`.
* `vpc_id` - The private network ID of the node.

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
