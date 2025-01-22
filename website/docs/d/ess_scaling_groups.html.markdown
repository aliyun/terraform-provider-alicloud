---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_groups"
sidebar_current: "docs-alicloud_ess_scaling_groups"
description: |-
    Provides a list of scaling groups available to the user.
---

# alicloud_ess_scaling_groups

This data source provides available scaling group resources. 

-> **NOTE:** Available since v1.39.0

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

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = local.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
}

data "alicloud_ess_scaling_groups" "scalinggroups_ds" {
  ids        = [alicloud_ess_scaling_group.default.id]
  name_regex = local.name
}

output "first_scaling_group" {
  value = "${data.alicloud_ess_scaling_groups.scalinggroups_ds.groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter resulting scaling groups by name.
* `ids` - (Optional, ForceNew) A list of scaling group IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of scaling group ids.
* `names` - A list of scaling group names.
* `groups` - A list of scaling groups. Each element contains the following attributes:
  * `id` - ID of the scaling group.
  * `name` - Name of the scaling group.
  * `active_scaling_configuration` -Active scaling configuration for scaling group.
  * `launch_template_id` - Active launch template ID for scaling group.
  * `launch_template_version` - Version of active launch template.
  * `region_id` - Region ID the scaling group belongs to.
  * `min_size` - The minimum number of ECS instances.
  * `max_size` - The maximum number of ECS instances.
  * `cooldown_time` - Default cooldown time of scaling group.
  * `removal_policies` - Removal policy used to select the ECS instance to remove from the scaling group.
  * `load_balancer_ids` - Slb instances id which the ECS instance attached to.
  * `db_instance_ids` - Db instances id which the ECS instance attached to.
  * `vswitch_ids` - Vswitches id in which the ECS instance launched.
  * `lifecycle_state` - Lifecycle state of scaling group.
  * `vpc_id` - The ID of the VPC to which the scaling group belongs.
  * `vswitch_id` - The ID of the vSwitch to which the scaling group belongs.
  * `health_check_type` - The health check method of the scaling group.
  * `suspended_processes` - The Process in suspension.
  * `group_deletion_protection` - Whether the scaling group deletion protection is enabled.
  * `modification_time` - The modification time.
  * `total_capacity` - Number of instances in scaling group.
  * `total_instance_count` - The number of all ECS instances in the scaling group.
  * `active_capacity` - Number of active instances in scaling group.
  * `init_capacity` - (Available since v1.242.0) The number of instances that are in the Initialized state and ready to be scaled out in the scaling group.
  * `pending_wait_capacity` - (Available since v1.242.0) The number of ECS instances that are in the Pending Add state in the scaling group.
  * `removing_wait_capacity` - (Available since v1.242.0) The number of ECS instances that are in the Pending Remove state in the scaling group.
  * `protected_capacity` - (Available since v1.242.0) The number of ECS instances that are in the Protected state in the scaling group.
  * `standby_capacity` - (Available since v1.242.0) The number of instances that are in the Standby state in the scaling group.
  * `spot_capacity` - (Available since v1.242.0) The number of preemptible instances in the scaling group.
  * `stopped_capacity` - (Available since v1.242.0) The number of instances that are in Economical Mode in the scaling group.
  * `pending_capacity` - (Available since v1.242.0) The number of ECS instances that are being added to the scaling group and still being configured.
  * `removing_capacity` - (Available since v1.242.0) The number of ECS instances that are being removed from the scaling group.
  * `system_suspended` - (Available since v1.242.0) Indicates whether Auto Scaling stops executing the scaling operation in the scaling group.
  * `monitor_group_id` - (Available since v1.242.0) The ID of the CloudMonitor application group that is associated with the scaling group.
  * `enable_desired_capacity` - (Available since v1.242.0) Indicates whether the Expected Number of Instances feature is enabled.
  * `creation_time` - Creation time of scaling group.
  * `tags` - A mapping of tags to assign to the resource.
  * `stop_instance_timeout` - (Available since v1.242.0) The period of time that is required by an ECS instance to enter the Stopped state during the scale-in process. Unit: seconds. 
  * `desired_capacity` - (Available since v1.242.0) The expected number of ECS instances in the scaling group. Auto Scaling automatically maintains the expected number of ECS instances that you specified.
  * `max_instance_lifetime` - (Available since v1.242.0) The maximum life span of each instance in the scaling group. Unit: seconds.
  * `multi_az_policy` - (Available since v1.242.0) The scaling policy of the multi-zone scaling group of the ECS type.
  * `group_type` - (Available since v1.242.0) The type of the instances in the scaling group. 
  * `resource_group_id` - (Available since v1.242.0) The ID of the resource group to which the scaling group that you want to query belongs.
  * `spot_instance_remedy` - (Available since v1.242.0) Indicates whether supplementation of preemptible instances is enabled. If this parameter is set to true, Auto Scaling creates an instance to replace a preemptible instance when Auto Scaling receives a system message indicating that the preemptible instance is to be reclaimed.
  * `spot_instance_pools` - (Available since v1.242.0) The number of instance types. Auto Scaling creates preemptible instances of multiple instance types that are provided at the lowest price.
  * `on_demand_percentage_above_base_capacity` - (Available since v1.242.0) The percentage of pay-as-you-go instances in the excess instances when the minimum number of pay-as-you-go instances is reached. OnDemandBaseCapacity specifies the minimum number of pay-as-you-go instances that must be contained in the scaling group. 
  * `on_demand_base_capacity` - (Available since v1.242.0) The lower limit of the number of pay-as-you-go instances in the scaling group.
  * `spot_allocation_strategy` - (Available since v1.242.0) The allocation policy of preemptible instances. This parameter indicates the method used by Auto Scaling to select instance types to create the required number of preemptible instances. This parameter takes effect only if you set multi_az_policy to COMPOSABLE.
  * `allocation_strategy` -  (Available since v1.242.0) The allocation policy of instances. Auto Scaling selects instance types based on the allocation policy to create instances. The allocation policy applies to pay-as-you-go and preemptible instances.
  * `az_balance` - (Available since v1.242.0) Indicates whether instances in the scaling group are evenly distributed across multiple zones.
  * `scaling_policy` - (Available since v1.242.0) The reclaim mode of the scaling group.