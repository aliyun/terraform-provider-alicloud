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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_image_pipeline&exampleId=0a8d95e4-07a9-e56f-43f9-c3992f0d2e1744f83f55&activeTab=example&spm=docs.r.ecs_image_pipeline.0.0a8d95e407&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "default" {
  image_id = data.alicloud_images.default.ids.0
}
data "alicloud_account" "default" {
}
resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_ecs_image_pipeline" "default" {
  add_account                = [data.alicloud_account.default.id]
  base_image                 = data.alicloud_images.default.ids.0
  base_image_type            = "IMAGE"
  build_content              = "RUN yum update -y"
  delete_instance_on_failure = false
  image_name                 = "terraform-example"
  name                       = "terraform-example"
  description                = "terraform-example"
  instance_type              = data.alicloud_instance_types.default.ids.0
  resource_group_id          = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  internet_max_bandwidth_out = 20
  system_disk_size           = 40
  to_region_id               = ["cn-qingdao", "cn-zhangjiakou"]
  vswitch_id                 = alicloud_vswitch.default.id
  tags = {
    Created = "TF"
    For     = "Acceptance-test"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_image_pipeline&spm=docs.r.ecs_image_pipeline.example&intl_lang=EN_US)

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

```shell
$ terraform import alicloud_ecs_image_pipeline.example <id>
```