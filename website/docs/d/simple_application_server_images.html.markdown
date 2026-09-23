---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_images"
sidebar_current: "docs-alicloud-datasource-simple-application-server-images"
description: |-
  Provides a list of Simple Application Server Images to the user.
---

# alicloud\_simple\_application\_server\_images

This data source provides the Simple Application Server Images of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_simple_application_server_images" "instanceImageType" {
  instance_image_type = "system"
}
output "simple_application_server_image_id_1" {
  value = data.alicloud_simple_application_server_images.ids.images.0.id
}
```


The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Image IDs.
* `image_type` - (Optional, ForceNew) The image type. Valid values: `app`, `custom`, `system`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Image name.
* `platform` - (Available in v1.161.0) The platform of Image supported. Valid values: ["Linux", "Windows"].
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of  Image names.
* `images` - A list of Simple Application Server Images. Each element contains the following attributes:
	* `description` - The description of the image.
	* `id` - The ID of the Instance Image.
	* `image_id` - The ID of the image.
	* `image_name` - The name of the resource.
	* `platform` - (Available in v1.161.0) The platform of Plan supported.
	* `image_type` - The type of the image. Valid values: `app`, `custom`, `system`.
		* `system`: operating system (OS) image.
		* `app`: application image.
		* `custom`: custom image.
