---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_configurations"
sidebar_current: "docs-alicloud_ess_scaling_configurations"
description: |-
    Provides a list of scaling configurations available to the user.
---

# alicloud_ess_scaling_configurations

This data source provides available scaling configuration resources. 

-> **NOTE:** Available since v1.240.0

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

locals {
  name = "${var.name}-${random_integer.default.result}"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = local.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = local.name
}

resource "alicloud_security_group" "default" {
  security_group_name = local.name
  vpc_id              = alicloud_vpc.default.id
}

resource "alicloud_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "172.16.0.0/24"
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = local.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id           = alicloud_ess_scaling_group.default.id
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  security_group_id          = alicloud_security_group.default.id
  force_delete               = true
  active                     = true
  scaling_configuration_name = "scaling_configuration_name"
}

data "alicloud_ess_scaling_configurations" "scalingconfigurations_ds" {
  scaling_group_id = alicloud_ess_scaling_group.default.id
  ids              = [alicloud_ess_scaling_configuration.default.id]
  name_regex       = "scaling_configuration_name"
}

output "first_scaling_configuration" {
  value = data.alicloud_ess_scaling_configurations.scalingconfigurations_ds.configurations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional, ForceNew) Scaling group id the scaling configurations belong to.
* `name_regex` - (Optional, ForceNew) A regex string to filter resulting scaling configurations by name.
* `ids` - (Optional, ForceNew) A list of scaling configuration IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of scaling configuration ids.
* `names` - A list of scaling configuration names.
* `configurations` - A list of scaling rules. Each element contains the following attributes:
  * `id` - ID of the scaling rule.
  * `scaling_group_id` - ID of the scaling group.
  * `name` - Name of the scaling configuration.
  * `image_id` - Image ID of the scaling configuration.
  * `instance_type` - Instance type of the scaling configuration.
  * `security_group_id` - Security group ID of the scaling configuration.
  * `internet_charge_type` - Internet charge type of the scaling configuration.
  * `internet_max_bandwidth_in` - Internet max bandwidth in of the scaling configuration.
  * `internet_max_bandwidth_out` - Internet max bandwidth of the scaling configuration.
  * `credit_specification` - Performance mode of the t5 burstable instance.
  * `system_disk_category` - System disk category of the scaling configuration.
  * `system_disk_size` - System disk size of the scaling configuration.
  * `system_disk_performance_level` - The performance level of the ESSD used as the system disk.
  * `data_disks` - Data disks of the scaling configuration.
    * `size` - Size of data disk.
    * `category` - Category of data disk.
    * `snapshot_id` - Size of data disk.
    * `device` - Device attribute of data disk.
    * `delete_with_instance` - Delete_with_instance attribute of data disk.
    * `performance_level` - The performance level of the ESSD used as data disk.
  * `lifecycle_state` - Lifecycle state of the scaling configuration.
  * `creation_time` - Creation time of the scaling configuration.
  * `instance_name` - (Optional,Available since v1.143.0) InstanceName of an ECS instance.
  * `host_name` - (Optional,Available since v1.143.0) Hostname of an ECS instance.
  * `spot_strategy` - (Optional, Available since v1.151.0) The spot strategy for a Pay-As-You-Go instance.
  * `spot_price_limit` - (Optional, Available since v1.151.0) The maximum price hourly for instance types.
    * `instance_type` - Resource type of an ECS instance.
    * `price_limit` - Price limit hourly of instance type.
  * `instance_pattern_info` - (Optional, Available since v1.240.0) intelligent configuration mode. In this mode, you only need to specify the number of vCPUs, memory size, instance family, and maximum price. The system selects an instance type that is provided at the lowest price based on your configurations to create ECS instances. This mode is available only for scaling groups that reside in virtual private clouds (VPCs). This mode helps reduce the failures of scale-out activities caused by insufficient inventory of instance types. 
    * `cores` - The number of vCPUs that are specified for an instance type in instancePatternInfo.
    * `instance_family_level` - The instance family level in instancePatternInfo.
    * `max_price` - The maximum hourly price for a pay-as-you-go instance or a preemptible instance in instancePatternInfo.
    * `memory` - The memory size that is specified for an instance type in instancePatternInfo.
    * `burstable_performance` - Specifies whether to include burstable instance types.  Valid values: Exclude, Include, Required.
    * `excluded_instance_types` - Instance type N that you want to exclude. You can use wildcard characters, such as an asterisk (*), to exclude an instance type or an instance family.
    * `architectures` -  Architecture N of instance type N. Valid values: X86, Heterogeneous, BareMetal, Arm, SuperComputeCluster.
