---
layout: "alicloud"
page_title: "Alicloud: alicloud_instance_types"
sidebar_current: "docs-alicloud-datasource-instance-types"
description: |-
    Provides a list of ECS Instance Types to be used by the alicloud_instance resource.
---

# alicloud\_instance\_types

This data source provides the ECS instance types of Alibaba Cloud.

~> **NOTE:** By default, only the upgraded instance types are returned. If you want to get outdated instance types, you must set `is_outdated` to true.

~> **NOTE:** If one instance type is sold out, it will not be exported.

## Example Usage

```
# Declare the data source
data "alicloud_instance_types" "types_ds" {
  cpu_core_count = 1
  memory_size = 2
}

# Create ECS instance with the first matched instance_type

resource "alicloud_instance" "instance" {
  instance_type = "${data.alicloud_instance_types.types_ds.instance_types.0.id}"

  # Other properties...
}

```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The zone where instance types are supported.
* `cpu_core_count` - (Optional) Filter the results to a specific number of cpu cores.
* `memory_size` - (Optional) Filter the results to a specific memory size in GB.
* `instance_type_family` - (Optional) Filter the results based on their family name. For example: 'ecs.n4'.
* `instance_charge_type` - (Optional) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `network_type` - (Optional) Filter the results by network type. Valid values: `Classic` and `Vpc`.
* `spot_strategy` - - (Optional) Filter the results by ECS spot type. Valid values: `NoSpot`, `SpotWithPriceLimit` and `SpotAsPriceGo`. Default to `NoSpot`.
* `is_outdated` - (Optional, type: bool) If true, outdated instance types are included in the results. Default to false.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instance_types` - A list of image types. Each element contains the following attributes:
  * `id` - ID of the instance type.
  * `cpu_core_count` - Number of CPU cores.
  * `memory_size` - Size of memory, measured in GB.
  * `family` - The instance type family.
  * `availability_zones` - List of availability zones that support the instance type.
  * `gpu` - The GPU attribution of an instance type:
    * `amount` - The amount of GPU of an instance type.
    * `category` - The category of GPU of an instance type.
  * `burstable_instance` - The burstable instance attribution:
    * `initial_credit` - The initial CPU credit of a burstable instance.
    * `baseline_credit` - The compute performance benchmark CPU credit of a burstable instance.
  * `eni_amount` - The maximum number of network interfaces that an instance type can be attached to.
  * `local_storage` - Local storage of an instance type:
    * `capacity` - The capacity of a local storage in GB.
    * `amount` - The number of local storage devices that an instance has been attached to.
    * `category` - The category of local storage that an instance has been attached to.