---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_role_attachment"
sidebar_current: "docs-alicloud-resource-ram-role-attachment"
description: |-
  Provides a RAM role attachment resource to bind role for several ECS instances.
---

# alicloud\_ram\_role\_attachment

Provides a RAM role attachment resource to bind role for several ECS instances.

## Example Usage

```
resource "alicloud_ram_role" "role" {
  name = "test_role"
  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
  ram_users = ["acs:ram::${your_account_id}:root", "acs:ram::${other_account_id}:user/username"]
  description = "this is a role test."
  force = true
}

resource "alicloud_instance" "instance" {
  instance_name = "test-keypair-${format(var.count_format, count.index+1)}"
  image_id = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
  instance_type = "ecs.n4.small"
  count = 2
  availability_zone = "${var.availability_zones}"
  ...
}

resource "alicloud_ram_role_attachment" "attach" {
  role_name = "${alicloud_ram_role.role.name}"
  instance_ids = ["${alicloud_instance.instance.*.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `role_name` - (Required, Forces new resource) The name of role used to bind. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-", "_", and must not begin with a hyphen.
* `instance_ids` - (Required, Forces new resource) The list of ECS instance's IDs.

## Attributes Reference

The following attributes are exported:

* `role_name` - The name of the role.
* `instance_ids` The list of ECS instance's IDs.