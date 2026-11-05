---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_reserved_instances"
sidebar_current: "docs-alicloud-datasource-reserved-instances"
description: |-
  Provides a list of ECS reserved instances.
---

# alicloud_reserved_instances

The `alicloud_reserved_instances` data source provides a collection of ECS reserved instances (reserved instance coupons) in the current region. You can filter the results by ID, name, zone, instance type, scope, payment type, allocation type, status, lock reason and tags.

For information about the API and its fields, see [DescribeReservedInstances](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-describereservedinstances).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The data source reads every page returned by `DescribeReservedInstances`, so the result is complete regardless of the number of reserved instances. It starts at `page_number` and stops when the rows read reach the rows that are readable from that page or when the read has run for 30 minutes. The read fails instead of writing an incomplete result to state when a page ends early while rows that are readable from `page_number` are still missing, when a resource ID repeats across pages, when the reported total changes between pages or disagrees with the rows read, when the API keeps returning full pages for 1,000 requests, or when the time budget is exhausted.

-> **NOTE:** Filtering by `tags` returns at most 1,000 reserved instances. Use [ListTagResources](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-listtagresources) to enumerate the tagged reserved instances beyond that limit.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

data "alicloud_reserved_instances" "default" {
  status = ["Active"]
  scope  = "Region"
}

output "reserved_instance_ids" {
  value = data.alicloud_reserved_instances.default.ids
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Computed, List of String) A list of reserved instance IDs. A maximum of 100 entries is allowed. Blank and repeated IDs are sent once.
* `reserved_instance_name` - (Optional, String) The name of the reserved instance. The match is exact and is sent to the server.
* `zone_id` - (Optional, String) The zone ID. It is required when `scope` is `Zone`, and the read is rejected without it. When `zone_id` is set without `scope`, the data source requests `Scope = Zone`; when neither is set, the query is not restricted to a scope.
* `lock_reason` - (Optional, String) The reason why the reserved instance is locked. Valid values: `financial` and `security`.
* `instance_type` - (Optional, String) The instance type. The value follows `ecs.<instance type family>.<specification>`, for example `ecs.g6.large`.
* `instance_type_family` - (Optional, String) The instance type family, for example `ecs.g6`.
* `scope` - (Optional, String) The scope of the reserved instance. Valid values: `Region` and `Zone`. `zone_id` is required when the value is `Zone`.
* `offering_type` - (Optional, String) The payment type. Valid values: `No Upfront`, `Partial Upfront` and `All Upfront`.
* `allocation_type` - (Optional, String) The allocation type. Valid values: `Normal` and `Shared`. Defaults to `Normal` on the server side.
* `status` - (Optional, List of String) A list of statuses to filter by. A maximum of 100 entries is allowed. Valid values: `Creating`, `Active`, `Retired`, `Updating` and `Expired`. **Note:** `Expired` is deprecated by `DescribeReservedInstances` and may be removed from the API. It is still accepted for compatibility; new configurations should filter with the other values.
* `tags` - (Optional, Map of String) A mapping of tags to filter by. A maximum of 20 entries is allowed. A key and a value must be 1 to 128 characters, so an empty value is rejected; a key must not start with `aliyun` or `acs:`, a value must not start with `acs:`, and neither may contain `http://` or `https://`. A null value is rejected as well: omit the entry to leave the tag out of the request. A tag filter returns at most 1,000 reserved instances; when more may match, query `ListTagResources` instead.
* `page_number` - (Optional, Int) The number of the first page to read. Valid values: `1` or greater. Defaults to `1`. The rows of the pages before `page_number` are not read and are not part of the result; `total_count` still reports them. The data source continues with the following pages until the rows that are readable from `page_number` have been read, or the read budget is exhausted.
* `page_size` - (Optional, Int) The number of entries per page. Valid values: `1` to `100`. Defaults to `100`. A page with fewer entries than this value ends the read only when the rows read account for every row that is readable from `page_number`; otherwise the read is reported as incomplete.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Computed, String) The ID of the data source, derived from the matched reserved instance IDs.
* `total_count` - (Computed, Int) The `TotalCount` reported by `DescribeReservedInstances` for the configured filters. It counts every matching reserved instance in the region, not only the rows returned from `page_number`. A total that changes between pages, or that disagrees with the rows read, is reported as an error instead of being written.
* `names` - (Computed, List of String) A list of reserved instance names that match the configured filters.
* `instances` - (Computed, List of Object) A list of reserved instances. Each element contains the following attributes:
  * `reserved_instance_id` - (String) The ID of the reserved instance.
  * `reserved_instance_name` - (String) The name of the reserved instance.
  * `status` - (String) The status of the reserved instance. Valid values: `Creating`, `Active`, `Retired`, `Updating` and `Expired`.
  * `creation_time` - (String) The time when the reserved instance was created, in UTC and formatted as `2026-11-05T10:00:00Z`.
  * `start_time` - (String) The time when the reserved instance takes effect, in UTC and formatted as `2026-11-05T10:00:00Z`.
  * `expired_time` - (String) The time when the reserved instance expires, in UTC and formatted as `2026-11-05T10:00:00Z`.
  * `instance_type` - (String) The instance type, for example `ecs.g6.large`.
  * `instance_amount` - (Int) The number of instances that the reserved instance covers.
  * `region_id` - (String) The ID of the region.
  * `zone_id` - (String) The ID of the zone.
  * `scope` - (String) The scope of the reserved instance as returned by the service, for example `region`. Use the request-side `scope` values `Region` and `Zone` when filtering.
  * `platform` - (String) The platform of the instances covered by the reserved instance, for example `Linux` or `Windows`.
  * `offering_type` - (String) The payment type. Valid values: `No Upfront`, `Partial Upfront` and `All Upfront`.
  * `allocation_status` - (String) The allocation status of the reserved instance. This value is meaningful only when `allocation_type` is `Shared`. Valid values: `allocated` (allocated to another user) and `beAllocated` (allocated by another user).
  * `description` - (String) The description of the reserved instance.
  * `resource_group_id` - (String) The ID of the resource group.
  * `operation_locks` - (List of Object) The locks applied to the reserved instance. Each element contains:
    * `lock_reason` - (String) The reason why the reserved instance is locked. Valid values: `financial` and `security`.
  * `tags` - (Map of String) A mapping of tags.
