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

For information about ECS Image Pipeline Execution and how to use it, see [What is Image Pipeline Execution](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-startimagepipelineexecution).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_image_pipeline_execution&exampleId=51ea8dfb-0104-8d61-3482-a8835b0b0df89a34328f&activeTab=example&spm=docs.r.ecs_image_pipeline_execution.0.51ea8dfb01&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `status` - (Optional, Computed) The status of the image build task. Valid values:
  - CANCELLED: canceled. The build process has been canceled.

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