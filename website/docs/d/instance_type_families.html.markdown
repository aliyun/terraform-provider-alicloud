---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_instance_type_families"
sidebar_current: "docs-alicloud-datasource-instance-type-families"
description: |-
    Provides a list of ECS Instance Type Families to be used by the alicloud_instance resource.
---

# alicloud\_instance\_type\_families

This data source provides the ECS instance type families of Alibaba Cloud.

-> **NOTE:** Available in 1.54.0+

## Example Usage

```
data "alicloud_instance_type_families" "default" {
  instance_charge_type = "PrePaid"
}

output "first_instance_type_family_id" {
  value = "${data.alicloud_instance_type_families.default.families.0.id}"
}

output "instance_ids" {
  value = "${data.alicloud_instance_type_families.default.ids}"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional, ForceNew) The Zone to launch the instance.
* `generation` - (Optional) The generation of the instance type family, Valid values: `ecs-1`, `ecs-2`, `ecs-3` and `ecs-4`. For more information, see [Instance type families](https://www.alibabacloud.com/help/doc-detail/25378.htm). 
* `instance_charge_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`, Default to `PostPaid`.
* `spot_strategy` - (Optional, ForceNew) Filter the results by ECS spot type. Valid values: `NoSpot`, `SpotWithPriceLimit` and `SpotAsPriceGo`. Default to `NoSpot`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance type family IDs.
* `instance_types` - A list of image type families. Each element contains the following attributes:
  * `id` - ID of the instance type family.
  * `generation` - The generation of the instance type family.
  * `zone_ids` - A list of Zone to launch the instance.
 
