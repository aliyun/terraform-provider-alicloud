---
subcategory: "Cloudauth"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloudauth_face_configs"
sidebar_current: "docs-alicloud-datasource-cloudauth-face-configs"
description: |-
  Provides a list of Cloudauth Face Configs to the user.
---

# alicloud\_cloudauth\_face\_configs

This data source provides the Cloudauth Face Configs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.137.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloudauth_face_config" "default" {
  biz_name = "example-value"
  biz_type = "example-value"
}

data "alicloud_cloudauth_face_configs" "default" {
  ids        = [alicloud_cloudauth_face_config.default.id]
  name_regex = alicloud_cloudauth_face_config.default.biz_name
}

output "face_config" {
  value = data.alicloud_cloudauth_face_configs.default.configs.0
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Face Config IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by biz_name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `configs` - A list of Cloudauth Face Configs. Each element contains the following attributes:
	* `biz_name` - Scene name.
	* `biz_type` - Scene type. **NOTE:** The biz_type cannot exceed 32 characters and can only use English letters, numbers and dashes (-).
	* `gmt_updated` - The Update Time.
