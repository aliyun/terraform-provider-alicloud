---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_sharding_instance"
sidebar_current: "docs-alicloud-resource-mongodb-sharding-instance"
description: |-
  Provides a MongoDB sharding instance resource.
---

# alicloud_mongodb_sharding_instance

Provides a MongoDB Sharding Instance resource supports replica set instances only. the MongoDB provides stable, reliable, and automatic scalable database services.
It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.
You can see detail product introduction [here](https://www.alibabacloud.com/help/doc-detail/26558.htm)

-> **NOTE:** Available since v1.40.0.

-> **NOTE:**  The following regions don't support create Classic network MongoDB Sharding Instance.
[`cn-zhangjiakou`,`cn-huhehaote`,`ap-southeast-3`,`ap-southeast-5`,`me-east-1`,`ap-northeast-1`,`eu-west-1`]

-> **NOTE:**  Create MongoDB Sharding instance or change instance type and storage would cost 10~20 minutes. Please make full preparation.

## Example Usage

### Create a Mongodb Sharding instance

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_mongodb_sharding_instance&exampleId=4a2adea6-e53f-129f-cdc7-a9a46d5c3bc985009ca2&activeTab=example&spm=docs.r.mongodb_sharding_instance.0.4a2adea6e5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_mongodb_zones" "default" {
}

locals {
  index   = length(data.alicloud_mongodb_zones.default.zones) - 1
  zone_id = data.alicloud_mongodb_zones.default.zones[local.index].id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = local.zone_id
}

resource "alicloud_mongodb_sharding_instance" "default" {
  engine_version = "4.2"
  vswitch_id     = alicloud_vswitch.default.id
  zone_id        = local.zone_id
  name           = var.name
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = "10"
  }
  shard_list {
    node_class        = "dds.shard.standard"
    node_storage      = "20"
    readonly_replicas = "1"
  }
}
```

## Module Support

You can use to the existing [mongodb-sharding module](https://registry.terraform.io/modules/terraform-alicloud-modules/mongodb-sharding/alicloud)
to create a MongoDB Sharding Instance resource one-click.

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/en/doc-detail/61884.htm) `EngineVersion`. **NOTE:** From version 1.225.1, `engine_version` can be modified.
* `storage_engine` (Optional, ForceNew) The storage engine of the instance. Default value: `WiredTiger`. Valid values: `WiredTiger`, `RocksDB`.
* `storage_type` - (Optional, Available since v1.225.1) The storage type of the instance. Valid values: `cloud_essd1`, `cloud_essd2`, `cloud_essd3`, `cloud_auto`, `local_ssd`. **NOTE:** From version 1.229.0, `storage_type` can be modified. However, `storage_type` can only be modified to `cloud_auto`.
* `provisioned_iops` - (Optional, Int, Available since v1.229.0) The provisioned IOPS. Valid values: `0` to `50000`.
* `protocol_type` - (Optional, ForceNew, Available since v1.161.0) The type of the access protocol. Valid values: `mongodb` or `dynamodb`.
* `vpc_id` - (Optional, ForceNew, Available since v1.161.0) The ID of the VPC. -> **NOTE:** `vpc_id` is valid only when `network_type` is set to `VPC`.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance. MongoDB Sharding Instance does not support multiple-zone.
  If it is a multi-zone and `vswitch_id` is specified, the vswitch must in one of them.
* `security_group_id` - (Optional, Available since v1.76.0) The Security Group ID of ECS.
* `network_type` - (Optional, ForceNew, Available since v1.161.0) The network type of the instance. Valid values:`Classic` or `VPC`.
* `name` - (Optional) The name of DB instance. It must be 2 to 256 characters in length.
* `instance_charge_type` - (Optional) The billing method of the instance. Default value: `PostPaid`. Valid values: `PrePaid`, `PostPaid`. **NOTE:** It can be modified from `PostPaid` to `PrePaid` after version v1.141.0.
* `period` - (Optional, Int) The duration that you will buy DB instance (in month). It is valid when `instance_charge_type` is `PrePaid`. Default value: `1`. Valid values: [1~9], 12, 24, 36.
* `security_ip_list` - (Optional, List) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]). System default to `["127.0.0.1"]`.
* `account_password` - (Optional, Sensitive) Password of the root account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `kms_encrypted_password` - (Optional, Available since v1.57.1) An KMS encrypts password used to a instance. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available since v1.57.1) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `resource_group_id` - (Optional, Available since v1.161.0) The ID of the Resource Group.
* `auto_renew` - (Optional, Bool, Available since v1.141.0) Auto renew for prepaid. Default value: `false`. Valid values: `true`, `false`.
* `backup_time` - (Optional, Available since v1.42.0) Sharding Instance backup time. It is required when `backup_period` was existed. In the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. If not set, the system will return a default, like "23:00Z-24:00Z".
* `backup_period` - (Optional, List, Available since v1.42.0) MongoDB Instance backup period. It is required when `backup_time` was existed. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]
* `backup_retention_policy_on_cluster_deletion` - (Optional, Int, Available since v1.235.0) The backup retention policy configured for the instance. Valid values:
  - `0`: All backup sets are immediately deleted when the instance is released.
  - `1 `: Automatic backup is performed when the instance is released and the backup set is retained for a long period of time.
  - `2 `: Automatic backup is performed when the instance is released and all backup sets are retained for a long period of time.
* `tde_status` - (Optional, Available since v1.76.0) The TDE(Transparent Data Encryption) status. It can be updated from version 1.160.0.
* `mongo_list` - (Required, Set) The Mongo nodes of the instance. The mongo-node count can be purchased is in range of [2, 32]. See [`mongo_list`](#mongo_list) below.
* `shard_list` - (Required, Set) The Shard nodes of the instance. The shard-node count can be purchased is in range of [2, 32]. See [`shard_list`](#shard_list) below.
* `config_server_list` - (Optional, ForceNew, Set, Available since v1.223.0) The ConfigServer nodes of the instance. See [`config_server_list`](#config_server_list) below.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `order_type` - (Optional, Available since v1.134.0) The type of configuration changes performed. Default value: `DOWNGRADE`. Valid values:
  - `UPGRADE`: The specifications are upgraded.
  - `DOWNGRADE`: The specifications are downgraded.
**NOTE:** `order_type` is only applicable to instances when `instance_charge_type` is `PrePaid`.

### `mongo_list`

The mongo_list supports the following:

* `node_class` -(Required) The instance type of the mongo node. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).

### `shard_list`

The shard_list supports the following:

* `node_class` - (Required) The instance type of the shard node. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
* `node_storage` - (Required, Int) The storage space of the shard node.
  - Custom storage space; value range: [10, 1,000]
  - 10-GB increments. Unit: GB.
* `readonly_replicas` - (Optional, Int, Available since v1.126.0) The number of read-only nodes in shard node Default value: `0`. Valid values: `0` to `5`.

### `config_server_list`

The config_server_list supports the following:

* `node_class` - (Optional, ForceNew) The instance type of the ConfigServer node. Valid values: `mdb.shard.2x.xlarge.d`, `dds.cs.mid`.
* `node_storage` - (Optional, ForceNew, Int) The storage space of the ConfigServer node.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Sharding Instance.
* `retention_period` - (Available since v1.42.0) Instance data backup retention days.
* `mongo_list` - The mongo nodes of the instance.
  * `node_id` - The ID of the mongo node.
  * `connect_string` - The endpoint of the mongo node.
  * `port` - The port number that is used to connect to the mongo node.
* `shard_list` - The information of the shard node.
  * `node_id` - The ID of the shard node.
* `config_server_list` - The information of the ConfigServer nodes.
  * `node_id` - The ID of the Config Server node.
  * `connect_string` - The connection address of the Config Server node.
  * `port` - The connection port of the Config Server node.
  * `max_connections` - The max connections of the Config Server node.
  * `max_iops` - The maximum IOPS of the Config Server node.
  * `node_description` - The description of the Config Server node.

## Timeouts

-> **NOTE:** Available since v1.126.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the Sharding Instance (until it reaches the initial `Running` status).
* `update` - (Defaults to 30 mins) Used when updating the Sharding Instance (until it reaches the initial `Running` status).
* `delete` - (Defaults to 30 mins) Used when deleting the Sharding Instance.

## Import

MongoDB Sharding Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_sharding_instance.example dds-bp1291daeda44195
```
