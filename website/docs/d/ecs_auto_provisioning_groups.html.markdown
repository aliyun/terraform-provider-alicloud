---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_auto_provisioning_groups"
sidebar_current: "docs-alicloud-datasource-ecs-auto-provisioning-groups"
description: |-
  Provides a list of ECS auto provisioning groups.
---

# alicloud_ecs_auto_provisioning_groups

The `alicloud_ecs_auto_provisioning_groups` data source provides a collection of ECS auto provisioning groups in the current region. You can filter the results by ID, name, resource group, status, type and tags.

For information about the API and its fields, see [DescribeAutoProvisioningGroups](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-describeautoprovisioninggroups).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The data source reads every page returned by `DescribeAutoProvisioningGroups`, so the result is complete regardless of the number of groups. It starts at `page_number` and stops when the rows read reach the rows that are readable from that page or when the read has run for 30 minutes. The read fails instead of writing an incomplete result to state when a page ends early while rows that are readable from `page_number` are still missing, when a group ID repeats across pages, when the reported total changes between pages or disagrees with the rows read, when the API keeps returning full pages for 1,000 requests, or when the time budget is exhausted.

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

* `ids` - (Optional, Computed, List of String) A list of auto provisioning group IDs. A maximum of 20 entries may be specified. The `DescribeAutoProvisioningGroups` operation description limits the ID list to 20 entries even though the published schema metadata still declares `maxItems` 100; the provider follows the documented operation limit. Blank IDs are dropped and repeated IDs are sent once, so the result does not depend on the server's duplicate handling. The attribute also reports the IDs that matched.
* `auto_provisioning_group_name` - (Optional, String) The name of the auto provisioning group. The match is exact and is sent to the server.
* `resource_group_id` - (Optional, String) The ID of the resource group to which the auto provisioning group belongs.
* `status` - (Optional, List of String) A list of statuses to filter by. A maximum of 10 entries is allowed. Valid values: `submitted`, `active`, `deleted`, `delete-running` and `modifying`.
* `types` - (Optional, List of String) A list of auto provisioning group types to filter by. A maximum of 10 entries is allowed. Valid values: `request`, `maintain` and `candidate`.
* `tags` - (Optional, Map of String) A mapping of tags to filter by. A maximum of 20 entries is allowed. A key must be 1 to 128 characters, must not start with `aliyun` or `acs:`, and must not contain `http://` or `https://`. A value may be empty, is limited to 128 characters and must not contain `http://` or `https://`. A null value is rejected: omit the entry to leave the tag out of the request, or use an empty string to send an empty value.
* `page_number` - (Optional, Int) The number of the first page to read. Valid values: `1` to `2147483647` (`PageNumber` is an int32 on the wire). Defaults to `1`. The rows of the pages before `page_number` are not read and are not part of the result; `total_count` still reports them. The data source continues with the following pages until the rows that are readable from `page_number` have been read, or the read budget is exhausted.
* `page_size` - (Optional, Int) The number of entries per page. Valid values: `1` to `100`. Defaults to `100`. A page with fewer entries than this value ends the read only when the rows read account for every row that is readable from `page_number`; otherwise the read is reported as incomplete.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Computed, String) The ID of the data source, derived from the matched auto provisioning group IDs.
* `total_count` - (Computed, Int) The `TotalCount` reported by `DescribeAutoProvisioningGroups` for the configured filters. It counts every matching group in the region, not only the rows returned from `page_number`. A total that changes between pages, or that disagrees with the rows read, is reported as an error instead of being written. A response that does not report `TotalCount` falls back to the number of rows read.
* `names` - (Computed, List of String) A list of auto provisioning group names that match the configured filters.
* `groups` - (Computed, List of Object) A list of auto provisioning groups. Each element contains the following attributes:
  * `auto_provisioning_group_id` - (String) The ID of the auto provisioning group.
  * `auto_provisioning_group_name` - (String) The name of the auto provisioning group.
  * `auto_provisioning_group_type` - (String) The delivery type. Valid values: `request` and `maintain`.
  * `status` - (String) The status of the auto provisioning group. Valid values: `submitted`, `active`, `deleted`, `delete-running` and `modifying`.
  * `state` - (String) The scheduling state of the auto provisioning group. Valid values: `fulfilled`, `pending-fulfillment`, `pending-termination` and `error`.
  * `region_id` - (String) The ID of the region.
  * `resource_group_id` - (String) The ID of the resource group.
  * `launch_template_id` - (String) The ID of the launch template.
  * `launch_template_version` - (String) The version of the launch template.
  * `valid_from` - (String) The time when the auto provisioning group takes effect, in UTC and formatted as `2026-11-05T10:00:00Z`.
  * `valid_until` - (String) The time when the auto provisioning group expires, in UTC and formatted as `2026-11-05T10:00:00Z`.
  * `terminate_instances` - (Bool) Whether the instances in the auto provisioning group are released when the group is deleted. When `false`, the instances are retained.
  * `terminate_instances_with_expiration` - (Bool) Whether the instances in the auto provisioning group are released when the group expires. When `false`, the instances are only removed from the group and are not released.
  * `excess_capacity_termination_policy` - (String) The policy applied when the current capacity exceeds the target capacity and a scale-in is triggered. `termination` releases the scaled-in instances; `no-termination` only removes them from the auto provisioning group without releasing them. Valid values: `no-termination` and `termination`.
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
  * `candidate_options` - (List of Object) A single-element list with the candidate auto provisioning group options. Each element contains:
    * `timeout_minutes` - (Int) The timeout, in minutes, that the group tries to reach the target capacity.
  * `capacity_specification` - (List of Object) A single-element list with the current capacity distribution of the group. Each element contains:
    * `pay_as_you_go_capacity` - (String) The current capacity of pay-as-you-go instances.
    * `pre_paid_capacity` - (String) The current capacity of subscription instances.
    * `spot_capacity` - (String) The current capacity of preemptible instances.
    * `total_capacity` - (String) The current total capacity of the group.
  * `suspended_processes` - (List of String) The processes that are suspended for the group. Valid values: `ScaleIn` and `ScaleOut`.
