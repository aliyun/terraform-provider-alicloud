---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_instance_types"
sidebar_current: "docs-alicloud-datasource-instance-types"
description: |-
    Provides a list of ECS Instance Types to be used by the alicloud_instance resource.
---

# alicloud_instance_types

This data source provides the ECS instance types of Alibaba Cloud.

~> **NOTE:** By default, only the upgraded instance types are returned. If you want to get outdated instance types, you must set `is_outdated` to true.

~> **NOTE:** If one instance type is sold out, it will not be exported.

-> **NOTE:** Available since v1.0.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

# Declare the data source
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.default.id
  security_group_ids     = [alicloud_security_group.default.id]
}

resource "alicloud_instance" "default" {
  count                      = 14
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, ForceNew) The zone where instance types are supported.
* `cpu_core_count` - (Optional, ForceNew) Filter the results to a specific number of cpu cores.
* `memory_size` - (Optional, ForceNew) Filter the results to a specific memory size in GB.
* `sorted_by` - (Optional, ForceNew) Sort mode, valid values: `CPU`, `Memory`, `Price`.
* `gpu_amount` - (Optional, ForceNew, Available in 1.69.0+) The GPU amount of an instance type.
* `gpu_spec` - (Optional, ForceNew, Available in 1.69.0+) The GPU spec of an instance type.
* `instance_type_family` - (Optional, ForceNew) Filter the results based on their family name. For example: 'ecs.n4'.
* `instance_type` - (Optional, ForceNew, Available since 1.222.0) Instance specifications. For more information, see instance Specification Family, or you can call the describe instance types interface to get the latest specification table.
* `instance_charge_type` - (Optional, ForceNew) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `network_type` - (Optional, ForceNew) Filter the results by network type. Valid values: `Classic` and `Vpc`.
* `spot_strategy` - (Optional, ForceNew) Filter the results by ECS spot type. Valid values: `NoSpot`, `SpotWithPriceLimit` and `SpotAsPriceGo`. Default to `NoSpot`.
* `eni_amount` - (Optional, ForceNew) Filter the result whose network interface number is no more than `eni_amount`.
* `kubernetes_node_role` - (Optional, ForceNew) Filter the result which is used to create a [kubernetes cluster](https://www.terraform.io/docs/providers/alicloud/r/cs_kubernetes)
 and [managed kubernetes cluster](https://www.terraform.io/docs/providers/alicloud/r/cs_managed_kubernetes). Optional Values: `Master` and `Worker`.
* `is_outdated` - (Optional, type: bool) If true, outdated instance types are included in the results. Default to false.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `system_disk_category` - (Optional, ForceNew, Available since 1.120.0) Filter the results by system disk category. Valid values: `cloud`, `ephemeral_ssd`, `cloud_essd`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd_entry`, `cloud_auto`. 
  **NOTE**: Its default value `cloud_efficiency` has been removed from the version v1.150.0.
* `image_id` - (Optional, ForceNew, Available in 1.163.0+) The ID of the image.
* `minimum_eni_ipv6_address_quantity` (Optional, ForceNew, Available since 1.193.0) The minimum number of IPv6 addresses per ENI. **Note:** If an instance type supports fewer IPv6 addresses per ENI than the specified value, information about the instance type is not queried.
* `minimum_eni_private_ip_address_quantity` (Optional, ForceNew, Available since 1.223.1) The minimum expected IPv4 address upper limit of a single ENI when querying instance specifications. **Note:** If an instance type supports fewer IPv4 addresses per ENI than the specified value, information about the instance type is not queried.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance type IDs.
* `instance_types` - A list of image types. Each element contains the following attributes:
  * `id` - ID of the instance type.
  * `price` - The price of instance type.
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
  * `eni_amount` - (Deprecated since v1.239.0) The maximum number of ENIs per instance. It sames as `eni_quantity`.
  * `eni_quantity` - (Available since v1.239.0) The maximum number of ENIs per instance.
  * `primary_eni_queue_number` - (Available since v1.239.0) The default number of queues per primary ENI.
  * `secondary_eni_queue_number` - (Available since v1.239.0) The default number of queues per secondary ENI.
  * `eni_ipv6_address_quantity` - (Available since v1.239.0) The maximum number of IPv6 addresses per ENI. 
  * `maximum_queue_number_per_eni` - (Available since v1.239.0) The maximum number of queues per ENI, including primary and secondary ENIs.
  * `total_eni_queue_quantity` - (Available since v1.239.0) The maximum number of queues on ENIs that the instance type supports. 
  * `eni_private_ip_address_quantity` - (Available since v1.239.0) The maximum number of IPv4 addresses per ENI.
  * `local_storage` - Local storage of an instance type:
    - capacity: The capacity of a local storage in GB.
    - amount:  The number of local storage devices that an instance has been attached to.
    - category: The category of local storage that an instance has been attached to.
