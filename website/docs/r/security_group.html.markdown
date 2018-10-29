---
layout: "alicloud"
page_title: "Alicloud: alicloud_security_group"
sidebar_current: "docs-alicloud-resource-security-group"
description: |-
  Provides a Alicloud Security Group resource.
---

# alicloud\_security\_group

Provides a security group resource.

~> **NOTE:** `alicloud_security_group` is used to build and manage a security group, and `alicloud_security_group_rule` can define ingress or egress rules for it.

~> **NOTE:** From version 1.7.2, `alicloud_security_group` has supported to segregate different ECS instance in which the same security group.

## Example Usage

Basic Usage

```
resource "alicloud_security_group" "group" {
  name        = "terraform-test-group"
  description = "New security group"
}
```
Basic usage for vpc

```
resource "alicloud_security_group" "group" {
  name   = "new-group"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the security group. Defaults to null.
* `description` - (Optional, Forces new resource) The security group description. Defaults to null.
* `vpc_id` - (Optional, Forces new resource) The VPC ID.
* `inner_access` - (Optional) Whether to allow both machines to access each other on all ports in the same security group.
* `tags` - (Optional) A mapping of tags to assign to the resource.

Combining security group rules, the policy can define multiple application scenario. Default to true. It is valid from verison `1.7.2`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the security group
* `vpc_id` - The VPC ID.
* `name` - The name of the security group
* `description` - The description of the security group
* `inner_access` - Whether to allow inner network access.
* `tags` - The instance tags, use jsonencode(item) to display the value.

## Import

Security Group can be imported using the id, e.g.

```
$ terraform import alicloud_security_group.example sg-abc123456
```