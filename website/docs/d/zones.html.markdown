---
layout: "alicloud"
page_title: "Alicloud: alicloud_zones"
sidebar_current: "docs-alicloud-datasource-zones"
description: |-
    Provides a list of availability zones that can be used by an Alibaba Cloud account.
---

# alicloud\_zones

This data source provides availability zones that can be accessed by an Alibaba Cloud account within the region configured in the provider.


~> **NOTE:** If one zone is sold out, it will not be exported.

## Example Usage

```
# Declare the data source
data "alicloud_zones" "zones_ds" {
  "available_instance_type" = "ecs.n4.large"
  "available_disk_category" = "cloud_ssd"
}

# Create an ECS instance with the first matched zone
resource "alicloud_instance" "instance" {
  availability_zone = "${data.alicloud_zones.zones_ds.zones.0.id}"

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `available_instance_type` - (Optional) Filter the results by a specific instance type.
* `available_resource_creation` - (Optional) Filter the results by a specific resource type. The following values are allowed: `Instance`, `Disk`, `VSwitch` and `Rds`.
* `available_disk_category` - (Optional) Filter the results by a specific disk category. Can be either `cloud`, `cloud_efficiency` or `cloud_ssd`.
* `multi` - (Optional, type: bool) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch RDS instances.
* `instance_charge_type` - (Optional) Filter the results by a specific ECS instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `network_type` - (Optional) Filter the results by a specific network type. Valid values: `Classic` and `Vpc`.
* `spot_strategy` - - (Optional) Filter the results by a specific ECS spot type. Valid values: `NoSpot`, `SpotWithPriceLimit` and `SpotAsPriceGo`. Default to `NoSpot`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

~> **NOTE:** The disk category `cloud` has been outdated and can only be used by non-I/O Optimized ECS instances. Many availability zones don't support it. It is recommended to use `cloud_efficiency` or `cloud_ssd`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `local_name` - Name of the zone in the local language.
  * `available_instance_types` - Allowed instance types.
  * `available_resource_creation` - Type of resources that can be created.
  * `available_disk_categories` - Set of supported disk categories.
