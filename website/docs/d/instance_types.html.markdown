---
subcategory: "ECS"
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
  memory_size    = 2
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
* `gpu_amount` - (Optional, Available in 1.69.0+) The GPU amount of an instance type.
* `gpu_spec` - (Optional, Available in 1.69.0+) The GPU spec of an instance type.
* `instance_type_family` - (Optional) Filter the results based on their family name. For example: 'ecs.n4'.
* `instance_charge_type` - (Optional) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `network_type` - (Optional) Filter the results by network type. Valid values: `Classic` and `Vpc`.
* `spot_strategy` - (Optional) Filter the results by ECS spot type. Valid values: `NoSpot`, `SpotWithPriceLimit` and `SpotAsPriceGo`. Default to `NoSpot`.
* `eni_amount` - (Optional) Filter the result whose network interface number is no more than `eni_amount`.
* `kubernetes_node_role` - (Optional) Filter the result which is used to create a [kubernetes cluster](https://www.terraform.io/docs/providers/alicloud/r/cs_kubernetes)
 and [managed kubernetes cluster](https://www.terraform.io/docs/providers/alicloud/r/cs_managed_kubernetes). Optional Values: `Master` and `Worker`.
* `is_outdated` - (Optional, type: bool) If true, outdated instance types are included in the results. Default to false.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `system_disk_category` - (Optional, Available in 1.120.0+) Filter the results by system disk category. Valid values: `cloud`, `ephemeral_ssd`, `cloud_essd`, `cloud_efficiency`, `cloud_ssd`. 
  **NOTE**: Its default value `cloud_efficiency` has been removed from the version v1.150.0.
* `image_id` - (Optional, Available in 1.163.0+) The ID of the image.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance type IDs.
* `instance_types` - A list of image types. Each element contains the following attributes:
  * `id` - ID of the instance type.
  * `cpu_core_count` - Number of CPU cores.
  * `memory_size` - Size of memory, measured in GB.
  * `family` - The instance type family.
  * `availability_zones` - List of availability zones that support the instance type.
  * `nvme_support` - Indicates whether the cloud disk can be attached by using the nonvolatile memory express (NVMe) protocol. Valid values:
    - required: The cloud disk can be attached by using the NVMe protocol.
    - unsupported: The cloud disk cannot be attached by using the NVMe protocol.  
  * `gpu` - The GPU attribution of an instance type:
    - amount: The amount of GPU of an instance type.
    - category: The category of GPU of an instance type.
  * `burstable_instance` - The burstable instance attribution:
    - initial_credit: The initial CPU credit of a burstable instance.
    - baseline_credit:  The compute performance benchmark CPU credit of a burstable instance.
  * `eni_amount` - The maximum number of network interfaces that an instance type can be attached to.
  * `local_storage` - Local storage of an instance type:
    - capacity: The capacity of a local storage in GB.
    - amount:  The number of local storage devices that an instance has been attached to.
    - category: The category of local storage that an instance has been attached to.
