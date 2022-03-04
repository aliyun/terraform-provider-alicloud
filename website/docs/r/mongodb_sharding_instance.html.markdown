---
subcategory: "MongoDB"
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
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.default.zones[0].id
  name       = var.name
}

resource "alicloud_mongodb_sharding_instance" "foo" {
  zone_id        = data.alicloud_zones.default.zones[0].id
  vswitch_id     = alicloud_vswitch.default.id
  engine_version = "3.4"
  name           = var.name
  dynamic "shard_list" {
    for_each = [var.shard]
    content {
      # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
      # which keys might be set in maps assigned here, so it has
      # produced a comprehensive set here. Consider simplifying
      # this after confirming which keys can be set in practice.

      node_class   = shard_list.value.node_class
      node_storage = shard_list.value.node_storage
    }
  }
  dynamic "shard_list" {
    for_each = [var.shard]
    content {
      # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
      # which keys might be set in maps assigned here, so it has
      # produced a comprehensive set here. Consider simplifying
      # this after confirming which keys can be set in practice.

      node_class   = shard_list.value.node_class
      node_storage = shard_list.value.node_storage
    }
  }
  dynamic "mongo_list" {
    for_each = [var.mongo]
    content {
      # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
      # which keys might be set in maps assigned here, so it has
      # produced a comprehensive set here. Consider simplifying
      # this after confirming which keys can be set in practice.

      node_class = mongo_list.value.node_class
    }
  }
  dynamic "mongo_list" {
    for_each = [var.mongo]
    content {
      # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
      # which keys might be set in maps assigned here, so it has
      # produced a comprehensive set here. Consider simplifying
      # this after confirming which keys can be set in practice.

      node_class = mongo_list.value.node_class
    }
  }
}
```

## Module Support

You can use to the existing [mongodb-sharding module](https://registry.terraform.io/modules/terraform-alicloud-modules/mongodb-sharding/alicloud) 
to create a MongoDB sharding instance resource one-click.

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/en/doc-detail/61884.htm) `EngineVersion`. 
* `storage_engine` (Optional, ForceNew) Storage engine: WiredTiger or RocksDB. System Default value: WiredTiger.
* `name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `instance_charge_type` - (Optional) Valid values are `PrePaid`, `PostPaid`,System default to `PostPaid`. **NOTE:** It can be modified from `PostPaid` to `PrePaid` after version v1.141.0.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. System default to 1.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance. MongoDB sharding instance does not support multiple-zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in one of them.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `account_password` - (Optional, Sensitive) Password of the root account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `kms_encrypted_password` - (Optional, Available in 1.57.1+) An KMS encrypts password used to a instance. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `security_ip_list` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]). System default to `["127.0.0.1"]`.
* `security_group_id` - (Optional, Available in 1.76.0+) The Security Group ID of ECS.
* `tde_status` - (Optional, Available in 1.76.0+) The TDE(Transparent Data Encryption) status. It can be updated from version 1.160.0+.
* `mongo_list` - (Required) The mongo-node count can be purchased is in range of [2, 32].
    * `node_class` -(Required) Node specification. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
* `shard_list` - (Required) the shard-node count can be purchased is in range of [2, 32].
    * `node_class` -(Required) Node specification. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
    * `node_storage` - (Required)
        - Custom storage space; value range: [10, 1,000]
        - 10-GB increments. Unit: GB.
    * `readonly_replicas` - (Optional, Available in 1.126.0+) The number of read-only nodes in shard node. Valid values: 0 to 5. Default value: 0.
* `backup_period` - (Optional, Available in 1.42.0+) MongoDB Instance backup period. It is required when `backup_time` was existed. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]
* `backup_time` - (Optional, Available in 1.42.0+) MongoDB instance backup time. It is required when `backup_period` was existed. In the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. If not set, the system will return a default, like "23:00Z-24:00Z".
* `order_type` - (Optional, Available in v1.134.0+) The type of configuration changes performed. Default value: DOWNGRADE. Valid values:
  * UPGRADE: The specifications are upgraded.
  * DOWNGRADE: The specifications are downgraded. 
    Note: This parameter is only applicable to instances when `instance_charge_type` is PrePaid.
* `auto_renew` - (Optional, Available in v1.141.0+) Auto renew for prepaid, true of false. Default is false.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `network_type` - (Optional, ForceNew, Computed, Available in v1.161.0+) The network type of the instance. Valid values:`Classic` or `VPC`. Default value: `Classic`.
* `vpc_id` - (Optional, ForceNew, Computed, Available in v1.161.0+) The ID of the VPC. -> **NOTE:** This parameter is valid only when NetworkType is set to VPC.
* `protocol_type` - (Optional, ForceNew, Computed, Available in v1.161.0+) The type of the access protocol. Valid values: `mongodb` or `dynamodb`.
* `resource_group_id` - (Optional, Computed, Available in v1.161.0+) The ID of the Resource Group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MongoDB.
* `mongo_list`
    * `node_id` - The ID of the mongo-node.
    * `connect_string` - Mongo node connection string
    * `port` - Mongo node port
* `shard_list`
    * `node_id` - The ID of the shard-node.
* `retention_period` - Instance log backup retention days. **NOTE:** Available in 1.42.0+.
* `config_server_list` - The node information list of config server. The details see Block `config_server_list`. **NOTE:** Available in v1.140+.

#### config_server_list
The config_server_list supports the following:
* `max_iops` - The maximum IOPS of the Config Server node.
* `connect_string` - The connection address of the Config Server node.
* `node_class` - The node class of the Config Server node.
* `max_connections` - The max connections of the Config Server node.
* `port` - The connection port of the Config Server node.
* `node_description` - The description of the Config Server node.
* `node_id` - The ID of the Config Server node.
* `node_storage` - The node storage of the Config Server node.

### Timeouts

-> **NOTE:** Available in 1.126.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the MongoDB instance (until it reaches the initial `Running` status).
* `update` - (Defaults to 30 mins) Used when updating the MongoDB instance (until it reaches the initial `Running` status).
* `delete` - (Defaults to 30 mins) Used when terminating the MongoDB instance.

## Import

MongoDB can be imported using the id, e.g.

```
$ terraform import alicloud_mongodb_sharding_instance.example dds-bp1291daeda44195
```
