---
subcategory: "SelectDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_selectdb_db_instances"
sidebar_current: "docs-alicloud-datasource-selectdb-db-instances"
description: |-
  Provides a list of SelectDB DBInstance to the user.
---

# alicloud_selectdb_db_instances

This data source provides the SelectDB DBInstance of the current Alibaba Cloud user.

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
  cache_size              = 200
  payment_type            = "PayAsYouGo"
  engine_minor_version    = "3.0.12"
  vpc_id                  = data.alicloud_vswitches.default.vswitches.0.vpc_id
  zone_id                 = data.alicloud_vswitches.default.vswitches.0.zone_id
  vswitch_id              = data.alicloud_vswitches.default.vswitches.0.id
}

data "alicloud_selectdb_db_instances" "default" {
  ids = [alicloud_selectdb_db_instance.default.id]
}
output "db_instance" {
  value = data.alicloud_selectdb_db_instances.default.ids.0
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of DBInstance IDs.
* `tags` - (Optional) A mapping of tags to assign to the resource. Used for instance searching.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of SelectDB DBInstance. Each element contains the following attributes:
  * `id` - The resource ID in terraform of DBInstance.
  * `db_instance_id` - The instance ID.
  * `engine` - The Engine of the DBInstance.
  * `engine_version` - The engine version of the DBInstance.
  * `engine_minor_version` - The engine minor version of the DBInstance.
  * `db_instance_description` - The DBInstance description.
  * `status` - The status of the DBInstance. Valid values: `ACTIVATION`,`CREATING`,`DELETING`,`RESTARTING`,`ORDER_PREPARING`.
  * `payment_type` - The payment type of the resource. Valid values: `PayAsYouGo`,`Subscription`.
  * `cpu_prepaid` - The sum of cpu resource amount for every `Subscription` clusters in DBInstance.
  * `memory_prepaid` - The sum of memory resource amount offor every `Subscription` clusters in DBInstance.
  * `cache_size_prepaid` - The sum of cache size for every `Subscription` clusters in DBInstance.
  * `cluster_count_prepaid` - The sum of cluster counts for `Subscription` clusters in DBInstance.
  * `cpu_postpaid` - The sum of cpu resource amount for every `PayAsYouGo` clusters in DBInstance.
  * `memory_postpaid` - The sum of memory resource amount offor every `PayAsYouGo` clusters in DBInstance.
  * `cache_size_postpaid` - The sum of cache size for every `PayAsYouGo` clusters in DBInstance.
  * `cluster_count_postpaid` - The sum of cluster counts for `PayAsYouGo` clusters in DBInstance.
  * `region_id` - The ID of region for DBInstance.
  * `zone_id` - The ID of zone for DBInstance.
  * `vpc_id` - The ID of the VPC for DBInstance.
  * `vswitch_id` - The ID of vswitch for DBInstance.
  * `sub_domain` - The sub domain of DBInstance.
  * `gmt_created` - The time when DBInstance is created.
  * `gmt_modified` - The time when DBInstance is modified.
  * `gmt_expired` - The time when DBInstance will be expired. Available on `Subscription` DBInstance.
  * `lock_mode` - The lock mode of the instance. Set the value to lock, which specifies that the instance is locked when it automatically expires or has an overdue payment.
  * `lock_reason` - The reason why the instance is locked.
