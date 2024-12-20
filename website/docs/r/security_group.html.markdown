---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_group"
description: |-
  Provides a Alicloud ECS Security Group resource.
---

# alicloud_security_group

Provides a ECS Security Group resource.



For information about ECS Security Group and how to use it, see [What is Security Group](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-createsecuritygroup).

-> **NOTE:** Available since v1.0.0.

-> **NOTE:** `alicloud_security_group` is used to build and manage a security group, and `alicloud_security_group_rule` can define ingress or egress rules for it.

-> **NOTE:** From version 1.7.2, `alicloud_security_group` has supported to segregate different ECS instance in which the same security group.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_security_group&exampleId=a779d5f6-6a79-09ec-6f53-24a3c18ea4b1edff5bfe&activeTab=example&spm=docs.r.security_group.0.a779d5f66a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_security_group" "default" {
  security_group_name = "terraform-example"
}
```

Basic Usage for VPC

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_security_group&exampleId=99236f01-ecf3-77e3-3c27-670b6276d29fcf85b305&activeTab=example&spm=docs.r.security_group.1.99236f01ec&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_security_group" "default" {
  security_group_name = "terraform-example"
  vpc_id              = alicloud_vpc.default.id
}
```

## Module Support

You can use the existing [security-group module](https://registry.terraform.io/modules/alibaba/security-group/alicloud) 
to create a security group and add several rules one-click.

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the security group. The description must be `2` to `256` characters in length. It cannot start with `http://` or `https://`.
* `inner_access_policy` - (Optional, Available since v1.55.3) The internal access control policy of the security group. Valid values:
  - `Accept`: The internal interconnectivity policy.
  - `Drop`: The internal isolation policy.
* `resource_group_id` - (Optional, Available since v1.58.0) The ID of the resource group to which the security group belongs. **NOTE:** From version 1.115.0, `resource_group_id` can be modified.
* `security_group_name` - (Optional, Available since v1.239.0) The name of the security group. The name must be `2` to `128` characters in length. The name must start with a letter and cannot start with `http://` or `https://`. The name can contain Unicode characters under the Decimal Number category and the categories whose names contain Letter. The name can also contain colons (:), underscores (\_), periods (.), and hyphens (-).
* `security_group_type` - (Optional, ForceNew, Available since v1.58.0) The type of the security group. Default value: `normal`. Valid values:
  - `normal`: Basic security group.
  - `enterprise`: Advanced security group For more information, see [Advanced security groups](https://www.alibabacloud.com/help/en/ecs/advanced-security-groups).
* `tags` - (Optional, Map) A mapping of tags to assign to the resource.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC in which you want to create the security group.
* `name` - (Optional, Deprecated since v1.239.0) Field `name` has been deprecated from provider version 1.239.0. New field `security_group_name` instead.
* `inner_access` - (Optional, Bool, Deprecated since v1.55.3) Field `inner_access` has been deprecated from provider version 1.55.3. New field `inner_access_policy` instead.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Security Group.
* `create_time` - (Available since v1.239.0) The time when the security group was created.

## Timeouts

-> **NOTE:** Available since v1.214.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Security Group.
* `delete` - (Defaults to 10 mins) Used when delete the Security Group.
* `update` - (Defaults to 5 mins) Used when update the Security Group.

## Import

ECS Security Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_security_group.example <id>
```
