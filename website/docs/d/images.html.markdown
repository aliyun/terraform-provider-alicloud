---
subcategory: "ECS"
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
  owners     = "system"
  name_regex = "^centos_6"
}

output "first_image_id" {
  value = "${data.alicloud_images.images_ds.images.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter resulting images by name.
* `status` - (Optional, Available in 1.95.0+) The status of the image. The following values are available, Separate multiple parameter values by using commas (,). Default value: `Available`. Valid values: 
    * `Creating`: The image is being created. 
    * `Waiting`: The image is waiting to be processed. 
    * `Available`: The image is available.
    * `UnAvailable`: The image is unavailable.
    * `CreateFailed`: The image failed to be created.
    * `Deprecated`: The image is discontinued.
* `image_id` - (Optional, Available in 1.145.0+) The ID of the image.
* `image_name` - (Optional, Available in 1.145.0+) The name of the image.
* `snapshot_id` - (Optional, Available in 1.95.0+) The ID of the snapshot used to create the custom image.
* `resource_group_id` - (Optional, Available in 1.95.0+) The ID of the resource group to which the custom image belongs.
* `image_family` - (Optional, Available in 1.95.0+) The name of the image family. You can set this parameter to query images of the specified image family. This parameter is empty by default.
* `instance_type` - (Optional, Available in 1.95.0+) The instance type for which the image can be used.
* `is_support_io_optimized` - (Optional, Available in 1.95.0+) Specifies whether the image can be used on I/O optimized instances.
* `is_support_cloud_init` - (Optional, Available in 1.95.0+) Specifies whether the image supports cloud-init.
* `os_type` - (Optional, Available in 1.95.0+) The operating system type of the image. Valid values: `windows` and `linux`.
* `architecture` - (Optional, Available in 1.95.0+) The image architecture. Valid values: `i386` and `x86_64`.
* `action_type` - (Optional, Available in 1.95.0+) The scenario in which the image will be used. Default value: `CreateEcs`. Valid values:                                                
    * `CreateEcs`: instance creation.
    * `ChangeOS`: replacement of the system disk or operating system.
* `usage` - (Optional, Available in 1.95.0+) Specifies whether to check the validity of the request without actually making the request. Valid values:                                           
    * `instance`: The image is already in use and running on an ECS instance.
    * `none`: The image is not in use.
* `dry_run` - (Optional, Available in 1.95.0+) Specifies whether the image is running on an ECS instance. Default value: `false`. Valid values:                                           
    * `true`: The validity of the request is checked but resources are not queried. Check items include whether your AccessKey pair is valid, whether RAM users are authorized, and whether the required parameters are specified. If the check fails, the corresponding error message is returned. If the check succeeds, the DryRunOperation error code is returned.
    * `false`: The validity of the request is checked, and a 2XX HTTP status code is returned and resources are queried if the check succeeds.    
* `most_recent` - (Optional, type: bool) If more than one result are returned, select the most recent one.
* `owners` - (Optional) Filter results by a specific image owner. Valid items are `system`, `self`, `others`, `marketplace`.
* `tags` - (Optional, Available in 1.95.0+) A mapping of tags to assign to the resource.
* `image_owner_id` - (Optional, Available in 1.165.0+) The ID of the Alibaba Cloud account to which the image belongs. This parameter takes effect only when you query shared images or community images.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

-> **NOTE:** At least one of the `name_regex`, `most_recent` and `owners` must be set.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of image IDs.
* `images` - A list of images. Each element contains the following attributes:
  * `id` - ID of the image.
  * `architecture` - Platform type of the image system: i386 or x86_64.
  * `creation_time` - Time of creation.
  * `description` - Description of the image.
  * `image_owner_alias` - Alias of the image owner.
  * `os_name` - Display Chinese name of the OS.
  * `os_name_en` - Display English name of the OS.
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
