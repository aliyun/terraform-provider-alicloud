---
layout: "alicloud"
page_title: "Alicloud: alicloud_images"
sidebar_current: "docs-alicloud-datasource-images"
description: |-
    Provides a list of images available to the user.
---

# alicloud\_images

This data source provides available image resources. It contains user's private images, system images provided by Alibaba Cloud, 
other public images and the ones available on the image market. 

## Example Usage

```
data "alicloud_images" "images_ds" {
  owners = "system"
  name_regex = "^centos_6"
}

output "first_image_id" {
  value = "${data.alicloud_images.images_ds.images.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter resulting images by name. 
* `most_recent` - (Optional, type: bool) If more than one result are returned, select the most recent one.
* `owners` - (Optional) Filter results by a specific image owner. Valid items are `system`, `self`, `others`, `marketplace`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `images` - A list of images. Each element contains the following attributes:
  * `id` - ID of the image.
  * `architecture` - Platform type of the image system: i386 or x86_64.
  * `creation_time` - Time of creation.
  * `description` - Description of the image.
  * `image_owner_alias` - Alias of the image owner.
  * `os_name` - Display name of the OS.
  * `status` - Status of the image. Possible values: `UnAvailable`, `Available`, `Creating` and `CreateFailed`.
  * `size` - Size of the image.
  * `disk_device_mappings` - Description of the system with disks and snapshots under the image.
    * `device` - Device information of the created disk: such as /dev/xvdb.
    * `size` - Size of the created disk.
    * `snapshot_id` - Snapshot ID.
  * `product_code` - Product code of the image on the image market.
  * `is_subscribed` - Whether the user has subscribed to the terms of service for the image product corresponding to the ProductCode.
  * `image_version` - Version of the image.
  * `progress` - Progress of image creation, presented in percentages.
