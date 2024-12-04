---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_pipeline_execution"
description: |-
  Provides a Alicloud ECS Image Pipeline Execution resource.
---

# alicloud_ecs_image_pipeline_execution

Provides a ECS Image Pipeline Execution resource.

The mirror template performs the build mirror task.

For information about ECS Image Pipeline Execution and how to use it, see [What is Image Pipeline Execution](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "pipelineExecution-vpc" {
  description = "example-pipeline"
  enable_ipv6 = true
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vs" {
  description  = "pipelineExecution-start"
  vpc_id       = alicloud_vpc.pipelineExecution-vpc.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = format("%s1", var.name)
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_ecs_image_pipeline" "pipelineExection-pipeline" {
  base_image_type            = "IMAGE"
  description                = "example"
  system_disk_size           = "40"
  vswitch_id                 = alicloud_vswitch.vs.id
  add_account                = ["1284387915995949"]
  image_name                 = "example-image-pipeline"
  delete_instance_on_failure = true
  internet_max_bandwidth_out = "5"
  to_region_id               = ["cn-beijing"]
  base_image                 = "aliyun_3_x64_20G_dengbao_alibase_20240819.vhd"
  build_content              = "COMPONENT ic-bp122acttbs2sxdyq2ky"
}


resource "alicloud_ecs_image_pipeline_execution" "default" {
  image_pipeline_id = alicloud_ecs_image_pipeline.pipelineExection-pipeline.id
}
```

### Deleting `alicloud_ecs_image_pipeline_execution` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ecs_image_pipeline_execution`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `image_pipeline_id` - (Required, ForceNew) The ID of the image template.
* `status` - (Optional, Computed) The status of the image build task. You can set multiple values at the same time. Each value is separated by a comma (,). The format is 'BUILDING, distribut '. Value range:
  - PREPARING: PREPARING. Create resources such as temporary transit instances.
  - REPAIRING: REPAIRING. Repair the source image.
  - BUILDING: under construction. Execute user-defined commands and create images.
  - TESTING: TESTING. Execute user-defined test commands.
  - DISTRIBUTING: DISTRIBUTING. Perform mirror replication and sharing.
  - RELEASING: Resource Recovery. Temporary resources generated during the build process.
  - SUCCESS: SUCCESS. Build successfully.
  - PARTITION_SUCCESS: Partial success. The image has been built successfully, but there may be an exception in the distribution or resource cleanup steps.
  - FAILED: FAILED. Failed to build image.
  - TEST_FAILED: The test failed. The image was created successfully, but the test failed.
  - Canceling: canceling. The build process is being canceled.
  - Canceled: canceled. The build process has been canceled.

-> **NOTE:**  When the parameter value is empty, the image build tasks in all states are queried by default.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the image build task was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Image Pipeline Execution.
* `update` - (Defaults to 5 mins) Used when update the Image Pipeline Execution.

## Import

ECS Image Pipeline Execution can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_image_pipeline_execution.example <id>
```