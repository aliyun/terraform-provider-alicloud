---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_pipelines"
sidebar_current: "docs-alicloud-datasource-ecs-image-pipelines"
description: |-
  Provides a list of Ecs Image Pipelines to the user.
---

# alicloud\_ecs\_image\_pipelines

This data source provides the Ecs Image Pipelines of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.163.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_image_pipelines" "ids" {
  ids = ["example_value"]
}
output "ecs_image_pipeline_id_1" {
  value = data.alicloud_ecs_image_pipelines.ids.pipelines.0.id
}

data "alicloud_ecs_image_pipelines" "nameRegex" {
  name_regex = "^my-ImagePipeline"
}
output "ecs_image_pipeline_id_2" {
  value = data.alicloud_ecs_image_pipelines.nameRegex.pipelines.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Image Pipeline IDs.
* `name` - (Optional, ForceNew) The name of the image template.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Image Pipeline name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group to which the image template belongs.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Image Pipeline ids.
* `names` - A list of Image Pipeline names.
* `pipelines` - A list of Ecs Image Pipelines. Each element contains the following attributes:
  * `base_image` - The source image.
  * `base_image_type` - The type of the source image.
  * `build_content` - The content of the image template.
  * `creation_time` - The time when the image template was created.
  * `delete_instance_on_failure` - Indicates whether the intermediate instance was released when the image failed to be created.
  * `description` - The description of the image template.
  * `id` - The ID of the Image Pipeline.
  * `image_name` - The name prefix of the created image.
  * `image_pipeline_id` - The ID of the image template.
  * `instance_type` - The instance type of the intermediate instance.
  * `internet_max_bandwidth_out` - The size of the outbound public bandwidth for the intermediate instance. Unit: `Mbit/s`.
  * `name` - The name of the image template.
  * `resource_group_id` - The ID of the resource group to which the image template belongs.
  * `system_disk_size` - The system disk size of the intermediate instance. Unit: `GiB`.
  * `vswitch_id` - The vswitch id.
  * `add_account` - The IDs of Alibaba Cloud accounts to which the image was shared.
  * `to_region_id` - The IDs of regions to which to distribute the created image.
  * `tags` - A mapping of tags to assign to the resource.