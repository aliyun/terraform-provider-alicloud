---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_custom_images"
sidebar_current: "docs-alicloud-datasource-simple-application-server-custom-images"
description: |-
  Provides a list of Simple Application Server Custom Images to the user.
---

# alicloud\_simple\_application\_server\_custom\_images

This data source provides the Simple Application Server Custom Images of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_simple_application_server_custom_images" "ids" {
  ids = ["example_id"]
}
output "simple_application_server_custom_image_id_1" {
  value = data.alicloud_simple_application_server_custom_images.ids.images.0.id
}

data "alicloud_simple_application_server_custom_images" "nameRegex" {
  name_regex = "^my-CustomImage"
}
output "simple_application_server_custom_image_id_2" {
  value = data.alicloud_simple_application_server_custom_images.nameRegex.images.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Custom Image IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Custom Image name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Custom Image names.
* `images` - A list of Simple Application Server Custom Images. Each element contains the following attributes:
	* `custom_image_id` - The first ID of the resource.
	* `custom_image_name` - The name of the resource.
	* `description` - Image description information.
	* `id` - The ID of the Custom Image.
	* `platform` - The type of operating system used by the Mirror. Valid values: `Linux`, `Windows`.