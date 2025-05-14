---
subcategory: "Tair (Redis OSS-Compatible) And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_instance"
sidebar_current: "docs-alicloud-resource-kvstore-instance"
description: |-
  Provides a Tair (Redis OSS-Compatible) And Memcache (KVStore) Classic Instance resource.
---

# alicloud_kvstore_instance

Provides  Tair (Redis OSS-Compatible) And Memcache (KVStore) Classic Instance resource. A DB instance is an isolated database environment in the cloud. It support be associated with IP whitelists and backup configuration which are separate resource providers. For information about Alicloud KVStore DBInstance more and how to use it, see [What is Resource Alicloud KVStore DBInstance](https://www.alibabacloud.com/help/en/redis/developer-reference/api-r-kvstore-2015-01-01-createinstances-redis).

-> **NOTE:** Available since v1.14.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kvstore_instance&exampleId=34a3c4a6-2047-e76c-8908-f436040a49e1fd671eac&activeTab=example&spm=docs.r.kvstore_instance.0.34a3c4a620&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_kvstore_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_kvstore_zones.default.zones.0.id
}

resource "alicloud_kvstore_instance" "default" {
  db_instance_name  = var.name
  vswitch_id        = alicloud_vswitch.default.id
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  zone_id           = data.alicloud_kvstore_zones.default.zones.0.id
  instance_class    = "redis.master.large.default"
  instance_type     = "Redis"
  engine_version    = "5.0"
  security_ips      = ["10.23.12.24"]
  config = {
    appendonly             = "yes"
    lazyfree-lazy-eviction = "yes"
  }
  tags = {
    Created = "TF",
    For     = "example",
  }
}
```

Launching a PrePaid instance

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kvstore_instance&exampleId=e922fa80-fae0-feb0-67b9-cebfa4431999756f5560&activeTab=example&spm=docs.r.kvstore_instance.1.e922fa80fa&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example-prepaid"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_kvstore_zones" "default" {
  instance_charge_type = "PrePaid"
}

// PrePaid instance can not deleted and there suggests using an existing vpc and vswitch, like default vpc.
data "alicloud_vpcs" "default" {
  is_default = true
}

data "alicloud_vswitches" "default" {
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
  vpc_id  = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_kvstore_instance" "default" {
  db_instance_name  = var.name
  vswitch_id        = data.alicloud_vswitches.default.ids.0
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  zone_id           = data.alicloud_kvstore_zones.default.zones.0.id
  secondary_zone_id = data.alicloud_kvstore_zones.default.zones.1.id
  instance_class    = "redis.master.large.default"
  instance_type     = "Redis"
  engine_version    = "5.0"
  payment_type      = "PrePaid"
  period            = "12"
  security_ips      = ["10.23.12.24"]
  config = {
    appendonly             = "no"
    lazyfree-lazy-eviction = "no"
    EvictionPolicy         = "volatile-lru"
  }
  tags = {
    Created = "TF",
    For     = "example",
  }
}
```

Setting Private Connection String

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kvstore_instance&exampleId=e71b3868-131e-85db-2ff0-dea29fe822e37d35849d&activeTab=example&spm=docs.r.kvstore_instance.2.e71b386813&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example-with-connection"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_kvstore_zones" "default" {
  product_type = "OnECS"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_kvstore_zones.default.zones.0.id
}

resource "alicloud_kvstore_instance" "default" {
  db_instance_name          = var.name
  vswitch_id                = alicloud_vswitch.default.id
  resource_group_id         = data.alicloud_resource_manager_resource_groups.default.ids.0
  zone_id                   = data.alicloud_kvstore_zones.default.zones.0.id
  secondary_zone_id         = data.alicloud_kvstore_zones.default.zones.1.id
  instance_class            = "redis.shard.small.ce"
  instance_type             = "Redis"
  engine_version            = "7.0"
  maintain_start_time       = "04:00Z"
  maintain_end_time         = "06:00Z"
  backup_period             = ["Wednesday"]
  backup_time               = "11:00Z-12:00Z"
  private_connection_prefix = "exampleconnectionprefix"
  private_connection_port   = 4011
  security_ips              = ["10.23.12.24"]
  config = {
    appendonly             = "yes"
    lazyfree-lazy-eviction = "yes"
    EvictionPolicy         = "volatile-lru"
  }
  tags = {
    Created = "TF",
    For     = "example",
  }
}
```

### Deleting `alicloud_kvstore_instance` or removing it from your configuration

The `alicloud_kvstore_instance` resource allows you to manage `payment_type = "Prepaid"` db instance, but Terraform cannot destroy it.
From version 1.201.0, deleting the subscription resource or removing it from your configuration will remove it
from your state file and management, but will not destroy the DB Instance.
You can resume managing the subscription db instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `instance_name` - (Deprecated since v1.101.0) It has been deprecated from provider version 1.101.0 and `db_instance_name` instead.
* `db_instance_name` - (Optional, Available since v1.101.0) The name of KVStore DBInstance. It is a string of 2 to 256 characters.
* `password`- (Optional, Sensitive) The password of the KVStore DBInstance. The password that is used to connect to the instance. The password must be 8 to 32 characters in length and must contain at least three of the following character types: uppercase letters, lowercase letters, special characters, and digits. Special characters include: `! @ # $ % ^ & * ( ) _ + - =`
* `kms_encrypted_password` - (Optional, Available since v1.57.1) An KMS encrypts password used to an instance. If the `password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available since v1.57.1) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `instance_class` - (Optional) Type of the applied Tair (Redis OSS-Compatible) And Memcache (KVStore) Classic Instance. It can be retrieved by data source [alicloud_kvstore_instance_classes](https://www.terraform.io/docs/providers/alicloud/d/kvstore_instance_classes)
  or referring to help-docs [Instance type table](https://www.alibabacloud.com/help/doc-detail/26350.htm).
* `capacity` - (Optional, ForceNew, Int, Available since v1.101.0) The storage capacity of the KVStore DBInstance. Unit: MB.
* `availability_zone` - (Deprecated since v1.101.0) It has been deprecated from provider version 1.101.0 and `zone_id` instead.
* `zone_id` - (Optional, Available since v1.101.0) The ID of the zone.
* `secondary_zone_id` - (Optional, Available since v1.128.0) The ID of the secondary zone to which you want to migrate the Tair (Redis OSS-Compatible) And Memcache (KVStore) Classic Instance.
  -> **NOTE:** If you specify this parameter, the master node and replica node of the instance can be deployed in different zones and disaster recovery is implemented across zones. The instance can withstand failures in data centers.
* `instance_charge_type` - (Deprecated since v1.101.0) It has been deprecated from provider version 1.101.0 and `payment_type` instead.
* `payment_type` - (Optional, Available since v1.101.0) The billing method of the KVStore DBInstance. Valid values: `PrePaid`, `PostPaid`. Default value: `PostPaid`.
* `period` - (Optional) The duration that you will buy KVStore DBInstance (in month). It is valid when payment_type is `PrePaid`. Valid values: `[1~9]`, `12`, `24`, `36`.
* `auto_renew` - (Optional, Bool, Available since v1.36.0) Whether to renewal a KVStore DBInstance automatically or not. It is valid when payment_type is `PrePaid`. Default value: `false`.
* `auto_renew_period` - (Optional, Int, Available since v1.36.0) Auto-renewal period of an KVStore DBInstance, in the unit of the month. It is valid when payment_type is `PrePaid`. Valid values: [1~12]. Default value: `1`.
* `instance_type` - (Optional, ForceNew) The engine type of the KVStore DBInstance. Valid values: `Redis` or `Memcache`. Default value: `Redis`.
* `vswitch_id` - (Optional) The ID of VSwitch.
* `engine_version`- (Optional) The engine version of the KVStore DBInstance. Valid values: ["2.8", "4.0", "5.0", "6.0", "7.0"]. Default value: `5.0`.
  **NOTE:** When `instance_type = Memcache`, the `engine_version` only supports "4.0".
* `tags` - (Optional, Available since v1.55.3) A mapping of tags to assign to the resource.
* `security_ips`- (Optional, List) The IP addresses in the whitelist group. The maximum number of IP addresses in the whitelist group is 1000.
* `security_ip_group_attribute`- (Optional, Available since v1.101.0) The value of this parameter is empty by default. The attribute of the whitelist group. The console does not display the whitelist group whose value of this parameter is hidden.
* `security_ip_group_name`- (Optional, Available since v1.101.0) The name of the whitelist group.
* `security_group_id` - (Optional, Available since v1.76.0) The ID of security groups. Separate multiple security group IDs with commas (,), such as `sg-***,sg-***,sg-***`.
* `private_ip`- (Optional, ForceNew) The internal IP address of the instance.
* `backup_id`- (Optional, ForceNew) The ID of the backup file of the source instance.
* `srcdb_instance_id`- (Optional, ForceNew, Available since v1.101.0) The ID of the source instance.
* `restore_time`- (Optional, ForceNew, Available since v1.101.0) The point in time of a backup file.
* `vpc_auth_mode`- (Optional) Only meaningful if instance_type is `Redis` and network type is VPC. Valid values: `Close`, `Open`. Default value: `Open`. `Close` means the redis instance can be accessed without authentication. `Open` means authentication is required.
* `parameters` - (Deprecated since v1.101.0) It has been deprecated from provider version 1.101.0 and `config` instead. See [`parameters`](#parameters) below.
* `config` - (Optional, MapString, Available since v1.101.0) The configuration of the KVStore DBInstance. Available parameters can refer to the latest docs [Instance configurations table](https://www.alibabacloud.com/help/en/redis/user-guide/supported-parameters).
* `maintain_start_time` - (Optional, Available since v1.56.0) The start time of the operation and maintenance time period of the KVStore DBInstance, in the format of HH:mmZ (UTC time).
* `maintain_end_time` - (Optional, Available since v1.56.0) The end time of the operation and maintenance time period of the KVStore DBInstance, in the format of HH:mmZ (UTC time).
* `effective_time` - (Optional, Available since v1.204.0) The time when the database is switched after the instance is migrated, or when the major version is upgraded, or when the instance class is upgraded. Valid values:
  - `Immediately` (Default): The configurations are immediately changed.
  - `MaintainTime`: The configurations are changed within the maintenance window. You can set `maintain_start_time` and `maintain_end_time` to change the maintenance window.
* `resource_group_id` - (Optional, Available since v1.86.0) The ID of resource group which the resource belongs.
* `enable_public` - (Deprecated since v1.101.0) It has been deprecated from provider version 1.101.0 and resource `alicloud_kvstore_connection` instead.
* `connection_string_prefix` - (Deprecated since v1.101.0) It has been deprecated from provider version 1.101.0 and resource `alicloud_kvstore_connection` instead.
* `port` - (Optional, Int, Available since v1.94.0) It has been deprecated from provider version 1.101.0 and resource `alicloud_kvstore_connection` instead.
* `order_type`- (Optional, Available since v1.101.0) Specifies a change type when you change the configuration of a subscription instance. Valid values: `UPGRADE`, `DOWNGRADE`. Default value: `UPGRADE`. `UPGRADE` means upgrades the configuration of a subscription instance. `DOWNGRADE` means downgrades the configuration of a subscription instance.
* `node_type`- (Deprecated since v1.120.1) Node type, valid values:
  - `MASTER_SLAVE`: High availability (dual copies)
  - `STAND_ALONE`: Single copy
  - `double`: Dual copies
  - `single`: Single copy
-> **NOTE:** For cloud native instances, please select MASTER_SLAVE or STAND_ALONE. For classic instances, please select double or single.
* `ssl_enable`- (Optional, Available since v1.101.0) Modifies the SSL status. Valid values: `Disable`, `Enable` and `Update`.
  **NOTE:** This functionality is supported by Cluster mode (Redis 2.8, 4.0, 5.0) and Standard mode( Redis 2.8 only).
* `force_upgrade`- (Optional, Bool, Available since v1.101.0) Specifies whether to forcibly change the type. Default value: `true`.
* `dedicated_host_group_id`- (Optional, ForceNew, Available since v1.101.0) The ID of the dedicated cluster. This parameter is required when you create a Tair (Redis OSS-Compatible) And Memcache (KVStore) Classic Instance in a dedicated cluster.
* `coupon_no`- (Optional, Available since v1.101.0) The coupon code. **NOTE:** The default value `youhuiquan_promotion_option_id_for_blank` removed since v1.216.0, this can cause your status file to change even if it has not been modified, so please review your change plan before apply change plan.
* `business_info`- (Optional, Available since v1.101.0) The ID of the event or the business information.
* `auto_use_coupon`- (Optional, Bool, ForceNew, Available since v1.101.0) Specifies whether to use a coupon. Default value: `false`.
* `instance_release_protection`- (Optional, Bool, Available since v1.101.0) Whether to open the release protection.
* `global_instance_id`- (Optional, Available since v1.101.0) The ID of distributed cache.
* `global_instance`- (Optional, ForceNew, Bool, Available since v1.101.0) Whether to create a distributed cache. Default value: `false`.
* `backup_period`- (Optional, List, Available since v1.104.0) Backup period.
* `backup_time`- (Optional, Available since v1.104.0) Backup time, the format is HH:mmZ-HH:mmZ (UTC time).
* `enable_backup_log`- (Optional, Int, Available since v1.104.0) Turn on or off incremental backup. Valid values: `1`, `0`. Default value: `0`
* `private_connection_prefix`- (Optional, Available since v1.105.0) Private network connection prefix, used to modify the private network connection address. Only supports updating private network connections for existing instance.
* `private_connection_port`- (Optional, Available since v1.124.0) Private network connection port, used to modify the private network connection port.
* `dry_run` - (Optional, Bool, Available since v1.128.0) Specifies whether to precheck the request. Valid values:
  - `true`: prechecks the request without creating an instance. The system prechecks the required parameters, request format, service limits, and available resources. If the request fails the precheck, the corresponding error message is returned. If the request passes the precheck, the DryRunOperation error code is returned.
  - `false`: checks the request. After the request passes the check, an instance is created.
* `tde_status`- (Optional, Available since v1.200.0) Specifies whether to enable TDE. Valid values: `Enabled`.
**NOTE:**: TDE cannot be disabled after it is enabled. Before you enable it, evaluate whether this feature affects.
  your business. For more information, see [Enable TDE](https://www.alibabacloud.com/help/en/redis/user-guide/enable-tde).
* `encryption_name`- (Optional, Available since v1.200.0) The Encryption algorithm, default AES-CTR-256.Note that this parameter is only available when the TDEStatus parameter is Enabled.
* `encryption_key`- (Optional, Available since v1.200.0) The Custom key ID, which you can get by calling DescribeEncryptionKeyList.If this parameter is not passed, the key is automatically generated by the key management service. To create a custom key, you can call the CreateKey interface of the key management service.
* `role_arn`- (Optional, Available since v1.200.0) The Specify the global resource descriptor ARN (Alibaba Cloud Resource Name) information of the role to be authorized, and use the related key management services after the authorization is completed, in the format: `acs:ram::$accountID:role/$roleName`.
* `shard_count`- (Optional, Int, Available since v1.208.0) The number of data shards. This parameter is available only if you create a cluster instance that uses cloud disks. You can use this parameter to specify a custom number of data shards. **NOTE:** From version 1.216.0, `shard_count` can be modified.
* `read_only_count` - (Optional, Int, Available since v1.226.0) The number of read replicas in the primary zone. Valid values: `1` to `9`.
* `slave_read_only_count` - (Optional, Int, Available since v1.226.0) The number of read replicas in the secondary zone. **NOTE:**: When you create a multi-zone read/write splitting instance, you must specify both `secondary_zone_id` and `slave_read_only_count`.
-> **NOTE:** The sum of `read_only_count` and `slave_read_only_count` cannot be greater than `9`.
* `is_auto_upgrade_open` - (Optional, Available since v1.228.0) Specifies whether to enable automatic minor version update. Valid values:
  - `1`: Enables automatic minor version update.
  - `0`: Disables automatic minor version update.
* `bandwidth` - (Optional, Int) The total bandwidth of the instance. **NOTE:** From version 1.232.0, `bandwidth` can be set. If the instance is a cluster instance, `bandwidth` must be divisible by the number of `shard_count` in the instance, and if the instance is a read/write splitting instance, `bandwidth` cannot be set.
* `connection_string` - (Deprecated since v1.101.0) Indicates whether the address is a private endpoint.
* `modify_mode`- (Removed since v1.216.0) The method of modifying the whitelist. **NOTE:** Field `modify_mode` has been removed from provider version 1.216.0.

-> **NOTE:** The start time to the end time must be 1 hour. For example, the MaintainStartTime is 01:00Z, then the MaintainEndTime must be 02:00Z.

-> **NOTE:** You must specify at least one of the `capacity` and `instance_class` parameters when you call create instance operation.

-> **NOTE:** The `private_ip` must be in the Classless Inter-Domain Routing (CIDR) block of the VSwitch to which the instance belongs.

-> **NOTE:** If you specify the `srcdb_instance_id` parameter, you must specify the `backup_id` or `restore_time` parameter.

### `parameters`

The parameters supports the following:

* `name` (Deprecated since v1.101.0) Field `parameters` has been deprecated from provider version 1.101.0 and `config` instead.
* `value` (Deprecated since v1.101.0) Field `parameters` has been deprecated from provider version 1.101.0 and `config` instead.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of KVStore DBInstance.
* `end_time` - The expiration time of the prepaid instance.
* `qps` - Theoretical maximum QPS value.
* `connection_domain`- Intranet connection address of the KVStore instance.
* `status` - The status of KVStore DBInstance.

## Timeouts

-> **NOTE:** Available since v1.54.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the KVStore instance (until it reaches the initial `Normal` status).
* `update` - (Defaults to 40 mins) Used when updating the KVStore instance (until it reaches the initial `Normal` status).
* `delete` - It has been deprecated from provider version 1.101.0.

## Import

KVStore instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_kvstore_instance.example r-abc12345678
```
