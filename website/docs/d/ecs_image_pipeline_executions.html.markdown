---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_pipeline_executions"
sidebar_current: "docs-alicloud-datasource-ecs-image-pipeline-executions"
description: |-
  Provides a list of Ecs Image Pipeline Executions to the user.
---

# alicloud\_ecs\_image\_pipeline\_executions

This data source provides the Ecs Image Pipeline Executions of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.166.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_image_pipeline_executions" "ids" {
  ids = ["example_value"]
}
output "ecs_image_pipeline_execution_id_1" {
  value = data.alicloud_ecs_image_pipeline_executions.ids.executions.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Image Pipeline Execution IDs.
* `image_pipeline_id` - (Required, ForceNew) The ID of the image template.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the mirror build task. Valid values: `BUILDING`, `CANCELLED`, `CANCELLING`, `DISTRIBUTING`, `FAILED`, `RELEASING`, `SUCCESS`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `executions` - A list of Ecs Image Pipeline Executions. Each element contains the following attributes:
	* `create_time` - The time when the image build task was created.
	* `id` - The ID of the Image Pipeline Execution.
	* `image_id` - The ID of the target image.
	* `image_pipeline_execution_id` - The first ID of the resource.
	* `image_pipeline_id` - The ID of the image template.
	* `message` - Execution result information.
	* `modified_time` - The time when the task was last updated.
	* `resource_group_id` - The ID of the enterprise resource group.
	* `status` - The status of the mirror build task.