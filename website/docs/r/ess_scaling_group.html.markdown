---
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_14.*_64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alicloud_security_group.default.id}"
  cidr_ip           = "172.16.0.0/24"
}

resource "alicloud_vswitch" "default2" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}-bar"
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "${var.name}"
  default_cooldown   = 20
  vswitch_ids        = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
  removal_policies   = ["OldestInstance", "NewestInstance"]
}
```

## Argument Reference

The following arguments are supported:

* `min_size` - (Required) Minimum number of ECS instances in the scaling group. Value range: [0, 1000].
* `max_size` - (Required) Maximum number of ECS instances in the scaling group. Value range: [0, 1000].
* `scaling_group_name` - (Optional) Name shown for the scaling group, which must contain 2-40 characters (English or Chinese). If this parameter is not specified, the default value is ScalingGroupId.
* `default_cooldown` - (Optional) Default cool-down time (in seconds) of the scaling group. Value range: [0, 86400]. The default value is 300s.
* `vswitch_id` - (Deprecated) It has been deprecated from version 1.7.1 and new field 'vswitch_ids' replaces it.
* `vswitch_ids` - (Optional) List of virtual switch IDs in which the ecs instances to be launched.
* `removal_policies` - (Optional) RemovalPolicy is used to select the ECS instances you want to remove from the scaling group when multiple candidates for removal exist. Optional values:
    - OldestInstance: removes the first ECS instance attached to the scaling group.
    - NewestInstance: removes the first ECS instance attached to the scaling group.
    - OldestScalingConfiguration: removes the ECS instance with the oldest scaling configuration.
    - Default values: OldestScalingConfiguration and OldestInstance. You can enter up to two removal policies.
* `db_instance_ids` - (Optional) If an RDS instance is specified in the scaling group, the scaling group automatically attaches the Intranet IP addresses of its ECS instances to the RDS access whitelist.
    - The specified RDS instance must be in running status.
    - The specified RDS instance’s whitelist must have room for more IP addresses.
* `loadbalancer_ids` - (Optional) If a Server Load Balancer instance is specified in the scaling group, the scaling group automatically attaches its ECS instances to the Server Load Balancer instance.
    - The Server Load Balancer instance must be enabled.
    - At least one listener must be configured for each Server Load Balancer and it HealthCheck must be on. Otherwise, creation will fail (it may be useful to add a `depends_on` argument
      targeting your `alicloud_slb_listener` in order to make sure the listener with its HealthCheck configuration is ready before creating your scaling group).
    - The Server Load Balancer instance attached with VPC-type ECS instances cannot be attached to the scaling group.
    - The default weight of an ECS instance attached to the Server Load Balancer instance is 50.
* `vserver_groups` - (Optional, Available in 1.52.3+) If a LoadBalancer VServer Group is specified in the scaling group, the scaling group automatically attaches its ECS instances to VServer Group. See [Block vserver_group](#block-vserver_group) below for details.
* `multi_az_policy` - (Optional, ForceNew) Multi-AZ scaling group ECS instance expansion and contraction strategy. PRIORITY or BALANCE.

-> **NOTE:** When detach loadbalancers, instances in group will be remove from loadbalancer's `Default Server Group`; On the contrary, When attach loadbalancers, instances in group will be added to loadbalancer's `Default Server Group`.

-> **NOTE:** When detach dbInstances, private ip of instances in group will be remove from dbInstance's `WhiteList`; On the contrary, When attach dbInstances, private ip of instances in group will be added to dbInstance's `WhiteList`.

-> **NOTE:** A VServer Group is exclusively defined by loadbalancer_id, vserver_group_id and port, therefore, if you want to update the weight attribute, you need first detach VServer Group and then, attch it with newly changed weigth.

## Block vserver_group

the vserver_group supports the following:

* `loadbalancer_id` - (Required) Loadbalancer server ID of VServer Group.
* `vserver_attributes` - (Required) A Set of VServer Group attributes. See [Block vserver_groups](#block-vserver_attribute) below for details.

## Block vserver_attribute

* `vserver_group_id` - (Required) ID of VServer Group.
* `port` - (Required) - The port will be used for VServer Group backend server.
* `weight` - (Required) The weight of an ECS instance attached to the VServer Group.

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

## Import

ESS scaling group can be imported using the id, e.g.

```
$ terraform import alicloud_ess_scaling_group.example asg-abc123456
```
