---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_instances"
sidebar_current: "docs-alicloud-datasource-kvstore-instances"
description: |-
    Provides a collection of kvstore instances according to the specified filters.
---

# alicloud\_kvstore\_instances

The `alicloud_kvstore_instances` data source provides a collection of kvstore instances available in Alicloud account.
Filters support regular expression for the instance name, searches by tags, and other filters which are listed below.

## Example Usage

```terraform
data "alicloud_kvstore_instances" "default" {
  name_regex = "testname"
}
output "first_instance_name" {
  value = data.alicloud_kvstore_instances.default.instances.0.name
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the instance name.
* `ids` - (Optional, Available 1.52.2+) A list of KVStore DBInstance IDs.
* `instance_type` - (Optional) The engine type of the KVStore DBInstance. Options are `Memcache`, and `Redis`. If no value is specified, all types are returned.
* `status` - (Optional) The status of the KVStore DBInstance. Valid values: `Changing`, `CleaningUpExpiredData`, `Creating`, `Flushing`, `HASwitching`, `Inactive`, `MajorVersionUpgrading`, `Migrating`, `NetworkModifying`, `Normal`, `Rebooting`, `SSLModifying`, `Transforming`, `ZoneMigrating`.
* `instance_class`- (Optional) Type of the applied ApsaraDB for Redis instance. For more information, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/61135.htm).
* `vpc_id` - (Optional) Used to retrieve instances belong to specified VPC.
* `vswitch_id` - (Optional) Used to retrieve instances belong to specified `vswitch` resources.
* `tags` - (Optional) Query the instance bound to the tag. The format of the incoming value is `json` string, including `TagKey` and `TagValue`. `TagKey` cannot be null, and `TagValue` can be empty. Format example `{"key1":"value1"}`.
* `architecture_type` - (Optional, Available in 1.101.0+) The type of the architecture. Valid values: `cluster`, `standard` and `SplitRW`.
* `edition_type` - (Optional, Available in 1.101.0+) Used to retrieve instances belong to specified `vswitch` resources.  Valid values: `Enterprise`, `Community`.
* `engine_version` - (Optional, Available in 1.101.0+) The engine version. Valid values: `2.8`, `4.0`, `5.0`, `6.0`.
* `expired` - (Optional, Available in 1.101.0+) The expiration status of the instance.
* `global_instance` - (Optional, Available in 1.101.0+) Whether to create a distributed cache.
* `network_type` - (Optional, Available in 1.101.0+) The type of the network. Valid values: `CLASSIC`, `VPC`.
* `payment_type` - (Optional, Available in 1.101.0+) The payment type. Valid values: `PostPaid`, `PrePaid`.
* `resource_group_id` - (Optional, Available in 1.101.0+) The ID of the resource group.
* `search_key` - (Optional, Available in 1.101.0+) The name of the instance.
* `zone_id` - (Optional, Available in 1.101.0+) The ID of the zone.
* `enable_details` - (Optional, Available in 1.101.0+) Default to `false`. Set it to true can output more details.
* `output_file` - (Optional) The name of file that can save the collection of instances after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of KVStore Instance IDs.
* `names` - A list of KVStore Instance names.
* `instances` - A list of KVStore Instances. Its every element contains the following attributes:
  * `id` - The ID of the instance.
  * `name` - It has been deprecated from provider version 1.101.0 and `db_instance_name` instead.
  * `db_instance_name` - The name of the instance.
  * `charge_type` - It has been deprecated from provider version 1.101.0 and `payment_type` instead.
  * `payment_type` - Billing method. Valid Values: `PostPaid` for  Pay-As-You-Go and `PrePaid` for subscription.
  * `expire_time` - It has been deprecated from provider version 1.101.0 and `end_time` instead.
  * `end_time` - Expiration time. Pay-As-You-Go instances are never expire.
  * `availability_zone` - It has been deprecated from provider version 1.101.0 and `zone_id` instead.
  * `zone_id` - The ID of zone.
  * `connections` - IIt has been deprecated from provider version 1.101.0 and `max_connections` instead.
  * `max_connections` - Instance connection quantity limit. Unit: count.
  * `status` - Status of the instance.
  * `instance_type` - (Optional) Database type. Valid Values: `Memcache`, `Redis`. If no value is specified, all types are returned.
  * `instance_class`- (Optional) Type of the applied ApsaraDB for instance.
  For more information, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/61135.htm).
  * `vpc_id` - VPC ID the instance belongs to.
  * `vswitch_id` - VSwitch ID the instance belongs to.
  * `private_ip` - Private IP address of the instance.
  * `capacity` - Capacity of the applied ApsaraDB for the instance. Unit: MB.
  * `bandwidth` - Instance bandwidth limit. Unit: Mbit/s.
  * `config` - The parameter configuration of the instance.
  * `connection_mode` - The connection mode of the instance.
  * `db_instance_id` - The ID of the instance.
  * `destroy_time` - The time when the instance was destroyed.
  * `engine_version` - The engine version of the instance.
  * `has_renew_change_order` - Indicates whether there was an order of renewal with configuration change that had not taken effect.
  * `is_rds` - Indicates whether the instance is managed by Relational Database Service (RDS).
  * `network_type` - The network type of the instance.
  * `node_type` - The node type of the instance.
  * `package_type` - The type of the package.
  * `port` - The service port of the instance.
  * `qps` - The queries per second (QPS) supported by the instance.
  * `replacate_id` - The logical ID of the replica instance.
  * `vpc_cloud_instance_id` - Connection port of the instance.
  * `region_id` - Region ID the instance belongs to.
  * `create_time` - Creation time of the instance.
  * `user_name` - The username of the instance.
  * `connection_domain` - Instance connection domain (only Intranet access supported).
  * `secondary_zone_id` - (Optional, Available in 1.128.0+) The ID of the secondary zone to which you want to migrate the ApsaraDB for Redis instance.

    
    
