---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_group"
sidebar_current: "docs-alicloud-resource-security-group"
description: |-
  Provides a Alicloud Security Group resource.
---

# alicloud_security_group

Provides a Security Group resource.

For information about Security Group and how to use it, see [What is Security Group](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-createsecuritygroup).

-> **NOTE:** Available since v1.0.0.

-> **NOTE:** `alicloud_security_group` is used to build and manage a security group, and `alicloud_security_group_rule` can define ingress or egress rules for it.

-> **NOTE:** From version 1.7.2, `alicloud_security_group` has supported to segregate different ECS instance in which the same security group.

## Example Usage

Basic Usage

```terraform
resource "alicloud_security_group" "default" {
  name        = "terraform-example"
  description = "New security group"
}
```

Basic Usage for VPC

```terraform
resource "alicloud_vpc" "vpc" {
  vpc_name   = "terraform-example"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group" "group" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.vpc.id
}
```

## Module Support

You can use the existing [security-group module](https://registry.terraform.io/modules/alibaba/security-group/alicloud) 
to create a security group and add several rules one-click.

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Optional, ForceNew) The ID of the VPC.
* `security_group_type` - (Optional, ForceNew, Available since v1.58.0) The type of the security group. Valid values:
  - `normal`: basic security group.
  - `enterprise`: advanced security group For more information, see [Advanced security groups](https://www.alibabacloud.com/help/en/ecs/advanced-security-groups).
* `name` - (Optional) The name of the security group. Defaults to null.
* `description` - (Optional) The security group description. Defaults to null.
* `resource_group_id` - (Optional, Available since v1.58.0) The ID of the resource group to which the security group belongs. **NOTE:** From version 1.115.0, `resource_group_id` can be modified.
* `inner_access_policy` - (Optional, Available since v1.55.3) The internal access control policy of the security group. Valid values: `Accept`, `Drop`.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `inner_access` - (Deprecated since v1.55.3) Field `inner_access` has been deprecated from provider version 1.55.3. New field `inner_access_policy` instead.

Combining security group rules, the policy can define multiple application scenario. Default to true. It is valid from version `1.7.2`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Security Group.

## Timeouts

-> **NOTE:** Available since v1.214.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `delete` - (Defaults to 6 mins) Used when delete the Security Group.

## Import

Security Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_security_group.example sg-abc123456
```
