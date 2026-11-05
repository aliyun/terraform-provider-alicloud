---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_auto_provisioning_groups"
sidebar_current: "docs-alicloud-datasource-ecs-auto-provisioning-groups"
description: |-
  Provides a list of ECS auto provisioning groups.
---

# alicloud\_ecs\_auto\_provisioning\_groups

The `alicloud_ecs_auto_provisioning_groups` data source provides a collection of ECS auto provisioning groups in the current region. You can filter the results by ID, name, resource group, status, type and tags.

For information about the API and its fields, see [DescribeAutoProvisioningGroups](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-describeautoprovisioninggroups).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The data source reads every page returned by `DescribeAutoProvisioningGroups`, so the result is complete regardless of the number of groups.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

data "alicloud_ecs_auto_provisioning_groups" "default" {
  status = ["active"]
  types  = ["maintain"]
}

output "auto_provisioning_group_ids" {
  value = data.alicloud_ecs_auto_provisioning_groups.default.ids
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Computed, List of String) A list of auto provisioning group IDs. A maximum of 100 entries is allowed. The filter is sent to the server.
* `auto_provisioning_group_name` - (Optional, String) The name of the auto provisioning group. The match is exact and is sent to the server.
* `resource_group_id` - (Optional, String) The ID of the resource group to which the auto provisioning group belongs.
* `status` - (Optional, List of String) A list of statuses to filter by. A maximum of 10 entries is allowed. Valid values: `submitted`, `active`, `deleted`, `delete-running` and `modifying`.
* `types` - (Optional, List of String) A list of auto provisioning group types to filter by. A maximum of 10 entries is allowed. Valid values: `request`, `maintain` and `candidate`.
* `tags` - (Optional, Map of String) A mapping of tags to filter by. A maximum of 20 entries is allowed.
* `page_number` - (Optional, Int) The number of the first page to read. Valid values: `1` or greater. Defaults to `1`.
* `page_size` - (Optional, Int) The number of entries per page. Valid values: `1` to `100`. Defaults to `100`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Computed, String) The ID of the data source, derived from the matched auto provisioning group IDs.
* `total_count` - (Computed, Int) The number of auto provisioning groups that match all configured filters.
* `names` - (Computed, List of String) A list of auto provisioning group names that match the configured filters.
* `groups` - (Computed, List of Object) A list of auto provisioning groups. Each element contains the following attributes:
  * `auto_provisioning_group_id` - (String) The ID of the auto provisioning group.
  * `auto_provisioning_group_name` - (String) The name of the auto provisioning group.
  * `auto_provisioning_group_type` - (String) The type of the auto provisioning group. Valid values: `request`, `maintain` and `candidate`.
  * `status` - (String) The status of the auto provisioning group. Valid values: `submitted`, `active`, `deleted`, `delete-running` and `modifying`.
  * `state` - (String) The lifecycle state reported by the service, for example `Active`.
  * `region_id` - (String) The ID of the region.
  * `resource_group_id` - (String) The ID of the resource group.
  * `launch_template_id` - (String) The ID of the launch template.
  * `launch_template_version` - (String) The version of the launch template.
  * `valid_from` - (String) The time when the auto provisioning group takes effect, in UTC and formatted as `2026-11-05T10:00:00Z`.
  * `valid_until` - (String) The time when the auto provisioning group expires, in UTC and formatted as `2026-11-05T10:00:00Z`.
  * `terminate_instances` - (Bool) Whether instances of the auto provisioning group are released.
  * `terminate_instances_with_expiration` - (Bool) Whether preemptible instances are shut down when the auto provisioning group expires.
  * `excess_capacity_termination_policy` - (String) The shutdown policy for excess preemptible instances. Valid values: `no-termination` and `termination`.
  * `max_spot_price` - (String) The global maximum price for preemptible instances, as returned by the service.
  * `creation_time` - (String) The time when the auto provisioning group was created, in UTC and formatted as `2026-11-05T10:00:00Z`.
  * `tags` - (Map of String) A mapping of tags.
  * `launch_template_configs` - (List of Object) A list of extended configurations of the launch template. Each element contains:
    * `instance_type` - (String) The instance type.
    * `max_price` - (String) The maximum price.
    * `vswitch_id` - (String) The ID of the vSwitch.
    * `weighted_capacity` - (String) The weight of the instance type.
    * `priority` - (String) The priority of the instance type.
  * `target_capacity_specification` - (List of Object) A single-element list describing the target capacity. Each element contains:
    * `default_target_capacity_type` - (String) The capacity type used for supplemental instances. Valid values: `PayAsYouGo` and `Spot`.
    * `total_target_capacity` - (String) The total target capacity.
    * `pay_as_you_go_target_capacity` - (String) The target capacity of pay-as-you-go instances.
    * `spot_target_capacity` - (String) The target capacity of preemptible instances.
  * `spot_options` - (List of Object) A single-element list with the preemptible instance options. Each element contains:
    * `allocation_strategy` - (String) The scale-out policy. Valid values: `lowest-price` and `diversified`.
    * `instance_interruption_behavior` - (String) The behavior after a preemptible instance is shut down. Valid values: `stop` and `terminate`.
    * `instance_pools_to_use_count` - (Int) The number of instance pools to use when the strategy is `lowest-price`.
  * `pay_as_you_go_options` - (List of Object) A single-element list with the pay-as-you-go options. Each element contains:
    * `allocation_strategy` - (String) The scale-out policy. Valid values: `lowest-price` and `prioritized`.
