---
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_vsc"
description: |-
  Provides an ECS Vsc resource.
---

# alicloud_ecs_vsc

-> **NOTE:** Available since v1.295.0.

Provides an ECS Vsc resource.

## Example Usage

```terraform
provider "alicloud" {
  region = "cn-wulanchabu"
}

resource "alicloud_ecs_instance" "default" {
  instance_type        = "ecs.ebmg9i-vsc.48xlarge"
  image_id             = "aliyun_3_x64_22-04_x64_alibase"
  security_groups      = ["sg-xxx"]
  instance_name        = "tf-vsc-example-instance"
  vswitch_id           = "vsw-xxx"
  instance_charge_type = "PostPaid"
}

resource "alicloud_ecs_vsc" "default" {
  instance_id = alicloud_ecs_instance.default.id
  vsc_name    = "tf-vsc-example"
  description = "tf-vsc-example"
  tags = {
    Created = "tf"
    For     = "vsc-example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` (Optional) - The description of the Vsc. Changing this creates a new resource.
* `instance_id` (Required) - The ID of the ECS instance to which the Vsc belongs. Changing this creates a new resource.
* `resource_group_id` (Optional) - The ID of the resource group to which the Vsc belongs. Changing this creates a new resource.
* `vsc_name` (Optional) - The name of the Vsc. Changing this creates a new resource.
* `vsc_type` (Optional) - The type of the Vsc. Changing this creates a new resource.
* `tags` (Optional) - A mapping of tags to assign to the Vsc. See [`tags`](#tags) below.

### `tags`

The `tags` block supports:

* `key` (Required) - The tag key.
* `value` (Required) - The tag value.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Vsc.
* `status` - The status of the Vsc.

## Import

ECS Vsc can be imported using the id, e.g.

```shell
terraform import alicloud_ecs_vsc.example <id>
```
