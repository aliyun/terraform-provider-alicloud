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
resource "alicloud_image_export" "default" {
  image_id           = "m-bp1gxy***"
  oss_bucket         = "ecsimageexportconfig"
  oss_prefix         = "ecsExport"
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
