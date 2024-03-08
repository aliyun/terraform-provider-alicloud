---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_redis_tair_instance"
description: |-
  Provides a Alicloud Redis Tair Instance resource.
---

# alicloud_redis_tair_instance

Provides a Redis Tair Instance resource. Describe the creation, deletion and query of tair instances.

For information about Redis Tair Instance and how to use it, see [What is Tair Instance](https://www.alibabacloud.com/help/en/tair).

-> **NOTE:** Available since v1.206.0.

## Example Usage

Basic Usage

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

The `alicloud_redis_tair_instance` resource allows you to manage  `payment_type = "PayAsYouGo"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `auto_renew` - (Optional) Specifies whether to enable auto-renewal for the instance. Default value: false. Valid values: true(enables auto-renewal), false(disables auto-renewal).
* `auto_renew_period` - (Optional) The subscription duration that is supported by auto-renewal. Unit: months. Valid values: 1, 2, 3, 6, and 12. This parameter is required only if the AutoRenew parameter is set to true.
* `effective_time` - (Optional) The time when to change the configurations. Default value: Immediately. Valid values: Immediately (The configurations are immediately changed), MaintainTime (The configurations are changed within the maintenance window).
* `engine_version` - (Optional, Computed) Database version. Default value: 1.0.  Rules for transferring parameters of different tair product types:  tair_rdb:  Compatible with the Redis5.0 and Redis6.0 protocols, and is transmitted to 5.0 or 6.0. tair_scm: The Tair persistent memory is compatible with the Redis6.0 protocol and is passed 1.0. tair_essd: The disk (ESSD/SSD) is compatible with the Redis4.0 and Redis6.0 protocols, and is transmitted to 1.0 and 2.0 respectively.
* `force_upgrade` - (Optional) Specifies whether to forcefully change the configurations of the instance. Default value: true. Valid values: false (The system does not forcefully change the configurations), true (The system forcefully changes the configurations).
* `instance_class` - (Required) The instance type of the instance. For more information, see [Instance types](https://www.alibabacloud.com/help/en/apsaradb-for-redis/latest/instance-types).
* `instance_type` - (Required, ForceNew) The storage medium of the instance. Valid values: tair_rdb, tair_scm, tair_essd.
* `password` - (Optional) The password that is used to connect to the instance. The password must be 8 to 32 characters in length and contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include ! @ # $ % ^ & * ( ) _ + - =.
* `payment_type` - (Optional, ForceNew, Computed) Payment type: Subscription (prepaid), PayAsYouGo (postpaid). Default PayAsYouGo.
* `period` - (Optional) The subscription duration. Unit: months. Valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24,36, and 60. This parameter is required only if you set the PaymentType parameter to Subscription.
* `port` - (Optional, ForceNew, Computed) The Tair service port. The service port of the instance. Valid values: 1024 to 65535. Default value: 6379.
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the instance belongs.
* `secondary_zone_id` - (Optional, ForceNew) The ID of the secondary zone.This parameter is returned only if the instance is deployed in two zones.
* `shard_count` - (Optional) The number of data nodes in the instance. When 1 is passed, it means that the instance created is a standard architecture with only one data node. You can create an instance in the standard architecture that contains only a single data node. 2 to 32: You can create an instance in the cluster architecture that contains the specified number of data nodes. Only persistent memory-optimized instances can use the cluster architecture. Therefore, you can set this parameter to an integer from 2 to 32 only if you set the InstanceType parameter to tair_scm. It is not allowed to modify the number of shards by modifying this parameter after creating a master-slave architecture instance with or without passing 1.
* `storage_performance_level` - (Optional, ForceNew, Available since v1.211.1) The storage type. The value range is [PL1, PL2, and PL3]. The default value is PL1. When the value of instance_type is "tair_essd", this attribute takes effect and is required.
* `storage_size_gb` - (Optional, ForceNew, Computed, Available since v1.211.1) The value range of different specifications is different, see [ESSD-based instances](https://www.alibabacloud.com/help/en/tair/product-overview/essd-based-instances). When the value of instance_type is "tair_essd", this attribute takes effect and is required.
* `tags` - (Optional, Map, Available since v1.211.1) The tag of the resource.
* `tair_instance_name` - (Optional) The name of the resource.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch to which the instance is connected.
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC).
* `zone_id` - (Required, ForceNew) Zone ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the instance was created. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 30 mins) Used when create the Tair Instance.
* `delete` - (Defaults to 15 mins) Used when delete the Tair Instance.
* `update` - (Defaults to 30 mins) Used when update the Tair Instance.

## Import

Redis Tair Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_redis_tair_instance.example <id>
```