---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_instance"
sidebar_current: "docs-alicloud-resource-kvstore-instance"
description: |-
  Provides an ApsaraDB Redis / Memcache instance resource.
---

# alicloud\_kvstore\_instance

Provides an ApsaraDB Redis / Memcache instance resource. A DB instance is an isolated database environment in the cloud. It support be associated with IP whitelists and backup configuration which are separate resource providers. For information about Alicloud KVStore DBInstance more and how to use it, see [What is Resource Alicloud KVStore DBInstance](https://www.alibabacloud.com/help/doc-detail/60873.htm).

## Example Usage

Basic Usage

```terraform
resource "alicloud_kvstore_instance" "example" {
  db_instance_name = "tf-test-basic"
  vswitch_id       = "vsw-123456"
  security_ips = [
  "10.23.12.24"]
  instance_type  = "Redis"
  engine_version = "4.0"
  config = {
    appendonly             = "yes",
    lazyfree-lazy-eviction = "yes",
  }
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = "rg-123456"
  zone_id           = "cn-beijing-h"
  instance_class    = "redis.master.large.default"
}
```

Transform To PrePaid
```terraform
resource "alicloud_kvstore_instance" "example" {
  db_instance_name = "tf-test-basic"
  vswitch_id       = "vsw-123456"
  security_ips = [
  "10.23.12.24"]
  instance_type  = "Redis"
  engine_version = "4.0"
  config = {
    appendonly             = "yes",
    lazyfree-lazy-eviction = "yes",
  }
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = "rg-123456"
  zone_id           = "cn-beijing-h"
  instance_class    = "redis.master.large.default"
  payment_type      = "PrePaid"
  period            = "12"
}
```

Modify Private Connection String
```terraform
resource "alicloud_kvstore_instance" "example" {
  db_instance_name = "tf-test-basic"
  vswitch_id       = "vsw-123456"
  security_ips = [
  "10.23.12.24"]
  instance_type  = "Redis"
  engine_version = "4.0"
  config = {
    appendonly             = "yes",
    lazyfree-lazy-eviction = "yes",
  }
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id         = "rg-123456"
  zone_id                   = "cn-beijing-h"
  instance_class            = "redis.master.large.default"
  private_connection_prefix = "privateconnectionstringprefix"
}
```

## Argument Reference

The following arguments are supported:
* `instance_name` - (Optional) It has been deprecated from provider version 1.101.0 and `db_instance_name` instead.
* `db_instance_name` - (Optional, Available in 1.101.0+) The name of KVStore DBInstance. It is a string of 2 to 256 characters. 
* `password`- (Optional, Sensitive) The password of the KVStore DBInstance. The password that is used to connect to the instance. The password must be 8 to 32 characters in length and must contain at least three of the following character types: uppercase letters, lowercase letters, special characters, and digits. Special characters include: `! @ # $ % ^ & * ( ) _ + - =`
* `kms_encrypted_password` - (Optional, Available in 1.57.1+) An KMS encrypts password used to an instance. If the `password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `instance_class` - (Optional) Type of the applied ApsaraDB for Redis instance. It can be retrieved by data source [`alicloud_kvstore_instance_classes`](https://www.terraform.io/docs/providers/alicloud/d/kvstore_instance_classes)
or referring to help-docs [Instance type table](https://www.alibabacloud.com/help/doc-detail/26350.htm).
* `capacity` - (Optional, ForceNew, Available in 1.101.0+) The storage capacity of the KVStore DBInstance. Unit: MB.
* `availability_zone` - (Optional) It has been deprecated from provider version 1.101.0 and `zone_id` instead.
* `zone_id` - (Optional, Available in 1.101.0+) The ID of the zone. 
* `secondary_zone_id` - (Optional, Available in 1.128.0+) The ID of the secondary zone to which you want to migrate the ApsaraDB for Redis instance.

-> **NOTE:** If you specify this parameter, the master node and replica node of the instance can be deployed in different zones and disaster recovery is implemented across zones. The instance can withstand failures in data centers.

* `instance_charge_type` - (Optional) It has been deprecated from provider version 1.101.0 and `payment_type` instead.
* `payment_type` - (Optional, Available in 1.101.0+) The billing method of the KVStore DBInstance. Valid values: `PrePaid`, `PostPaid`. Default to `PostPaid`.
* `period` - (Optional) The duration that you will buy KVStore DBInstance (in month). It is valid when payment_type is `PrePaid`. Valid values: `[1~9]`, `12`, `24`, `36`.
* `auto_renew` - (Optional, Available in 1.36.0+) Whether to renewal a KVStore DBInstance automatically or not. It is valid when payment_type is `PrePaid`. Default to `false`.
* `auto_renew_period` - (Optional, Available in 1.36.0+) Auto-renewal period of an KVStore DBInstance, in the unit of the month. It is valid when payment_type is `PrePaid`. Valid value: [1~12], Default to `1`.
* `instance_type` - (Optional, ForceNew) The engine type of the KVStore DBInstance. Valid values: `Redis` or `Memcache`. Defaults to `Redis`.
* `vswitch_id` - (Optional) The ID of VSwitch.
* `engine_version`- (Optional) The engine version of the KVStore DBInstance. Valid values: ["2.8", "4.0", "5.0", "6.0"]. Default to "5.0".
  **NOTE:** When `instance_type = Memcache`, the `engine_version` only supports "4.0". 
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.
* `security_ips`- (Optional) The IP addresses in the whitelist group. The maximum number of IP addresses in the whitelist group is 1000. 
* `security_ip_group_attribute`- (Optional, Available in 1.101.0+) The value of this parameter is empty by default. The attribute of the whitelist group. The console does not display the whitelist group whose value of this parameter is hidden.
* `security_ip_group_name`- (Optional, Available in 1.101.0+) The name of the whitelist group.
* `modify_mode`- (Optional, Available in 1.101.0+) The method of modifying the whitelist. Valid values: `0`, `1` and `2`. Default to `0`. `0` means overwrites the original whitelist. `1` means adds the IP addresses to the whitelist. `2` means deletes the IP addresses from the whitelist.
* `security_group_id` - (Optional, Available in 1.76.0+) The ID of security groups.
* `private_ip`- (Optional) The internal IP address of the instance.
* `backup_id`- (Optional, ForceNew) The ID of the backup file of the source instance.
* `srcdb_instance_id`- (Optional, ForceNew, Available in 1.101.0+) The ID of the source instance.
* `restore_time`- (Optional, ForceNew, Available in 1.101.0+) The point in time of a backup file.
* `vpc_auth_mode`- (Optional) Only meaningful if instance_type is `Redis` and network type is VPC. Valid values: `Close`, `Open`. Defaults to `Open`.  `Close` means the redis instance can be accessed without authentication. `Open` means authentication is required.
* `parameters` - (Optional) It has been deprecated from provider version 1.101.0 and `config` instead..
* `config` - (Optional, Available in 1.101.0+) The configuration of the KVStore DBInstance. Available parameters can refer to the latest docs [Instance configurations table](https://www.alibabacloud.com/help/doc-detail/61209.htm) .
* `maintain_start_time` - (Optional, Available in v1.56.0+) The start time of the operation and maintenance time period of the KVStore DBInstance, in the format of HH:mmZ (UTC time).
* `maintain_end_time` - (Optional, Available in v1.56.0+) The end time of the operation and maintenance time period of the KVStore DBInstance, in the format of HH:mmZ (UTC time).
* `resource_group_id` - (Optional, Available in v1.86.0+) The ID of resource group which the resource belongs.
* `enable_public` - (Optional, Available in v1.94.0+) It has been deprecated from provider version 1.101.0 and resource `alicloud_kvstore_connection` instead.
* `connection_string_prefix` - (Optional, Available in v1.94.0+) It has been deprecated from provider version 1.101.0 and resource `alicloud_kvstore_connection` instead.
* `port` - (Optional, Available in v1.94.0+) It has been deprecated from provider version 1.101.0 and resource `alicloud_kvstore_connection` instead.
* `order_type`- (Optional, Available in 1.101.0+) Specifies a change type when you change the configuration of a subscription instance. Valid values: `UPGRADE`, `DOWNGRADE`. Default to `UPGRADE`. `UPGRADE` means upgrades the configuration of a subscription instance. `DOWNGRADE` means downgrades the configuration of a subscription instance.
* `node_type`- (Optional, Available in 1.101.0+) "Field 'node_type' has been deprecated from version 1.120.1". This parameter is determined by the `instance_class`.
* `ssl_enable`- (Optional, Available in 1.101.0+) Modifies the SSL status. Valid values: `Disable`, `Enable` and `Update`. 
  Note: This functionality is supported by Cluster mode (Redis 2.8, 4.0, 5.0) and Standard mode( Redis 2.8 only)
* `force_upgrade`- (Optional, Available in 1.101.0+) Specifies whether to forcibly change the type. Default to: `true`.
* `dedicated_host_group_id`- (Optional, Available in 1.101.0+) The ID of the dedicated cluster. This parameter is required when you create an ApsaraDB for Redis instance in a dedicated cluster.
* `coupon_no`- (Optional, Available in 1.101.0+) The coupon code. Default to: `youhuiquan_promotion_option_id_for_blank`.
* `business_info`- (Optional, Available in 1.101.0+) The ID of the event or the business information.
* `auto_use_coupon`- (Optional, ForceNew, Available in 1.101.0+) Specifies whether to use a coupon. Default to: `false`.
* `instance_release_protection`- (Optional, Computed, Available in 1.101.0+) Whether to open the release protection.
* `global_instance_id`- (Optional, Available in 1.101.0+) The ID of distributed cache.
* `global_instance`- (Optional, ForceNew, Available in 1.101.0+) Whether to create a distributed cache. Default to: `false`.
* `backup_period`- (Optional, Available in 1.104.0+) Backup period. 
* `backup_time`- (Optional, Available in 1.104.0+) Backup time, the format is HH:mmZ-HH:mmZ (UTC time).
* `enable_backup_log`- (Optional, Available in 1.104.0+) Turn on or off incremental backup. Valid values: `1`, `0`. Default to `0`
* `private_connection_prefix`- (Optional, Available in 1.105.0+) Private network connection prefix, used to modify the private network connection address. Only supports updating private network connections for existing instance.
* `private_connection_port`- (Optional, Available in 1.124.0+) Private network connection port, used to modify the private network connection port.
* `dry_run` - (Optional, Available in 1.128.0+) Specifies whether to precheck the request. Valid values:
  * true: prechecks the request without creating an instance. The system prechecks the required parameters, request format, service limits, and available resources. If the request fails the precheck, the corresponding error message is returned. If the request passes the precheck, the DryRunOperation error code is returned.
  * false: checks the request. After the request passes the check, an instance is created.
  
-> **NOTE:** The start time to the end time must be 1 hour. For example, the MaintainStartTime is 01:00Z, then the MaintainEndTime must be 02:00Z.

-> **NOTE:** You must specify at least one of the `capacity` and `instance_class` parameters when you call create instance operation.

-> **NOTE:** The `private_ip` must be in the Classless Inter-Domain Routing (CIDR) block of the VSwitch to which the instance belongs.

-> **NOTE:** If you specify the `srcdb_instance_id` parameter, you must specify the `backup_id` or `restore_time` parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of KVStore DBInstance.
* `bandwidth` - The bandwidth.
* `end_time` - The expiration time of the prepaid instance.
* `qps` - Theoretical maximum QPS value.
* `status` - The status of KVStore DBInstance.
* `connection_domain`- Intranet connection address of the KVStore instance.

### Timeouts

-> **NOTE:** Available in 1.54.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the KVStore instance (until it reaches the initial `Normal` status). 
* `update` - (Defaults to 40 mins) Used when updating the KVStore instance (until it reaches the initial `Normal` status). 
* `delete` - It has been deprecated from provider version 1.101.0.

## Import

KVStore instance can be imported using the id, e.g.

```
$ terraform import alicloud_kvstore_instance.example r-abc12345678
```
