---
layout: "alicloud"
page_title: "Alicloud: alicloud_instance_types"
sidebar_current: "docs-alicloud-datasource-instance-types"
description: |-
    Provides a list of Ecs Instance Types for use in alicloud_instance resource.
---

# alicloud\_instance\_types

The Instance Types data source list the ecs_instance_types of Alicloud.

~> **NOTE:** Default to provide upgraded instance types. If you want to get outdated instance types, you should set `is_outdated` to true.

~> **NOTE:** If one instance type is sold out, it will not be exported.

## Example Usage

```
# Declare the data source
data "alicloud_instance_types" "1c2g" {
  cpu_core_count = 1
  memory_size = 2
}

# Create ecs instance with the first matched instance_type

resource "alicloud_instance" "instance" {
  instance_type = "${data.alicloud_instance_types.1c2g.instance_types.0.id}"

  # Other properties...
}

```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The Zone that supports available instance types.
* `cpu_core_count` - (Optional) Limit search to specific cpu core count.
* `memory_size` - (Optional) Limit search to specific memory size.
* `instance_type_family` - (Optional) Allows to filter list of Instance Types based on their
family name, for example 'ecs.n4'.
* `instance_charge_type` - (Optional) According to ECS instance charge type to filter all results. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `network_type` - (Optional) According to network type to filter all results. Valid values: `Classic` and `Vpc`.
* `spot_strategy` - - (Optional) According to ECS spot type to filter all results. Valid values: `NoSpot`, `SpotWithPriceLimit` and `SpotAsPriceGo`. Default to `NoSpot`.
* `is_outdated` - (Optional) Whether to export outdated instance types. Default to false.
* `output_file` - (Optional) The name of file that can save instance types data source after running `terraform plan`.

## Attributes Reference

A list of instance types will be exported and its every element contains the following attributes:

* `id` - ID of the instance type.
* `cpu_core_count` - Number of CPU cores.
* `memory_size` - Size of memory, measured in GB.
* `family` - The instance type family.
* `availability_zones` - List of availability zones which support the instance types.
* `gpu` - The GPU attribution of an instance type:
    * `amount` - The amount of GPU of an instance type.
    * `category` - The category of GPU of an instance type.

* `burstable_instance` - The burstable instance's attribution:
    * `initial_credit` - The initial CPU credit of a burstable instance
    * `baseline_credit` - The compute performance benchmark CPU credit of a burstable instance

* `eni_amount` - The maximum number of network interface that an instance type can be attached to.
* `local_storage` - Local storage of an instance type:
    * `capacity` - The capacity of a local storage
    * `amount` - The number of local storages that an instance has been attached to
    * `category` - The category of local storage that an instance has been attached to