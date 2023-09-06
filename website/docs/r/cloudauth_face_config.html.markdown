---
subcategory: "Cloudauth"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloudauth_face_config"
sidebar_current: "docs-alicloud-resource-cloudauth-face-config"
description: |-
  Provides a Alicloud Cloudauth Face Config resource.
---

# alicloud_cloudauth_face_config

Provides a Cloudauth Face Config resource.

For information about Cloudauth Face Config and how to use it, see [What is Face Config](https://www.alibabacloud.com/help/en/document_detail/99173.html).

-> **NOTE:** Available since v1.137.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
provider "alicloud" {
  region = "cn-hangzhou"
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}
resource "alicloud_cloudauth_face_config" "example" {
  biz_name = format("%s-biz", var.name)
  biz_type = format("type-%s", random_integer.default.result)
}
```

## Argument Reference

The following arguments are supported:

* `biz_name` - (Required) Scene name.
* `biz_type` - (Required) Scene type. **NOTE:** The biz_type cannot exceed 32 characters and can only use English letters, numbers and dashes (-).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Face Config. Its value is same as `biz_type`.
* `gmt_modified` - Last Modified Date.

## Import

Cloudauth Face Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloudauth_face_config.example <lang>
```
