---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_pipeline_executions"
sidebar_current: "docs-alicloud-datasource-ecs-image-pipeline-executions"
description: |-
  Provides a list of ECS image pipeline executions.
---

# alicloud_ecs_image_pipeline_executions

The `alicloud_ecs_image_pipeline_executions` data source provides a collection of executions of ECS image pipelines in the current region. You can query the executions of one image pipeline or a single execution, and filter the results by status, tags and execution IDs.

For information about the API and its fields, see [DescribeImagePipelineExecutions](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-describeimagepipelineexecutions).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** At least one of `image_pipeline_id` and `execution_id` must be set. The data source reads every page returned by `DescribeImagePipelineExecutions` before applying `ids`, so the result is complete regardless of the number of executions.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_ecs_image_pipeline" "default" {
  base_image      = data.alicloud_images.default.ids.0
  base_image_type = "IMAGE"
  name            = "terraform-example"
}

data "alicloud_ecs_image_pipeline_executions" "default" {
  image_pipeline_id = alicloud_ecs_image_pipeline.default.id
}

output "execution_ids" {
  value = data.alicloud_ecs_image_pipeline_executions.default.ids
}
```

## Argument Reference

The following arguments are supported:

* `image_pipeline_id` - (Optional, String) The ID of the image pipeline whose executions are queried.
* `execution_id` - (Optional, String) The ID of a specific execution.
* `status` - (Optional, String) The execution status to filter by. Multiple statuses can be separated by commas. Valid values: `PREPARING`, `REPAIRING`, `BUILDING`, `TESTING`, `DISTRIBUTING`, `RELEASING`, `SUCCESS`, `PARTITION_SUCCESS`, `FAILED`, `TEST_FAILED`, and the cancellation statuses documented for `DescribeImagePipelineExecutions`.
* `tags` - (Optional, Map of String) A mapping of tags to filter by. A maximum of 20 entries is allowed.
* `ids` - (Optional, Computed, List of String) A list of execution IDs to filter by. The filter is applied on the provider side after all pages are read.
* `max_results` - (Optional, Int) The number of entries per page. Valid values: `1` to `500`. Defaults to `50`.
* `next_token` - (Optional, Computed, String) The pagination token to start reading from. The data source follows the remaining tokens itself; the value stored afterwards is the last token the server returned, and it is empty when every page has been read.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Computed, String) The ID of the data source, derived from the matched execution IDs.
* `total_count` - (Computed, Int) The number of execution records returned by this read after all configured filters, including `ids`.
* `executions` - (Computed, List of Object) A list of image pipeline executions. Each element contains the following attributes:
  * `execution_id` - (String) The ID of the execution.
  * `image_pipeline_id` - (String) The ID of the image pipeline.
  * `status` - (String) The status of the execution. Valid values are the same as the `status` filter.
  * `creation_time` - (String) The time when the execution was created, in UTC and formatted as `2026-11-05T10:00:00Z`.
  * `modified_time` - (String) The time when the execution was last modified, in UTC and formatted as `2026-11-05T10:00:00Z`.
  * `resource_group_id` - (String) The ID of the resource group.
  * `message` - (String) The message about the execution.
  * `image_id` - (String) The ID of the image created by the execution.
  * `error_code` - (String) The error code of the execution.
  * `error_info` - (String) The error details of the execution.
  * `tags` - (Map of String) A mapping of tags.
