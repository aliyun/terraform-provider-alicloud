---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_images"
sidebar_current: "docs-alicloud-datasource-ecd-images"
description: |-
  Provides a list of Ecd Images to the user.
---

# alicloud\_ecd\_images

This data source provides the Ecd Images of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = "your_simple_office_site_name"
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = "your_policy_group_name"
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = "example_value"
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = "example_value"
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}

resource "alicloud_ecd_desktop" "default" {
  office_site_id  = alicloud_ecd_simple_office_site.default.id
  policy_group_id = alicloud_ecd_policy_group.default.id
  bundle_id       = data.alicloud_ecd_bundles.default.bundles.1.id
  desktop_name    = "your_desktop_name"
}

resource "alicloud_ecd_image" "default" {
  image_name  = "your_image_name"
  desktop_id  = alicloud_ecd_desktop.default.id
  description = "example_value"
}

data "alicloud_ecd_images" "ids" {
  ids = [alicloud_ecd_image.default.id]
}
output "ecd_image_id_1" {
  value = data.alicloud_ecd_images.ids.images.0.id
}

data "alicloud_ecd_images" "nameRegex" {
  name_regex = alicloud_ecd_image.default.image_name
}
output "ecd_image_id_2" {
  value = data.alicloud_ecd_images.nameRegex.images.0.id
}
```

## Argument Reference

The following arguments are supported:


* `ids` - (Optional, ForceNew, Computed)  A list of Image IDs.
* `image_type` - (Optional, ForceNew) The image type of the image. Valid values: `SYSTEM`, `CUSTOM`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Image name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the image. Valid values: `Creating`, `Available`, `CreateFailed`.
* `os_type` - (Optional, ForceNew, Available in 1.170.0+) The operating system type of the image. Valid values: `Windows` and `Linux`.
* `desktop_instance_type` - (Optional, ForceNew, Available in 1.170.0+) The desktop type of the image.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Image names.
* `images` - A list of Ecd Images. Each element contains the following attributes:
	* `create_time` - The creation time of the image.
	* `data_disk_size` - The size of data disk of the image.
	* `description` - The description of the image.
	* `gpu_category` - The Gpu Category of the image.
	* `id` - The ID of the Image.
	* `image_id` - The image id of the image.
	* `image_name` - The image name.
	* `image_type` - The image type of the image. Valid values: `SYSTEM`, `CUSTOM`.
	* `os_type` - The os type of the image.
	* `progress` - The progress of the image.
	* `size` - The size of the image.
	* `status` - The status of the image. Valid values: `Creating`, `Available`, `CreateFailed`.