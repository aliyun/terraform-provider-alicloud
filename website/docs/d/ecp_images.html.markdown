---
subcategory: "Elastic Cloud Phone (ECP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecp_images"
sidebar_current: "docs-alicloud-datasource-ecp-images"
description: |-
  Provides a list of Ecp Images to the user.
---

# alicloud\_ecp\_images

This data source provides the Ecp Images of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.159.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecp_images" "ids" {}
output "ecp_image_id_1" {
  value = data.alicloud_ecp_images.ids.images.0.id
}

data "alicloud_ecp_images" "nameRegex" {
  name_regex = "^my-Image"
}
output "ecp_image_id_2" {
  value = data.alicloud_ecp_images.nameRegex.images.0.id
}


```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Image IDs.
* `image_category` - (Optional, ForceNew) Mirror source. Value range: system: Public image provided by Alibaba Cloud.
  self: The custom image you created. others: images shared with you by other Alibaba Cloud users.
* `image_name` - (Optional, ForceNew) The name of the mirror image. It must be 2 to 128 characters in length and must
  start with an uppercase letter or Chinese. It cannot start with http:// or https. It can contain Chinese, English,
  numbers, half-width colons (:), underscores (_), half-width periods (.), or dashes (-).
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Image name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) Query the mirror image in a certain state. Value range: Waiting: multitask queuing.
  Creating: The image is being created. Copying: The image is being copied. Importing: The image is being imported.
  Available: The image you can use. CreateFailed: The image that failed to be created.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Image names.
* `images` - A list of Ecp Images. Each element contains the following attributes:
    * `accounts` - User ID list.
    * `create_time` - Image creation time, in ISO 8601 format.
    * `description` - Mirror description.
    * `id` - The ID of the Image.
    * `image_category` - Mirror type.
    * `image_id` - Image ID.
    * `image_name` - The name of the mirror image.
    * `is_self_shared` - Whether the image has been shared with other users.
    * `os_name` - The Chinese display name of the operating system.
    * `os_name_en` - The English display name of the operating system.
    * `os_type` - Operating system type.
    * `platform` - Operating system distribution.
    * `progress` - The progress of mirror image production.
    * `status` - Mirror status.
    * `usage` - Whether the image is already running in the cloud phone instance. Value range:
      -none: The mirror image is idle and no cloud mobile phone instance is used for the time being. -instance: The
      image is running and is used by cloud mobile phone instances.