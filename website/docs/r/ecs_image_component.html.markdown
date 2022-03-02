---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_component"
sidebar_current: "docs-alicloud-resource-ecs-image-component"
description: |-
  Provides a Alicloud ECS Image Component resource.
---

# alicloud\_ecs\_image\_component

Provides a ECS Image Component resource.

For information about ECS Image Component and how to use it, see [What is Image Component](https://www.alibabacloud.com/help/en/doc-detail/200424.htm).

-> **NOTE:** Available in v1.159.0+.

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

* `component_type` - (Optional, ForceNew, Computed) The type of the image component. Only image building components are supported. Valid values: `Build`.
* `content` - (Required, ForceNew) The content of the image component. The content can consist of up to 127 commands.
* `description` - (Optional, ForceNew) The description of the image component. The description must be `2` to `256` characters in length and cannot start with `http://` or `https://`.
* `image_component_name` - (Optional, Computed, ForceNew) The name of the image component. The name must be `2` to `128` characters in length. It must start with a letter and cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (_), periods (.), and hyphens (-).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group to which to assign the image component.
* `system_type` - (Optional, ForceNew, Computed) The operating system type supported by the image component. Only Linux is supported. Valid values: `Linux`.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Image Component.

## Import

ECS Image Component can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_image_component.example <id>
```