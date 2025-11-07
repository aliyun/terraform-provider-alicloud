---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_attachment"
sidebar_current: "docs-alicloud-resource-ess-attachment"
description: |-
  Provides a ESS Attachment resource to attach or remove ECS instances.
---

# alicloud_ess_attachment

Attaches several ECS instances to a specified scaling group or remove them from it.

-> **NOTE:** ECS instances can be attached or remove only when the scaling group is active, and it has no scaling activity in progress.

-> **NOTE:** There are two types ECS instances in a scaling group: "AutoCreated" and "Attached". The total number of them can not larger than the scaling group "MaxSize".

-> **NOTE:** Available since v1.6.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ess_attachment&exampleId=fd33e833-99fb-73a7-2b4d-12381d2c3e857d343984&activeTab=example&spm=docs.r.ess_attachment.0.fd33e83399&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "ap-southeast-5"
}
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
  availability_zone    = data.alicloud_zones.default.zones[0].id
  cpu_core_count       = 2
  memory_size          = 8
  instance_type_family = "ecs.g9i"
}

data "alicloud_images" "default" {
  instance_type = data.alicloud_instance_types.default.instance_types[0].id
  most_recent   = true
  owners        = "system"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "10.4.0.0/16"
}

data "alicloud_vswitches" "default" {
  cidr_block = "10.4.0.0/24"
  vpc_id     = data.alicloud_vpcs.default.ids.0
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
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
  min_size           = 0
  max_size           = 2
  scaling_group_name = local.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [data.alicloud_vswitches.default.ids[0]]
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id     = alicloud_ess_scaling_group.default.id
  image_id             = data.alicloud_images.default.images[0].id
  instance_type        = data.alicloud_instance_types.default.instance_types[0].id
  security_group_id    = alicloud_security_group.default.id
  system_disk_category = "cloud_essd"
  force_delete         = true
  active               = true
  enable               = true
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  count                      = 2
  security_groups            = [alicloud_security_group.default.id]
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_essd"
  vswitch_id                 = data.alicloud_vswitches.default.ids[0]
  instance_name              = var.name
}

resource "alicloud_ess_attachment" "default" {
  scaling_group_id = alicloud_ess_scaling_group.default.id
  instance_ids     = [alicloud_instance.default[0].id, alicloud_instance.default[1].id]
  force            = true
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) ID of the scaling group of a scaling configuration.
* `instance_ids` - (Required) ID of the ECS instance to be attached to the scaling group. You can input up to 20 IDs.
* `force` - (Optional) Whether to remove forcibly "AutoCreated" ECS instances in order to release scaling group capacity "MaxSize" for attaching ECS instances. Default to false.
* `entrusted` - (Optional, Available since v1.220.0, ForceNew) Specifies whether the scaling group manages the lifecycles of the instances that are manually added to the scaling group.
* `lifecycle_hook` - (Optional, Available since v1.220.0, ForceNew) Specifies whether to trigger a lifecycle hook for the scaling group to which instances are being added.
* `load_balancer_weights` - (Optional, Available since v1.220.0, ForceNew) The weight of ECS instance N or elastic container instance N as a backend server of the associated Server Load Balancer (SLB) instance. Valid values of N: 1 to 20. Valid values of this parameter: 1 to 100.

-> **NOTE:** "AutoCreated" ECS instance will be deleted after it is removed from scaling group, but "Attached" will be not.

-> **NOTE:** Restrictions on attaching ECS instances:

   - The attached ECS instances and the scaling group must have the same region and network type(`Classic` or `VPC`).
   - The attached ECS instances and the instance with active scaling configurations must have the same instance type.
   - The attached ECS instances must in the running state.
   - The attached ECS instances has not been attached to other scaling groups.
   - The attached ECS instances supports Subscription and Pay-As-You-Go payment methods.

## Attributes Reference

The following attributes are exported:

* `id` - The ESS attachment resource ID.

## Import

ESS attachment can be imported using the id or scaling group id, e.g.

```shell
$ terraform import alicloud_ess_attachment.example asg-abc123456
```
