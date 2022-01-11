---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_access_strategies"
sidebar_current: "docs-alicloud-datasource-alidns-access-strategies"
description: |-
  Provides a list of Alidns Access Strategies to the user.
---

# alicloud\_alidns\_access\_strategies

This data source provides the Alidns Access Strategies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alidns_access_strategies" "ids" {
  instance_id   = "example_value"
  strategy_mode = "example_value"
  ids           = ["example_value-1", "example_value-2"]
  name_regex    = "the_resource_name"
}
output "alidns_access_strategy_id_1" {
  value = data.alicloud_alidns_access_strategies.ids.strategies.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Access Strategy IDs.
* `instance_id` - (Required, ForceNew) The Id of the associated instance.
* `lang` - (Optional, ForceNew) The lang.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Access Strategy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `strategy_mode` - (Required, ForceNew) The type of the access policy. Valid values:
  - `GEO`: based on geographic location. 
  - `LATENCY`: Based on delay.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Access Strategy names.
* `strategies` - A list of Alidns Access Strategies. Each element contains the following attributes:
 * `access_strategy_id` - The first ID of the resource.
 * `create_time` - The time when the access policy was created.
 * `create_timestamp` - The timestamp that indicates when the access policy was created.
 * `default_addr_pool_type` - The type of the primary address pool.
 * `default_addr_pools` - The address pools in the primary address pool group.
  * `addr_pool_id` - The ID of the address pool.
  * `lba_weight` - The weight of the address pool.
  * `name` - The name of the address pool.
  * `addr_count` - The number of addresses in the address pool.
 * `default_available_addr_num` - The number of addresses currently available in the primary address pool.
 * `default_latency_optimization` - Indicates whether scheduling optimization for latency resolution was enabled for the primary address pool group.
 * `default_lba_strategy` - The load balancing policy of the primary address pool group.
 * `default_max_return_addr_num` - The maximum number of addresses returned by the primary address pool set.
 * `default_min_available_addr_num` - The minimum number of available addresses for the primary address pool set.
 * `effective_addr_pool_group_type` - The type of the active address pool group.
 * `failover_addr_pool_type` - The type of the secondary address pool.
 * `failover_addr_pools` - The address pools in the secondary address pool group.
  * `lba_weight` - The weight of the address pool.
  * `name` - The name of the address pool.
  * `addr_count` - The number of addresses in the address pool.
  * `addr_pool_id` - The ID of the address pool.
 * `failover_available_addr_num` - The number of available addresses in the standby address pool.
 * `failover_latency_optimization` - Indicates whether scheduling optimization for latency resolution was enabled for the secondary address pool group.
 * `failover_lba_strategy` - The load balancing policy of the secondary address pool group.
 * `failover_max_return_addr_num` - The maximum number of returned addresses in the standby address pool.
 * `failover_min_available_addr_num` - The minimum number of available addresses in the standby address pool.
 * `id` - The ID of the Access Strategy.
 * `instance_id` - The Id of the associated instance.
 * `lines` - List of source regions.
  * `group_code` - The code of the source region group.
  * `group_name` - The name of the source region group.
  * `line_code` - The line code of the source region.
  * `line_name` - The line name of the source region.
 * `strategy_mode` - The type of the access policy.
 * `strategy_name` - The name of the access policy.
 * `access_mode` - The primary/secondary switchover policy for address pool groups.