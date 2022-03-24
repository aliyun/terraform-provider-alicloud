---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_pipeline"
sidebar_current: "docs-alicloud-resource-ecs-image-pipeline"
description: |-
  Provides a Alicloud ECS Image Pipeline resource.
---

# alicloud\_ecs\_image\_pipeline

Provides a ECS Image Pipeline resource.

For information about ECS Image Pipeline and how to use it, see [What is Image Pipeline](https://www.alibabacloud.com/help/en/doc-detail/200427.html).

-> **NOTE:** Available in v1.163.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
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
  add_account                = ["example_value"]
  base_image                 = data.alicloud_images.default.ids.0
  base_image_type            = "IMAGE"
  build_content              = "RUN yum update -y"
  delete_instance_on_failure = false
  image_name                 = "example_value"
  name                       = "example_value"
  description                = "example_value"
  instance_type              = data.alicloud_instance_types.default.ids.0
  resource_group_id          = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  internet_max_bandwidth_out = 20
  system_disk_size           = 40
  to_region_id               = ["cn-qingdao", "cn-zhangjiakou"]
  vswitch_id                 = data.alicloud_vswitches.default.ids.0
  tags = {
    Created = "TF"
    For     = "Acceptance-test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `add_account` - (Optional, ForceNew) The ID of Alibaba Cloud account to which to share the created image.
* `base_image` - (Required, ForceNew) The source image. When you set `base_image_type` to `IMAGE`, set `base_image` to the ID of a custom image. When you set `base_image_type` to `IMAGE_FAMILY`, set `base_image` to the name of an image family.
* `base_image_type` - (Required, ForceNew) The type of the source image. Valid values: `IMAGE`, `IMAGE_FAMILY`.
  - IMAGE: custom image.
  - IMAGE_FAMILY: image family.
* `build_content` - (Optional, ForceNew) The content of the image template. The content cannot be greater than 16 KB in size, and can contain up to 127 commands.
* `delete_instance_on_failure` - (Optional, ForceNew) Specifies whether to release the intermediate instance if the image cannot be created.
* `description` - (Optional, ForceNew) The description of the image template. The description must be `2` to `256` characters in length and cannot start with `http://` or `https://`. **Note:** If the intermediate instance cannot be started, the instance is released by default.
* `image_name` - (Optional, ForceNew) The name prefix of the image to be created. The prefix must be `2` to `64` characters in length. It must start with a letter and cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (_), periods (.),and hyphens (-).
* `instance_type` - (Optional, ForceNew) The instance type of the instance. You can call the DescribeInstanceTypes operation to query instance types. If you do not specify this parameter, an instance type that provides the fewest vCPUs and memory resources is automatically selected. This configuration is subject to resource availability of instance types. For example, the `ecs.g6.large` instance type is selected by default. If available `ecs.g6.large` resources are insufficient, the `ecs.g6.xlarge` instance type is selected.
* `internet_max_bandwidth_out` - (Optional, ForceNew) The size of the outbound public bandwidth for the intermediate instance. Unit: `Mbit/s`. Valid values: `0` to `100`. Default value: `0`.
* `name` - (Optional, ForceNew) The name of the image template. The name must be `2` to `128` characters in length. It must start with a letter and cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (_), periods (.),and hyphens (-).
* `system_disk_size` - (Optional, ForceNew) The size of the system disk of the intermediate instance. Unit: GiB. Valid values: `20` to `500`. Default value: `40`.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `resource_group_id` - (Optional, ForceNew, Computed) The ID of the resource group.
* `to_region_id` - (Optional, ForceNew) The ID of region to which to distribute the created image.
* `vswitch_id` - (Optional, ForceNew) The ID of the vSwitch. If you do not specify this parameter, a virtual private cloud (VPC) and a vSwitch are created by default.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Image Pipeline.

## Import

ECS Image Pipeline can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_image_pipeline.example <id>
```