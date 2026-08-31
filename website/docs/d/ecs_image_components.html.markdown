---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_components"
sidebar_current: "docs-alicloud-datasource-ecs-image-components"
description: |-
  Provides a list of Ecs Image Components to the user.
---

# alicloud\_ecs\_image\_components

This data source provides the Ecs Image Components of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.159.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_image_components" "ids" {
  ids = ["example_id"]
}
output "ecs_image_component_id_1" {
  value = data.alicloud_ecs_image_components.ids.components.0.id
}

data "alicloud_ecs_image_components" "nameRegex" {
  name_regex = "^my-ImageComponent"
}
output "ecs_image_component_id_2" {
  value = data.alicloud_ecs_image_components.nameRegex.components.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Image Component IDs.
* `image_component_name` - (Optional, ForceNew) The name of the component.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Image Component name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `owner` - (Optional, ForceNew) Mirror component type. Valid values: `SELF` or `ALIYUN`. Possible values:
  - SELF: The custom image component you created.
  - ALIYUN: System components provided by Alibaba Cloud.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Image Component names.
* `components` - A list of Ecs Image Components. Each element contains the following attributes:
	* `component_type` - The type of the image component.
	* `content` - The content of the image component.
	* `create_time` - The time when the image component was created.
	* `description` - The description of the image component.
	* `id` - The ID of the Image Component.
	* `image_component_id` - The ID of the image component.
	* `image_component_name` - The name of the image component.
	* `owner` - The type of the image component.
	* `resource_group_id` - The ID of the resource group.
	* `system_type` - The operating system type supported by the image component.
	* `tags` - List of label key-value pairs.
		* `tag_key` - The key of the tag.
		* `tag_value` - The value of the tag.