---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_group"
sidebar_current: "docs-alicloud-resource-ess-scaling-group"
description: |-
  Provides a ESS scaling group resource.
---

# alicloud_ess_scaling_group

Provides a ESS scaling group resource which is a collection of ECS instances with the same application scenarios.

It defines the maximum and minimum numbers of ECS instances in the group, and their associated Server Load Balancer instances, RDS instances, and other attributes.

-> **NOTE:** You can launch an ESS scaling group for a VPC network via specifying parameter `vswitch_ids`.

For information about ess scaling rule, see [CreateScalingGroup](https://www.alibabacloud.com/help/en/auto-scaling/latest/createscalinggroup).

-> **NOTE:** Available since v1.39.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ess_scaling_group&exampleId=4e5fd0ed-2555-7bf7-862a-260b88252787575ca4c6&activeTab=example&spm=docs.r.ess_scaling_group.0.4e5fd0ed25&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

resource "alicloud_vswitch" "default2" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = "${var.name}-bar"
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = local.name
  default_cooldown   = 20
  vswitch_ids        = [alicloud_vswitch.default.id, alicloud_vswitch.default2.id]
  removal_policies   = ["OldestInstance", "NewestInstance"]
}
```

## Module Support

You can use to the existing [autoscaling module](https://registry.terraform.io/modules/terraform-alicloud-modules/autoscaling/alicloud) 
to create a scaling group, configuration and lifecycle hook one-click.

## Argument Reference

The following arguments are supported:

* `min_size` - (Required) Minimum number of ECS instances in the scaling group. Value range: [0, 2000].
  **NOTE:** From version 1.204.1, `min_size` can be set to `2000`.
* `max_size` - (Required) Maximum number of ECS instances in the scaling group. Value range: [0, 2000].
  **NOTE:** From version 1.204.1, `max_size` can be set to `2000`.
* `desired_capacity` - (Optional, Available since v1.76.0) Expected number of ECS instances in the scaling group. Value range: [min_size, max_size].
* `scaling_group_name` - (Optional) Name shown for the scaling group, which must contain 2-64 characters (English or Chinese), starting with numbers, English letters or Chinese characters, and can contain numbers, underscores `_`, hyphens `-`, and decimal points `.`. If this parameter is not specified, the default value is ScalingGroupId.
* `default_cooldown` - (Optional) Default cool-down time (in seconds) of the scaling group. Value range: [0, 86400]. The default value is 300s.
* `vswitch_id` - (Deprecated) It has been deprecated from version 1.7.1 and new field 'vswitch_ids' replaces it.
* `vswitch_ids` - (Optional) List of virtual switch IDs in which the ecs instances to be launched.
* `removal_policies` - (Optional) RemovalPolicy is used to select the ECS instances you want to remove from the scaling group when multiple candidates for removal exist. Optional values:
    - OldestInstance: removes the ECS instance that is added to the scaling group at the earliest point in time.
    - NewestInstance: removes the ECS instance that is added to the scaling group at the latest point in time.
    - OldestScalingConfiguration: removes the ECS instance that is created based on the earliest scaling configuration.
    - Default values: Default value of RemovalPolicy.1: OldestScalingConfiguration. Default value of RemovalPolicy.2: OldestInstance.
* `db_instance_ids` - (Optional) If an RDS instance is specified in the scaling group, the scaling group automatically attaches the Intranet IP addresses of its ECS instances to the RDS access whitelist.
    - The specified RDS instance must be in running status.
    - The specified RDS instance’s whitelist must have room for more IP addresses.
* `loadbalancer_ids` - (Optional) If a Server Load Balancer instance is specified in the scaling group, the scaling group automatically attaches its ECS instances to the Server Load Balancer instance.
    - The Server Load Balancer instance must be enabled.
    - At least one listener must be configured for each Server Load Balancer and it HealthCheck must be on. Otherwise, creation will fail (it may be useful to add a `depends_on` argument
      targeting your `alicloud_slb_listener` in order to make sure the listener with its HealthCheck configuration is ready before creating your scaling group).
    - The Server Load Balancer instance attached with VPC-type ECS instances cannot be attached to the scaling group.
    - The default weight of an ECS instance attached to the Server Load Balancer instance is 50.
* `multi_az_policy` - (Optional, ForceNew) Multi-AZ scaling group ECS instance expansion and contraction strategy. PRIORITY, COMPOSABLE, BALANCE or COST_OPTIMIZED(Available since v1.54.0).
* `az_balance` - (Optional, Available since v1.225.1) Specifies whether to evenly distribute instances in the scaling group across multiple zones. This parameter takes effect only if you set MultiAZPolicy to COMPOSABLE.
* `spot_allocation_strategy` - (Optional, Available since v1.225.1) The allocation policy of preemptible instances. You can use this parameter to individually specify the allocation policy for preemptible instances. This parameter takes effect only if you set MultiAZPolicy to COMPOSABLE.
* `allocation_strategy` - (Optional, Available since v1.225.1) The allocation policy of instances. Auto Scaling selects instance types based on the allocation policy to create instances. The policy can be applied to pay-as-you-go instances and preemptible instances. This parameter takes effect only if you set MultiAZPolicy to COMPOSABLE.
* `on_demand_base_capacity` - (Optional, Available since v1.54.0) The minimum amount of the Auto Scaling group's capacity that must be fulfilled by On-Demand Instances. This base portion is provisioned first as your group scales.
* `compensate_with_on_demand` - (Optional, Available since v1.245.0) Specifies whether to automatically create pay-as-you-go instances to meet the requirement on the number of ECS instances when the expected capacity of preemptible instances cannot be provided due to reasons such as cost-related issues and insufficient resources. This parameter is supported only if you set 'multi_az_policy' to COST_OPTIMIZED. Valid values: true, false.
* `capacity_options_on_demand_base_capacity` - (Optional, Available since v1.245.0) The minimum number of pay-as-you-go instances that must be contained in the scaling group. When the actual number of pay-as-you-go instances in the scaling group drops below the value of this parameter, Auto Scaling preferentially creates pay-as-you-go instances. Valid values: 0 to 1000. If you set 'multi_az_policy' to COMPOSABLE, the default value of this parameter is 0.
* `capacity_options_on_demand_percentage_above_base_capacity` - (Optional, Available since v1.245.0) The percentage of pay-as-you-go instances in the excess instances when the minimum number of pay-as-you-go instances is reached. 'on_demand_base_capacity' specifies the minimum number of pay-as-you-go instances that must be contained in the scaling group. Valid values: 0 to 100. If you set 'multi_az_policy' to COMPOSABLE, the default value of this parameter is 100.
* `capacity_options_compensate_with_on_demand` - (Optional, Available since v1.245.0) Specifies whether to automatically create pay-as-you-go instances to meet the requirement on the number of ECS instances when the expected capacity of preemptible instances cannot be provided due to reasons such as cost-related issues and insufficient resources. This parameter is supported only if you set 'multi_az_policy' to COST_OPTIMIZED. Valid values: true, false.
* `capacity_options_spot_auto_replace_on_demand` - (Optional, Available since v1.245.0) Specifies whether to replace pay-as-you-go instances with preemptible instances. If you specify 'compensate_with_on_demand', it may result in a higher percentage of pay-as-you-go instances compared to the value of 'on_demand_percentage_above_base_capacity'. If you specify this parameter, Auto Scaling preferentially deploys preemptible instances to replace the surplus pay-as-you-go instances when preemptible instance types are available. If you specify 'compensate_with_on_demand', Auto Scaling creates pay-as-you-go instances when preemptible instance types are insufficient. To avoid retaining these pay-as-you-go instances for extended periods, Auto Scaling attempts to replace them with preemptible instances when sufficient preemptible instance types become available. Valid values: true, false.
* `capacity_options_price_comparison_mode` - (Optional, Available since v1.249.0) The price comparison mode. Valid values: PricePerUnit,PricePerVCpu. Default value: PricePerUnit.
* `on_demand_percentage_above_base_capacity` - (Optional, Available since v1.54.0) Controls the percentages of On-Demand Instances and Spot Instances for your additional capacity beyond OnDemandBaseCapacity.  
* `spot_instance_pools` - (Optional, Available since v1.54.0) The number of Spot pools to use to allocate your Spot capacity. The Spot pools is composed of instance types of lowest price.
* `spot_instance_remedy` - (Optional, Available since v1.54.0) Whether to replace spot instances with newly created spot/onDemand instance when receive a spot recycling message.
* `group_deletion_protection` - (Optional, Available since v1.102.0) Specifies whether the scaling group deletion protection is enabled. `true` or `false`, Default value: `false`.            
* `launch_template_id` - (Optional, Available since v1.141.0) Instance launch template ID, scaling group obtains launch configuration from instance launch template, see [Launch Template](https://www.alibabacloud.com/help/doc-detail/73916.html). Creating scaling group from launch template enable group automatically.
* `launch_template_version` - (Optional, Available since v1.159.0) The version number of the launch template. Valid values are the version number, `Latest`, or `Default`, Default value: `Default`.
* `group_type` - (Optional, ForceNew, Available since v1.164.0) Resource type within scaling group. Optional values: ECS, ECI. Default to ECS.
* `health_check_type` - (Optional, Available since v1.193.0) Resource type within scaling group. Optional values: ECS, ECI, NONE, LOAD_BALANCER. Default to ECS.
* `health_check_types` - (Optional, Available since v1.228.0) The health check modes of the scaling group. Valid values: ECS, NONE, LOAD_BALANCER.
* `instance_id` - (Optional, ForceNew, Available since v1.228.0) The ID of the instance from which Auto Scaling obtains the required configuration information and uses the information to automatically create a scaling configuration.
* `container_group_id` - (Optional, ForceNew, Available since v1.233.0) The ID of the elastic container instance.
* `scaling_policy` - (Optional, Available since v1.227.0) The reclaim mode of the scaling group. Optional values: recycle, release, forceRecycle, forceRelease. 
* `max_instance_lifetime` - (Optional, Available since v1.227.0) The maximum life span of an instance in the scaling group. Unit: seconds.
* `stop_instance_timeout` - (Optional, Available since v1.238.0) The period of time required by the ECS instance to enter the Stopped state. Unit: seconds. Valid values: 30 to 240.
* `tags` - (Optional, Available since v1.160.0) A mapping of tags to assign to the resource.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `protected_instances` - (Optional, Available since v1.182.0) Set or unset instances within group into protected status.
* `launch_template_override` - (Optional, Available since v1.216.0) The details of the instance types that are specified by using the Extend Instance Type of Launch Template feature.  See [`launch_template_override`](#launch_template_override) below for details.
* `resource_group_id` - (Optional, Available since v1.224.0) The ID of the resource group to which you want to add the scaling group.
* `alb_server_group` - (Optional, Available since v1.224.0) If a Serve ALB instance is specified in the scaling group, the scaling group automatically attaches its ECS instances to the Server ALB instance.  See [`alb_server_group`](#alb_server_group) below for details.

### `alb_server_group`

The AlbServerGroup mapping supports the following:

* `alb_server_group_id` - (Optional) The ID of ALB server group.
* `weight` - (Optional) The weight of the ECS instance as a backend server after Auto Scaling adds the ECS instance to ALB server group.
* `port` - (Optional) The port number used by an ECS instance after Auto Scaling adds the ECS instance to ALB server group.


### `launch_template_override`

The launchTemplateOverride mapping supports the following:

* `weighted_capacity` - (Optional) The weight of the instance type in launchTemplateOverride.
* `instance_type` - (Optional) The instance type in launchTemplateOverride.
* `spot_price_limit` - (Optional) The maximum bid price of instance type in launchTemplateOverride.


-> **NOTE:** When detach loadbalancers, instances in group will be remove from loadbalancer's `Default Server Group`; On the contrary, When attach loadbalancers, instances in group will be added to loadbalancer's `Default Server Group`.

-> **NOTE:** When detach dbInstances, private ip of instances in group will be remove from dbInstance's `WhiteList`; On the contrary, When attach dbInstances, private ip of instances in group will be added to dbInstance's `WhiteList`.

-> **NOTE:** `on_demand_base_capacity`,`on_demand_percentage_above_base_capacity`,`spot_instance_pools`,`spot_instance_remedy` are valid only if `multi_az_policy` is 'COST_OPTIMIZED'.


## Attributes Reference

The following attributes are exported:

* `id` - The scaling group ID.

## Import

ESS scaling group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_scaling_group.example asg-abc123456
```
