---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_component"
description: |-
  Provides a Alicloud ECS Image Component resource.
---

# alicloud_ecs_image_component

Provides a ECS Image Component resource. 

For information about ECS Image Component and how to use it, see [What is Image Component](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_resource_manager_resource_group" "ResourceGroup" {
  display_name        = "test"
  resource_group_name = var.name
}


resource "alicloud_ecs_image_component" "default" {
  image_component_name = var.name
  resource_group_id    = alicloud_resource_manager_resource_group.ResourceGroup.id
  content              = "RUN yum update -y"
}
```

## Argument Reference

The following arguments are supported:
* `component_type` - (Optional, ForceNew, Available since v1.159.0) Component type.
* `content` - (Optional, ForceNew, Available since v1.159.0) Component content.
* `description` - (Optional, ForceNew, Available since v1.159.0) Describe the information.
* `image_component_name` - (Optional, ForceNew, Available since v1.159.0) The name of the component.
* `resource_group_id` - (Optional, Computed, Available since v1.159.0) The ID of the resource group.
* `system_type` - (Optional, ForceNew, Available since v1.159.0) The operating system supported by the component.
* `tags` - (Optional, Map, Available since v1.159.0) List of label key-value pairs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Component creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Image Component.
* `delete` - (Defaults to 5 mins) Used when delete the Image Component.
* `update` - (Defaults to 5 mins) Used when update the Image Component.

## Import

ECS Image Component can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_image_component.example <id>
```