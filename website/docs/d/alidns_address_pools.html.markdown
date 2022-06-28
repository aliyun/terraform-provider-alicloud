---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_address_pools"
sidebar_current: "docs-alicloud-datasource-alidns-address-pools"
description: |-
  Provides a list of Alidns Address Pools to the user.
---

# alicloud\_alidns\_address\_pools

This data source provides the Alidns Address Pools of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alidns_address_pools" "ids" {
  instance_id = "example_value"
  ids         = ["example_value-1", "example_value-2"]
}
output "alidns_address_pool_id_1" {
  value = data.alicloud_alidns_address_pools.ids.pools.0.id
}

data "alicloud_alidns_address_pools" "nameRegex" {
  instance_id = "example_value"
  name_regex  = "^my-AddressPool"
}
output "alidns_address_pool_id_2" {
  value = data.alicloud_alidns_address_pools.nameRegex.pools.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Address Pool IDs.
* `instance_id` - (Required, ForceNew) The id of the instance.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Address Pool name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Address Pool names.
* `pools` - A list of Alidns Address Pools. Each element contains the following attributes:
  * `address_pool_id` - The first ID of the resource.
  * `address_pool_name` - The name of the address pool.
  * `create_time` - The time when the address pool was created.
  * `create_timestamp` - The timestamp that indicates when the address pool was created.
  * `id` - The ID of the Address Pool.
  * `instance_id` - The id of the instance.
  * `lba_strategy` - The load balancing policy of the address pool.
  * `type` - The type of the address pool.
  * `update_time` - The time when the address pool was updated.
  * `update_timestamp` - The timestamp that indicates when the address pool was updated.
  * `monitor_status` - Indicates whether health checks are configured.
  * `monitor_config_id` - The ID of the health check task.
  * `address` - The address lists of the Address Pool.
    * `address` - The address that you want to add to the address pool.
    * `attribute_info` - The source region of the address.
    * `lba_weight` - The weight of the address.
    * `mode` - The type of the address.
    * `remark` - The description of the address.