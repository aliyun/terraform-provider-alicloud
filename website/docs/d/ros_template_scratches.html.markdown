---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_template_scratches"
sidebar_current: "docs-alicloud-datasource-ros-template-scratches"
description: |-
  Provides a list of Ros Template Scratches to the user.
---

# alicloud\_ros\_template\_scratches

This data source provides the Ros Template Scratches of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.151.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ros_template_scratches" "ids" {
  ids = ["example_value"]
}
output "ros_template_scratch_id_1" {
  value = data.alicloud_ros_template_scratches.ids.scratches.0.id
}

data "alicloud_ros_template_scratches" "status" {
  status = "GENERATE_COMPLETE"
}
output "ros_template_scratch_id_2" {
  value = data.alicloud_ros_template_scratches.status.scratches.0.id
}

data "alicloud_ros_template_scratches" "templateScratchType" {
  template_scratch_type = "ResourceImport"
}
output "ros_template_scratch_id_3" {
  value = data.alicloud_ros_template_scratches.templateScratchType.scratches.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Template Scratch IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid Values: `GENERATE_IN_PROGRESS`, `GENERATE_COMPLETE` and `GENERATE_FAILED`.
* `template_scratch_type` - (Optional, ForceNew) The type of the template scratch. Valid Values: `ResourceImport`, `ArchitectureReplication`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `scratches` - A list of Ros Template Scratches. Each element contains the following attributes:
  * `create_time` - The creation time of the resource.
  * `description` - The description of the Template Scratch.
  * `id` - The ID of the Template Scratch.
  * `logical_id_strategy` - The Logical ID generation strategy of the Template Scratch.
  * `preference_parameters` - Priority parameter.
    * `parameter_key` - Priority parameter key.
    * `parameter_value` - Priority parameter value.
  * `source_tag` - The Source label list.
    * `resource_tags` - Source label.
    * `resource_type_filter` - Source resource type filter list.
  * `source_resource_group` - Source resource grouping.
    * `resource_type_filter` - Source resource type filter list.
    * `resource_group_id` - The ID of the Source Resource Group.
  * `source_resources` - Source resource.
    * `resource_id` - The ID of the Source Resource.
    * `resource_type` - The type of the Source resource.
  * `stacks` - A list of resource stacks associated with the resource scene.
    * `stack_id` - The ID of the Resource stack.
  * `status` - The status of the resource.
  * `template_scratch_id` - The ID of the Template Scratch.
  * `template_scratch_type` - The type of the Template Scratch.