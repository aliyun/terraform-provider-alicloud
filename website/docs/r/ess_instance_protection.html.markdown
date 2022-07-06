---
subcategory: "Auto Scaling(ESS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_instance_protection"
sidebar_current: "docs-alicloud-resource-ess-instance-protection"
description: |-
  Provides a ESS Instance Protection resource on  ECS instances of scaling group setting protection.
---

# alicloud\_ess\_protection

protect several ECS instances to a specified scaling group.

For information about a ECS instance is protected in scaling group, see [SetInstancesProtection](https://www.alibabacloud.com/help/en/auto-scaling/latest/setinstancesprotection).

-> **NOTE:** ECS instances can be protected only when the scaling group is active and it has no scaling activity in progress.
-> **NOTE:** Resource `alicloud_ess_instance_protection` is available in 1.178.0+.

## Example Usage

```
variable "name" {
  default = "essattachmentconfig"
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
  name              = var.name
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

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 2
  scaling_group_name = var.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  image_id          = data.alicloud_images.default.images[0].id
  instance_type     = data.alicloud_instance_types.default.instance_types[0].id
  security_group_id = alicloud_security_group.default.id
  force_delete      = true
  active            = true
  enable            = true
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  count                      = 2
  security_groups            = [alicloud_security_group.default.id]
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = var.name
}
resource "alicloud_ess_instance_protection" "default" {
  scaling_group_id = alicloud_ess_scaling_group.default.id
  instance_id     = alicloud_instance.default[0].id
  depends_on = ["alicloud_ess_attachment.default"]
}

resource "alicloud_ess_attachment" "default" {
  scaling_group_id = alicloud_ess_scaling_group.default.id
  instance_ids     = [alicloud_instance.default[0].id]
  force            = true
}

```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) ID of the scaling group of a scaling configuration.
* `instance_id` - (Required) ID of the ECS instance to be protected to the scaling group. 

-> **NOTE:** Restrictions on protecting ECS instances:

   - The protected ECS instances must in the scaling group.

## Attributes Reference

The following attributes are exported:

* `id` - (Required, ForceNew) The ESS instance protection resource ID，in the follwing format: scaling_group_id:instance_id.
* `instance_id` - (Required)ID of "Protected" ECS instance.

## Import

ESS instance protection can be imported using the id or scaling group id, e.g.

```
$ terraform import alicloud_ess_attachment.example asg-abc123456
```

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the resource.
* `delete` - (Defaults to 5 mins) Used when delete the resource.