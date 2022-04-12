---
subcategory: "Auto Scaling(ESS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_group"
sidebar_current: "docs-alicloud-resource-ess-scaling-group"
description: |-
  Provides a ESS scaling group resource.
---

# alicloud\_ess\_scaling\_group

Provides a ESS scaling group resource which is a collection of ECS instances with the same application scenarios.

It defines the maximum and minimum numbers of ECS instances in the group, and their associated Server Load Balancer instances, RDS instances, and other attributes.

-> **NOTE:** You can launch an ESS scaling group for a VPC network via specifying parameter `vswitch_ids`.

## Example Usage

```
variable "name" {
  default = "essscalinggroupconfig"
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
  vpc_name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
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
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.1.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = "${var.name}-bar"
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = var.name
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

* `min_size` - (Required) Minimum number of ECS instances in the scaling group. Value range: [0, 1000].
* `max_size` - (Required) Maximum number of ECS instances in the scaling group. Value range: [0, 1000].
* `desired_capacity` - (Optional,Available in 1.76.0+) Expected number of ECS instances in the scaling group. Value range: [min_size, max_size].
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
* `multi_az_policy` - (Optional, ForceNew) Multi-AZ scaling group ECS instance expansion and contraction strategy. PRIORITY, BALANCE or COST_OPTIMIZED(Available in 1.54.0+).
* `on_demand_base_capacity` - (Optional, Available in v1.54.0+) The minimum amount of the Auto Scaling group's capacity that must be fulfilled by On-Demand Instances. This base portion is provisioned first as your group scales.
* `on_demand_percentage_above_base_capacity` - (Optional, Available in v1.54.0+) Controls the percentages of On-Demand Instances and Spot Instances for your additional capacity beyond OnDemandBaseCapacity.  
* `spot_instance_pools` - (Optional, Available in v1.54.0+) The number of Spot pools to use to allocate your Spot capacity. The Spot pools is composed of instance types of lowest price.
* `spot_instance_remedy` - (Optional, Available in v1.54.0+) Whether to replace spot instances with newly created spot/onDemand instance when receive a spot recycling message.
* `group_deletion_protection` - (Optional, Available in v1.102.0+) Specifies whether the scaling group deletion protection is enabled. `true` or `false`, Default value: `false`.            
* `launch_template_id` - (Optional, Available in v1.141.0+) Instance launch template ID, scaling group obtains launch configuration from instance launch template, see [Launch Template](https://www.alibabacloud.com/help/doc-detail/73916.html). Creating scaling group from launch template enable group automatically.
* `launch_template_version` - (Optional, Available in v1.159.0+) The version number of the launch template. Valid values are the version number, `Latest`, or `Default`, Default value: `Default`.
* `group_type` - (Optional, Available in v1.164.0+) Resource type within scaling group. Optional values: ECS, ECI. Default to ECS.
* `tags` - (Optional, Available in v1.160.0+) A mapping of tags to assign to the resource.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

-> **NOTE:** When detach loadbalancers, instances in group will be remove from loadbalancer's `Default Server Group`; On the contrary, When attach loadbalancers, instances in group will be added to loadbalancer's `Default Server Group`.

-> **NOTE:** When detach dbInstances, private ip of instances in group will be remove from dbInstance's `WhiteList`; On the contrary, When attach dbInstances, private ip of instances in group will be added to dbInstance's `WhiteList`.

-> **NOTE:** `on_demand_base_capacity`,`on_demand_percentage_above_base_capacity`,`spot_instance_pools`,`spot_instance_remedy` are valid only if `multi_az_policy` is 'COST_OPTIMIZED'.


## Attributes Reference

The following attributes are exported:

* `id` - The scaling group ID.
* `min_size` - The minimum number of ECS instances.
* `max_size` - The maximum number of ECS instances.
* `scaling_group_name` - The name of the scaling group.
* `default_cooldown` - The default cool-down of the scaling group.
* `removal_policies` - The removal policy used to select the ECS instance to remove from the scaling group.
* `db_instance_ids` - The db instances id which the ECS instance attached to.
* `loadbalancer_ids` - The slb instances id which the ECS instance attached to.
* `vswitch_ids` - The vswitches id in which the ECS instance launched.
* `launch_template_id` - The instance launch template ID.
* `launch_template_version` - The version number of the launch template.

## Import

ESS scaling group can be imported using the id, e.g.

```
$ terraform import alicloud_ess_scaling_group.example asg-abc123456
```
