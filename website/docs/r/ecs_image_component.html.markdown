---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_component"
description: |-
  Provides a Alicloud Ecs Image Component resource.
---

# alicloud_ecs_image_component

Provides a Ecs Image Component resource. 

For information about Ecs Image Component and how to use it, see [What is Image Component](https://www.alibabacloud.com/help/en/doc-detail/200424.htm).

-> **NOTE:** Available since v1.159.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_ecs_image_component" "example" {
  component_type       = "Build"
  content              = "RUN yum update -y"
  description          = "example_value"
  image_component_name = "example_value"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  system_type          = "Linux"
  tags = {
    Created = "TF"
  }
}
```

## Argument Reference

The following arguments are supported:
* `component_type` - (Optional, ForceNew, Computed) The component type. Currently, only mirror build components are supported. Value: Build.  Default value: Build.
* `content` - (Required, ForceNew) Component content.
* `description` - (Optional, ForceNew) Describe the information.
* `image_component_name` - (Optional, ForceNew, Computed) The component name. The name must be 2 to 128 characters in length and must start with an uppercase letter or a Chinese character. It cannot start with http:// or https. Can contain Chinese, English, numbers, half-length colons (:), underscores (_), half-length periods (.), or dashes (-).  Note: If Name is not set, the return value of ImageComponentId is used by default.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `system_type` - (Optional, ForceNew, Computed) The operating system supported by the component. Currently, only Linux systems are supported. Value: Linux.  Default value: Linux.
* `tags` - (Optional, Map) List of label key-value pairs.

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

Ecs Image Component can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_image_component.example <id>
```