---
subcategory: "SelectDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_selectdb_db_instance"
sidebar_current: "docs-alicloud-resource-selectdb-db-instance"
description: |-
  Provides a Alicloud SelectDB DBInstance resource.
---

# alicloud_selectdb_db_instance

Provides a SelectDB DBInstance resource.

For information about SelectDB DBInstance and how to use it, see [What is DBInstance](https://www.alibabacloud.com/help/zh/selectdb/latest/api-selectdb-2023-05-22-createdbinstance).

-> **NOTE:** Available since v1.229.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

variable "name" {
  default = "terraform_example"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_selectdb_db_instance" "default" {
  db_instance_class       = "selectdb.xlarge"
  db_instance_description = var.name
  desired_cache_size      = 200
  payment_type            = "PayAsYouGo"
  desired_engine_version  = "3.0"
  vpc_id                  = data.alicloud_vswitches.default.vswitches.0.vpc_id
  zone_id                 = data.alicloud_vswitches.default.vswitches.0.zone_id
  vswitch_id              = data.alicloud_vswitches.default.vswitches.0.id
}

```

## Argument Reference

The following arguments are supported:

* `db_instance_class` - (Required) The class for default cluster in DBInstance. db_cluster_class has a range of class from `selectdb.xlarge` to `selectdb.256xlarge`.
* `cache_size` - (Required) The cache size in DBInstance on creating default cluster. The number should be divided by 100.
* `payment_type` - (Required) The payment type of the resource. Valid values: `PayAsYouGo`,`Subscription`.
* `db_instance_description` - (Required) The DBInstance description.
* `period` - (Optional) It is valid when payment_type is `Subscription`. Valid values are `Year`, `Month`.
* `period_time` - (Optional) The duration that you will buy DBInstance. It is valid when payment_type is `Subscription`. Valid values: [1~9], 12, 24, 36.
* `desired_engine_version` - (Optional) The version of DBInstance when creating. SelectDB `3.0` and `4.0` are available.
* `zone_id` - (Required, ForceNew) The ID of zone for DBInstance.
* `vpc_id` - (Required, ForceNew) The ID of the VPC for DBInstance.
* `vswitch_id` - (Required, ForceNew) The ID of vswitch for DBInstance.
* `enable_public_network` - (Optional) If DBInstance need to open public network, set it to `true`.
* `upgraded_engine_minor_version` - (Optional) The DBInstance minor version want to upgraded to.
* `tags` - (Optional) A mapping of tags to assign to the resource.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `desired_security_ip_lists` - (Optional) The modified IP address whitelists. See [`desired_security_ip_lists`](#desired_security_ip_lists) below.

### `desired_security_ip_lists`

The desired_security_ip_lists supports the following:

* `group_name` - (Optional) Security group name.
* `security_ip_list` - (Optional) The IP list of Security group. Each single IP value should be Separated by comma.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DBInstance. 
* `region_id` - The region ID of the instance.
* `engine` - The engine of DBInstance. Always `selectdb`.
* `engine_minor_version` - The current DBInstance minor version.
* `status` - The status of the resource. Valid values: `ACTIVE`,`STOPPED`,`STARTING`,`RESTART`.
* `cpu_prepaid` - The sum of cpu resource amount for every `Subscription` clusters in DBInstance.
* `memory_prepaid` - The sum of memory resource amount offor every `Subscription` clusters in DBInstance.
* `cache_size_prepaid` - The sum of cache size for every `Subscription` clusters in DBInstance.
* `cluster_count_prepaid` - The sum of cluster counts for `Subscription` clusters in DBInstance.
* `cpu_postpaid` - The sum of cpu resource amount for every `PayAsYouGo` clusters in DBInstance.
* `memory_postpaid` - The sum of memory resource amount offor every `PayAsYouGo` clusters in DBInstance.
* `cache_size_postpaid` - The sum of cache size for every `PayAsYouGo` clusters in DBInstance.
* `cluster_count_postpaid` - The sum of cluster counts for `PayAsYouGo` clusters in DBInstance.
* `sub_domain` - The sub domain of DBInstance.
* `gmt_created` - The time when DBInstance is created.
* `gmt_modified` - The time when DBInstance is modified.
* `gmt_expired` - The time when DBInstance will be expired. Available on `Subscription` DBInstance.
* `lock_mode` - The lock mode of the instance. Set the value to lock, which specifies that the instance is locked when it automatically expires or has an overdue payment.
* `lock_reason` - The reason why the instance is locked.
* `instance_net_infos` - The net infos for instances.
  * `db_ip` - The IP address of the instance.
  * `vpc_instance_id` - The VPC ID.
  * `connection_string` - The connection string of the instance.
  * `net_type` - The network type of the instance.
  * `vswitch_id` - The ID of vswitch.
  * `port_list` - A list for port provides SelectDB service.
    * `protocol` - The protocol of the port.
    * `port` - The port that is used to connect.
* `security_ip_lists` - The details about each IP address whitelist returned. 
  * `group_name` - Security group name.
  * `security_ip_type` - The IP address type. Valid values: `ipv4`, `ipv6` (not supported).
  * `security_ip_list` - The IP list of Security group. Each single IP value should be Separated by comma.
  * `group_tag` - The tag of Security group.
  * `list_net_type` - The network type of Security group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the SelectDB DBInstance (until it reaches the initial `ACTIVATION` status).
* `update` - (Defaults to 30 mins) Used when update the SelectDB DBInstance.
* `delete` - (Defaults to 10 mins) Used when delete the SelectDB DBInstance.

## Import

SelectDB DBInstance can be imported using the id, e.g.

```shell
$ terraform import alicloud_selectdb_db_instance.example <id>
```
