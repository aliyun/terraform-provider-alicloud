---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_image_export"
sidebar_current: "docs-alicloud-resource-image-export"
description: |-
  Provides an ECS image export resource.
---

# alicloud\_image\_export

Export a custom image to the OSS bucket in the same region as the custom image.

-> **NOTE:** If you create an ECS instance using a mirror image and create a system disk snapshot again, exporting a custom image created from the system disk snapshot is not supported.

-> **NOTE:** Support for exporting custom images that include data disk snapshot information in the image. The number of data disks cannot exceed 4 and the maximum capacity of a single data disk cannot exceed 500 GiB.

-> **NOTE:** Before exporting the image, you must authorize the cloud server ECS official service account to write OSS permissions through RAM.

-> **NOTE:** Available in 1.68.0+.

## Example Usage

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners     = "system"
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

resource "alicloud_security_group" "default" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_name              = "terraform-example"
  security_groups            = [alicloud_security_group.default.id]
  vswitch_id                 = alicloud_vswitch.default.id
  instance_type              = data.alicloud_instance_types.default.ids[0]
  image_id                   = data.alicloud_images.default.ids[0]
  internet_max_bandwidth_out = 10
}

resource "alicloud_image" "default" {
  instance_id = alicloud_instance.default.id
  image_name  = "terraform-example"
  description = "terraform-example"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "default" {
  bucket = "example-value-${random_integer.default.result}"
}

resource "alicloud_image_export" "default" {
  image_id   = alicloud_image.default.id
  oss_bucket = alicloud_oss_bucket.default.id
  oss_prefix = "ecsExport"
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required, ForceNew) The source image ID.
* `oss_bucket` - (Required, ForceNew) Save the exported OSS bucket.
* `oss_prefix` - (Optional, ForceNew) The prefix of your OSS Object. It can be composed of numbers or letters, and the character length is 1 ~ 30.
   
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when exporting the image (until it reaches the initial `Available` status). 
   
   
## Attributes Reference0
 
 The following attributes are exported:
 
* `id` - ID of the image.
