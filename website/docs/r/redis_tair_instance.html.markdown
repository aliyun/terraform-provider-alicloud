---
subcategory: "Tair (Redis OSS-Compatible) And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_redis_tair_instance"
description: |-
  Provides a Alicloud Tair (Redis OSS-Compatible) And Memcache (KVStore) Tair Instance resource.
---

# alicloud_redis_tair_instance

Provides a Tair (Redis OSS-Compatible) And Memcache (KVStore) Tair Instance resource.

Describe the creation, deletion and query of tair instances.

For information about Tair (Redis OSS-Compatible) And Memcache (KVStore) Tair Instance and how to use it, see [What is Tair Instance](https://www.alibabacloud.com/help/en/tair).

-> **NOTE:** Available since v1.206.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_redis_tair_instance&exampleId=620b10c9-7933-f986-971f-30e49bed804e188829d7&activeTab=example&spm=docs.r.redis_tair_instance.0.620b10c979&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_kvstore_zones" "default" {
  product_type = "Tair_rdb"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids.0
  zone_id    = data.alicloud_kvstore_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_redis_tair_instance" "default" {
  payment_type       = "Subscription"
  period             = "1"
  instance_type      = "tair_rdb"
  zone_id            = local.zone_id
  instance_class     = "tair.rdb.2g"
  shard_count        = "2"
  vswitch_id         = local.vswitch_id
  vpc_id             = data.alicloud_vpcs.default.ids.0
  tair_instance_name = var.name
}
```

### Deleting `alicloud_redis_tair_instance` or removing it from your configuration

The `alicloud_redis_tair_instance` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `auto_renew` - (Optional) Specifies whether to enable auto-renewal for the instance. Default value: false. Valid values: true(enables auto-renewal), false(disables auto-renewal).
* `auto_renew_period` - (Optional) The subscription duration that is supported by auto-renewal. Unit: months. Valid values: 1, 2, 3, 6, and 12. This parameter is required only if the AutoRenew parameter is set to true.
* `backup_id` - (Optional, Available since v1.233.1) You can set the BackupId parameter to the backup set ID of the source instance. The system uses the data stored in the backup set to create an instance. You can call the DescribeBackups operation to query backup set IDs. If the source instance is a cluster instance, set the BackupId parameter to the backup set IDs of all shards of the source instance, separated by commas (,).

  If your instance is a cloud-native cluster instance, we recommend that you use DescribeClusterBackupList to query the backup set ID of the cluster instance. Then, set the ClusterBackupId request parameter to the backup set ID to clone the cluster instance. This eliminates the need to specify the backup set ID of each shard.
* `cluster_backup_id` - (Optional, Available since v1.224.0) This parameter is supported for specific new cluster instances. You can query the backup set ID by calling the DescribeClusterBackupList operation. If this parameter is supported, you can specify the backup set ID. In this case, you do not need to specify the BackupId parameter. If this parameter is not supported, set the BackupId parameter to the IDs of backup sets in all shards of the source instance, separated by commas (,).
* `connection_string_prefix` - (Optional, Available since v1.235.0) The prefix of the endpoint the instance, which must consist of lowercase letters and numbers and start with a lowercase letter.
* `effective_time` - (Optional) The time when to change the configurations. Default value: Immediately. Valid values: Immediately (The configurations are immediately changed), MaintainTime (The configurations are changed within the maintenance window).
* `engine_version` - (Optional, Computed) Database version. Default value: 1.0.

  Rules for transferring parameters of different tair product types:

  tair_rdb:  Compatible with the Redis5.0 and Redis6.0 protocols, and is transmitted to 5.0 or 6.0.

  tair_scm: The Tair persistent memory is compatible with the Redis6.0 protocol and is passed 1.0.

  tair_essd: The disk (ESSD/SSD) is compatible with the Redis4.0 and Redis6.0 protocols, and is transmitted to 1.0 and 2.0 respectively.
* `force_upgrade` - (Optional) Specifies whether to forcefully change the configurations of the instance. Default value: true. Valid values: false (The system does not forcefully change the configurations), true (The system forcefully changes the configurations).
* `global_instance_id` - (Optional, Available since v1.233.1) The ID of a distributed (Global Distributed Cache) instance, which indicates whether to use the newly created instance as a sub-instance of a distributed instance. You can use this method to create a distributed instance.

  1. Enter true if you want the new instance to be the first child instance.

  2. If you want the new instance to be used as the second and third sub-instances, enter the distributed instance ID.

  3. Not as a distributed instance, you do not need to enter any values.
* `instance_class` - (Required) The instance type of the instance. For more information, see [Instance types](https://www.alibabacloud.com/help/en/apsaradb-for-redis/latest/instance-types).
* `instance_type` - (Required, ForceNew) The storage medium of the instance. Valid values: tair_rdb, tair_scm, tair_essd.
* `intranet_bandwidth` - (Optional, Computed, Int, Available since v1.233.1) Instance intranet bandwidth
* `modify_mode` - (Optional, Available since v1.233.1) The modification method when modifying the IP whitelist. The value includes Cover (default): overwrite the original whitelist; Append: Append the whitelist; Delete: Delete the whitelist.
* `node_type` - (Optional, Computed) The node type. For cloud-native instances, input MASTER_SLAVE (master-replica) or STAND_ALONE (standalone). For classic instances, input double (master-replica) or single (standalone).
* `param_no_loose_sentinel_enabled` - (Optional, Computed, Available since v1.233.1) sentinel compatibility mode, applicable to non-cluster instances. For more information about parameters, see yes or no in the https://www.alibabacloud.com/help/en/redis/user-guide/use-the-sentinel-compatible-mode-to-connect-to-an-apsaradb-for-redis-instance, valid values: yes, no. The default value is no. 
* `param_no_loose_sentinel_password_free_access` - (Optional, Computed, Available since v1.237.0) Whether to allow Sentinel commands to be executed without secrets when Sentinel mode is enabled. Value: yes: enabled. After the command is enabled, you can directly run the Sentinel command in the VPC without enabling the password-free feature. no: the default value, disabled. For parameters, see https://help.aliyun.com/zh/redis/user-guide/use-the-sentinel-compatible-mode-to-connect-to-an-apsaradb-for-redis-instance
* `param_repl_mode` - (Optional, Computed, Available since v1.233.1) The value is semisync or async. The default value is async.

  The default data synchronization mode is asynchronous replication. To modify the data synchronization mode, refer to https://www.alibabacloud.com/help/en/redis/user-guide/modify-the-synchronization-mode-of-a-persistent-memory-optimized-instance 。
* `param_semisync_repl_timeout` - (Optional, Computed, Available since v1.233.1) The degradation threshold time of the semi-synchronous replication mode. This parameter value is required only when semi-synchronous replication is enabled. The unit is milliseconds, and the range is 10ms to 60000ms. The default value is 500ms. Please refer to: https://www.alibabacloud.com/help/en/redis/user-guide/modify-the-synchronization-mode-of-a-persistent-memory-optimized-instance。
* `param_sentinel_compat_enable` - (Optional, Computed, Available since v1.233.1) sentinel compatibility mode, applicable to instances in the cluster architecture proxy connection mode or read/write splitting architecture. For more information about the parameters, see https://www.alibabacloud.com/help/en/redis/user-guide/use-the-sentinel-compatible-mode-to-connect-to-an-apsaradb-for-redis-instance. The value is 0 or 1. The default value is 0.
* `password` - (Optional) The password that is used to connect to the instance. The password must be 8 to 32 characters in length and contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include ! @ # $ % ^ & * ( ) _ + - =
* `payment_type` - (Optional, Computed) Payment type: Subscription (prepaid), PayAsYouGo (postpaid). Default Subscription.
* `period` - (Optional, Int) The subscription duration. Unit: months. Valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24,36, and 60. This parameter is required only if you set the PaymentType parameter to Subscription.
* `port` - (Optional, ForceNew, Computed, Int) The Tair service port. The service port of the instance. Valid values: 1024 to 65535. Default value: 6379.
* `read_only_count` - (Optional, Int) Number of read-only nodes in the primary zone. Valid values: 0 to 5. This parameter is only applicable to the following conditions:

  If the instance is in the cloud disk version standard architecture, you can set this parameter to a value greater than 0 to enable the read/write splitting architecture.

  If the instance is a cloud disk version read/write splitting architecture instance, you can use this parameter to customize the number of read-only nodes, or set this parameter to 0 to disable the read/write splitting architecture and switch the instance to the standard architecture.
* `recover_config_mode` - (Optional, Available since v1.233.1) Whether to restore the account, kernel parameters, and whitelist (config) information from the original backup set when creating an instance using a specified backup set. The default value is empty, indicating that the account, kernel parameters, and whitelist information are not restored from the original backup set. This parameter is only applicable to Cloud Native instances, and the account, kernel parameters, and whitelist information must have been saved in the original backup set.
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the instance belongs.
* `secondary_zone_id` - (Optional, ForceNew) The ID of the secondary zone.This parameter is returned only if the instance is deployed in two zones.
* `security_group_id` - (Optional) Security group id
* `security_ip_group_name` - (Optional, Computed, Available since v1.233.1) The name of the IP address whitelist. You cannot modify the whitelist that is generated by the system. If you do not specify this parameter, the default whitelist is modified by default.
* `security_ips` - (Optional, Computed, Available since v1.233.1) The IP addresses in the whitelist. Up to 1,000 IP addresses can be specified in a whitelist. Separate multiple IP addresses with a comma (,). Specify an IP address in the 0.0.0.0/0, 10.23.12.24, or 10.23.12.24/24 format. In CIDR block 10.23.12.24/24, /24 specifies the length of the prefix of an IP address. The prefix length ranges from 1 to 32.
* `shard_count` - (Optional, Computed, Int) The number of data nodes in the instance. When 1 is passed, it means that the instance created is a standard architecture with only one data node. You can create an instance in the standard architecture that contains only a single data node. 2 to 32: You can create an instance in the cluster architecture that contains the specified number of data nodes. Only persistent memory-optimized instances can use the cluster architecture. Therefore, you can set this parameter to an integer from 2 to 32 only if you set the InstanceType parameter to tair_scm. It is not allowed to modify the number of shards by modifying this parameter after creating a master-slave architecture instance with or without passing 1.
* `slave_read_only_count` - (Optional, Int) Specifies the number of read-only nodes in the secondary zone when creating a multi-zone read/write splitting instance.

  Note: To create a multi-zone read/write splitting instance, slaveadonlycount and SecondaryZoneId must be specified at the same time.
* `src_db_instance_id` - (Optional, Available since v1.233.1) If you want to create an instance based on the backup set of an existing instance, set this parameter to the ID of the source instance. preceding three parameters. After you specify the SrcDBInstanceId parameter, use the BackupId, ClusterBackupId (recommended for cloud-native cluster instances), or RestoreTime parameter to specify the backup set or the specific point in time that you want to use to create an instance. The SrcDBInstanceId parameter must be used in combination with one of the preceding three parameters.
* `ssl_enabled` - (Optional, Computed) Modifies SSL encryption configurations. Valid values: 1. Disable (The SSL encryption is disabled) 2. Enable (The SSL encryption is enabled)  3. Update (The SSL certificate is updated)
* `storage_performance_level` - (Optional, ForceNew) The storage type. Valid values: PL1, PL2, and PL3. This parameter is available only when the value of InstanceType is tair_essd, that is, when an ESSD disk instance is selected.

  If the ESSD instance type is 4C, 8C, or 16C, you can specify the storage type as PL1.

  If the type of ESSD instance you select is 8C, 16C, 32C, or 52C, you can specify the storage type as PL2.

  If the ESSD instance type is 16C, 32C, or 52C, you can specify the storage type as PL3.
* `storage_size_gb` - (Optional, ForceNew, Computed, Int) Different specifications have different value ranges. When the instance_type value is tair_essd and the disk type is ESSD, this attribute takes effect and is required. When a Tair disk is an SSD, see-https://help.aliyun.com/zh/redis/product-overview/capacity-storage-type. The capacity field is defined as different fixed values according to different specifications, and does not need to be specified.
* `tags` - (Optional, Map) The tag of the resource
* `tair_instance_name` - (Optional) The name of the resource.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch to which the instance is connected.
* `vpc_auth_mode` - (Optional, Computed, Available since v1.233.1) The VPC authentication mode. Valid values: Open (enables password authentication), Close (disables password authentication and enables [password-free access](https://www.alibabacloud.com/help/en/apsaradb-for-redis/latest/enable-password-free-access)).
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC).
* `zone_id` - (Required, ForceNew) Zone ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `architecture_type` - The architecture of the instance.  cluster, standard, rwsplit.
* `connection_domain` - The internal endpoint of the instance.
* `create_time` - The time when the instance was created. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `max_connections` - The maximum number of connections supported by the instance.
* `network_type` - The network type of the instance.  CLASSIC(classic network), VPC.
* `region_id` - Region Id
* `status` - The status of the resource
* `tair_instance_id` - The ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 60 mins) Used when create the Tair Instance.
* `delete` - (Defaults to 30 mins) Used when delete the Tair Instance.
* `update` - (Defaults to 60 mins) Used when update the Tair Instance.

## Import

Tair (Redis OSS-Compatible) And Memcache (KVStore) Tair Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_redis_tair_instance.example <id>
```