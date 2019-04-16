---
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_attachment"
sidebar_current: "docs-alicloud-resource-ess-attachment"
description: |-
  Provides a ESS Attachment resource to attach or remove ECS instances.
---

# alicloud\_ess\_attachment

Attaches several ECS instances to a specified scaling group or remove them from it.

-> **NOTE:** ECS instances can be attached or remove only when the scaling group is active and it has no scaling activity in progress.

-> **NOTE:** There are two types ECS instances in a scaling group: "AutoCreated" and "Attached". The total number of them can not larger than the scaling group "MaxSize".

## Example Usage

```
resource "alicloud_instance" "instance" {
  # Other parameters...
}
resource "alicloud_ess_scaling_group" "scaling" {
  min_size           = 0
  max_size           = 2
  removal_policies   = ["OldestInstance", "NewestInstance"]

  # Other parameters...
}

resource "alicloud_ess_scaling_configuration" "config" {
  scaling_group_id  = "${alicloud_ess_scaling_group.scaling.id}"
  image_id          = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
  instance_type     = "ecs.n4.large"
  security_group_id = "${alicloud_security_group.classic.id}"
  active = true
  enable = true
}

resource "alicloud_ess_attachment" "att" {
  scaling_group_id = "${alicloud_ess_scaling_group.scaling.id}"
  instance_ids = ["${alicloud_instance.instance.*.id}"]
  force = true
}

```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) ID of the scaling group of a scaling configuration.
* `instance_ids` - (Required) ID of the ECS instance to be attached to the scaling group. You can input up to 20 IDs.
* `force` - (Optional) Whether to remove forcibly "AutoCreated" ECS instances in order to release scaling group capacity "MaxSize" for attaching ECS instances. Default to false.

-> **NOTE:** "AutoCreated" ECS instance will be deleted after it is removed from scaling group, but "Attached" will be not.

-> **NOTE:** Restrictions on attaching ECS instances:

   - The attached ECS instances and the scaling group must have the same region and network type(`Classic` or `VPC`).
   - The attached ECS instances and the instance with active scaling configurations must have the same instance type.
   - The attached ECS instances must in the running state.
   - The attached ECS instances has not been attached to other scaling groups.
   - The attached ECS instances supports Subscription and Pay-As-You-Go payment methods.

## Attributes Reference

The following attributes are exported:

* `id` - (Required, ForceNew) The ESS attachment resource ID.
* `instance_ids` - (Required)ID of list "Attached" ECS instance.
* `force` - Whether to delete "AutoCreated" ECS instances.

## Import

ESS attachment can be imported using the id or scaling group id, e.g.

```
$ terraform import alicloud_ess_attachment.example asg-abc123456
```