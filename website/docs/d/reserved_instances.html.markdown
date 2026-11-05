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

-> **NOTE:** The data source reads every page returned by `DescribeReservedInstances`, so the result is complete regardless of the number of reserved instances. The `operation_locks` member of the API response is not exported.

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

* `ids` - (Optional, Computed, List of String) A list of reserved instance IDs. A maximum of 100 entries is allowed. The filter is sent to the server.
* `reserved_instance_name` - (Optional, String) The name of the reserved instance. The match is exact and is sent to the server.
* `zone_id` - (Optional, String) The zone ID. The value is only effective when `scope` is `Zone`; when it is set without `scope`, the data source requests `Scope = Zone`.
* `lock_reason` - (Optional, String) The reason why the reserved instance is locked. Valid values: `financial` and `security`.
* `instance_type` - (Optional, String) The instance type. The value follows `ecs.<instance type family>.<specification>`, for example `ecs.g6.large`.
* `instance_type_family` - (Optional, String) The instance type family, for example `ecs.g6`.
* `scope` - (Optional, String) The scope of the reserved instance. Valid values: `Region` and `Zone`.
* `offering_type` - (Optional, String) The payment type. Valid values: `No Upfront`, `Partial Upfront` and `All Upfront`.
* `allocation_type` - (Optional, String) The allocation type. Valid values: `Normal` and `Shared`. Defaults to `Normal` on the server side.
* `status` - (Optional, List of String) A list of statuses to filter by. A maximum of 100 entries is allowed. Valid values: `Creating`, `Active`, `Retired`, `Updating` and `Expired`.
* `tags` - (Optional, Map of String) A mapping of tags to filter by. A maximum of 20 entries is allowed.
* `page_number` - (Optional, Int) The number of the first page to read. Valid values: `1` or greater. Defaults to `1`.
* `page_size` - (Optional, Int) The number of entries per page. Valid values: `1` to `100`. Defaults to `100`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Computed, String) The ID of the data source, derived from the matched reserved instance IDs.
* `total_count` - (Computed, Int) The number of reserved instances that match all configured filters.
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
  * `scope` - (String) The scope of the reserved instance. Valid values: `Region` and `Zone`.
  * `platform` - (String) The platform of the instances covered by the reserved instance, for example `Linux` or `Windows`.
  * `offering_type` - (String) The payment type. Valid values: `No Upfront`, `Partial Upfront` and `All Upfront`.
  * `allocation_status` - (String) The allocation status reported by the service.
  * `description` - (String) The description of the reserved instance.
  * `resource_group_id` - (String) The ID of the resource group.
  * `tags` - (Map of String) A mapping of tags.
