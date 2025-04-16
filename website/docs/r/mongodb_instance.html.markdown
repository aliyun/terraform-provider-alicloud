---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_instance"
sidebar_current: "docs-alicloud-resource-mongodb-instance"
description: |-
  Provides a MongoDB instance resource.
---

# alicloud_mongodb_instance

Provides a MongoDB instance resource supports replica set instances only. the MongoDB provides stable, reliable, and automatic scalable database services.
It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.
You can see detail product introduction [here](https://www.alibabacloud.com/help/doc-detail/26558.htm)

-> **NOTE:** Available since v1.37.0.

-> **NOTE:**  The following regions don't support create Classic network MongoDB instance.
[`cn-zhangjiakou`,`cn-huhehaote`,`ap-southeast-3`,`ap-southeast-5`,`me-east-1`,`ap-northeast-1`,`eu-west-1`]

-> **NOTE:**  Create MongoDB instance or change instance type and storage would cost 5~10 minutes. Please make full preparation

## Example Usage

### Create a Mongodb instance

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mongodb_instance&exampleId=2769a854-b7b1-2ac5-c852-950339b763843c914a7d&activeTab=example&spm=docs.r.mongodb_instance.0.2769a854b7&intl_lang=EN_US" target="_blank">
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

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.2"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  vswitch_id          = alicloud_vswitch.default.id
  security_ip_list    = ["10.168.1.12", "100.69.7.112"]
  name                = var.name
  tags = {
    Created = "TF"
    For     = "example"
  }
}
```

## Module Support

You can use to the existing [mongodb module](https://registry.terraform.io/modules/terraform-alicloud-modules/mongodb/alicloud)
to create a MongoDB instance resource one-click.

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/61763.htm) `EngineVersion`. **NOTE:** From version 1.225.0, `engine_version` can be modified.
* `db_instance_class` - (Required) Instance specification. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
* `db_instance_storage` - (Required, Int) User-defined DB instance storage space.Unit: GB. Value range:
  - Custom storage space.
  - 10-GB increments.
* `storage_engine` (Optional, ForceNew) The storage engine of the instance. Default value: `WiredTiger`. Valid values: `WiredTiger`, `RocksDB`.
* `storage_type` - (Optional, Available since v1.199.0) The storage type of the instance. Valid values: `cloud_essd1`, `cloud_essd2`, `cloud_essd3`, `cloud_auto`, `local_ssd`. **NOTE:** From version 1.229.0, `storage_type` can be modified. However, `storage_type` can only be modified to `cloud_auto`.
* `provisioned_iops` - (Optional, Int, Available since v1.229.0) The provisioned IOPS. Valid values: `0` to `50000`.
* `vpc_id` - (Optional, ForceNew, Available since v1.161.0) The ID of the VPC. -> **NOTE:** `vpc_id` is valid only when `network_type` is set to `VPC`.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance. it supports multiple zone.
  If it is a multi-zone and `vswitch_id` is specified, the vswitch must in one of them.
  The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alicloud_zones`.
* `secondary_zone_id` - (Optional, ForceNew, Available since v1.199.0) Configure the available area where the slave node (Secondary node) is located to realize multi-available area deployment. **NOTE:** This parameter value cannot be the same as `zone_id` and `hidden_zone_id` parameter values.
* `hidden_zone_id` - (Optional, ForceNew, Available since v1.199.0) Configure the zone where the hidden node is located to deploy multiple zones. **NOTE:** This parameter value cannot be the same as `zone_id` and `secondary_zone_id` parameter values.
* `security_group_id` - (Optional, Available since v1.73.0) The Security Group ID of ECS.
* `replication_factor` - (Optional, Int) Number of replica set nodes. Valid values: `1`, `3`, `5`, `7`.
* `network_type` - (Optional, ForceNew, Available since v1.161.0) The network type of the instance. Valid values:`Classic`, `VPC`.
* `name` - (Optional) The name of DB instance. It must be 2 to 256 characters in length.
* `instance_charge_type` - (Optional) The billing method of the instance. Default value: `PostPaid`. Valid values: `PrePaid`, `PostPaid`. **NOTE:** It can be modified from `PostPaid` to `PrePaid` after version 1.63.0.
* `period` - (Optional, Int) The duration that you will buy DB instance (in month). It is valid when `instance_charge_type` is `PrePaid`. Default value: `1`. Valid values: [1~9], 12, 24, 36.
* `security_ip_list` - (Optional, List) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `account_password` - (Optional, Sensitive) Password of the root account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `kms_encrypted_password` - (Optional, Available since v1.57.1) An KMS encrypts password used to a instance. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available since v1.57.1) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `encrypted` - (Optional, ForceNew, Bool, Available since v1.212.0) Whether to enable cloud disk encryption. Default value: `false`. Valid values: `true`, `false`.
* `cloud_disk_encryption_key` - (Optional, ForceNew, Available since v1.212.0) The ID of the encryption key.
* `readonly_replicas` - (Optional, Int, Available since v1.199.0) The number of read-only nodes in the replica set instance. Default value: 0. Valid values: 0 to 5.
* `resource_group_id` - (Optional, Available since v1.161.0) The ID of the Resource Group.
* `auto_renew` - (Optional, Bool, Available since v1.141.0) Auto renew for prepaid. Default value: `false`. Valid values: `true`, `false`.
  -> **NOTE:** The start time to the end time must be 1 hour. For example, the MaintainStartTime is 01:00Z, then the MaintainEndTime must be 02:00Z.
* `backup_time` - (Optional, Available since v1.42.0) MongoDB instance backup time. It is required when `backup_period` was existed. In the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. If not set, the system will return a default, like "23:00Z-24:00Z".
* `backup_period` - (Optional, List, Available since v1.42.0) MongoDB Instance backup period. It is required when `backup_time` was existed. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday].
* `backup_retention_period` - (Optional, Int, Available since v1.213.1) The retention period of full backups.
* `backup_retention_policy_on_cluster_deletion` - (Optional, Int, Available since v1.235.0) The backup retention policy configured for the instance. Valid values:
  - `0`: All backup sets are immediately deleted when the instance is released.
  - `1 `: Automatic backup is performed when the instance is released and the backup set is retained for a long period of time.
  - `2 `: Automatic backup is performed when the instance is released and all backup sets are retained for a long period of time.
* `enable_backup_log` - (Optional, Int, Available since v1.230.1) Specifies whether to enable the log backup feature. Valid values:
  - `0`: The log backup feature is disabled.
  - `1 `: The log backup feature is enabled.
* `log_backup_retention_period` - (Optional, Int, Available since v1.230.1) The number of days for which log backups are retained. Valid values: `7` to `730`. **NOTE:** `log_backup_retention_period` is valid only when `enable_backup_log` is set to `1`.
* `snapshot_backup_type` - (Optional, Available since v1.212.0) The snapshot backup type. Default value: `Standard`. Valid values:
  - `Standard`: standard backup.
  - `Flash `: single-digit second backup.
* `backup_interval` - (Optional, Available since v1.212.0) The frequency at which high-frequency backups are created. Valid values: `-1`, `15`, `30`, `60`, `120`, `180`, `240`, `360`, `480`, `720`.
* `ssl_action` - (Optional, Available since v1.78.0) Actions performed on SSL functions. Valid values:
  - `Open`: turn on SSL encryption.
  - `Close`: turn off SSL encryption.
  - `Update`: update SSL certificate.
* `maintain_start_time` - (Optional, Available since v1.56.0) The start time of the operation and maintenance time period of the instance, in the format of HH:mmZ (UTC time).
* `maintain_end_time` - (Optional, Available since v1.56.0) The end time of the operation and maintenance time period of the instance, in the format of HH:mmZ (UTC time).
* `effective_time` - (Optional, Available since v1.215.0) The time when the changed configurations take effect. Valid values: `Immediately`, `MaintainTime`.
* `order_type` - (Optional, Available since v1.134.0) The type of configuration changes performed. Default value: `DOWNGRADE`. Valid values:
  - `UPGRADE`: The specifications are upgraded.
  - `DOWNGRADE`: The specifications are downgraded.
    **NOTE:** `order_type` is only applicable to instances when `instance_charge_type` is `PrePaid`.
* `tde_status` - (Optional, Available since v1.73.0) The TDE(Transparent Data Encryption) status. Note: `tde_status` cannot be set to `disabled` after it is enabled, see [Transparent Data Encryption](https://www.alibabacloud.com/help/en/mongodb/user-guide/configure-tde-for-an-apsaradb-for-mongodb-instance) for more details.
* `encryptor_name` - (Optional, Available since v1.212.0) The encryption method. **NOTE:** `encryptor_name` is valid only when `tde_status` is set to `enabled`.
* `encryption_key` - (Optional, Available since v1.212.0) The ID of the custom key.
* `role_arn` - (Optional, Available since v1.212.0) The Alibaba Cloud Resource Name (ARN) of the specified Resource Access Management (RAM) role.
* `parameters` - (Optional, Set, Available since v1.203.0) Set of parameters needs to be set after mongodb instance was launched. See [`parameters`](#parameters) below.
* `tags` - (Optional, Available since v1.66.0) A mapping of tags to assign to the resource.

### `parameters`

The parameters supports the following:

* `name` - (Required) The name of the parameter.
* `value` - (Required) The value of the parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MongoDB.
* `retention_period` - Instance data backup retention days. Available since v1.42.0.
* `replica_set_name` - The name of the mongo replica set.
* `ssl_status` - Status of the SSL feature.
* `replica_sets` - Replica set instance information.
  * `vpc_id` - The private network ID of the node.
  * `vswitch_id` - The virtual switch ID to launch DB instances in one VPC.
  * `network_type` - The network type of the node.
  * `vpc_cloud_instance_id` - VPC instance ID.
  * `replica_set_role` - The role of the node.
  * `connection_domain` - The connection address of the node.
  * `connection_port` - The connection port of the node.
  * `role_id` - The id of the role.

## Timeouts

-> **NOTE:** Available since v1.53.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the MongoDB instance (until it reaches the initial `Running` status).
* `update` - (Defaults to 30 mins) Used when updating the MongoDB instance (until it reaches the initial `Running` status).
* `delete` - (Defaults to 30 mins) Used when deleting the MongoDB instance.

## Import

MongoDB instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_instance.example dds-bp1291daeda44194
```
