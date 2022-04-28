---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_pipeline_execution"
sidebar_current: "docs-alicloud-resource-ecs-image-pipeline-execution"
description: |-
  Provides a Alicloud ECS Image Pipeline Execution resource.
---

# alicloud\_ecs\_image\_pipeline\_execution

Provides a ECS Image Pipeline Execution resource.

For information about ECS Image Pipeline Execution and how to use it, see [What is Image Pipeline Execution](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/startimagepipelineexecution).

-> **NOTE:** Available in v1.166.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "default" {
  image_id = data.alicloud_images.default.ids.0
}
resource "alicloud_ecs_image_pipeline" "default" {
  base_image                 = data.alicloud_images.default.ids.0
  base_image_type            = "IMAGE"
  build_content              = "RUN yum update -y"
  delete_instance_on_failure = false
  image_name                 = var.name
  name                       = var.name
  description                = var.name
  instance_type              = data.alicloud_instance_types.default.ids.0
  internet_max_bandwidth_out = 20
  system_disk_size           = 40
  to_region_id               = ["cn-qingdao", "cn-zhangjiakou"]
  vswitch_id                 = data.alicloud_vswitches.default.ids.0
  resource_group_id          = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  tags = {
    Created = "TF"
    For     = "Acceptance-test"
  }
}
resource "alicloud_ecs_image_pipeline_execution" "example" {
  image_pipeline_id = alicloud_ecs_image_pipeline.default.id
}
```

## Argument Reference

The following arguments are supported:

* `image_pipeline_id` - (Required, ForceNew) The ID of the image template.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Image Pipeline Execution.
* `status` - The status of the mirror build task.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Image Pipeline Execution.

## Import

ECS Image Pipeline Execution can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_image_pipeline_execution.example <id>
```