---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_bundle"
sidebar_current: "docs-alicloud-resource-ecd-bundle"
description: |-
  Provides a Alicloud ECD Bundle resource.
---

# alicloud\_ecd\_bundle

Provides a ECD Bundle resource.

For information about ECD Bundle and how to use it, see [What is Bundle](https://help.aliyun.com/document_detail/188883.html).

-> **NOTE:** Available in v1.170.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecd_images" "default" {
  image_type            = "SYSTEM"
  os_type               = "Windows"
  desktop_instance_type = "eds.hf.4c8g"
}

data "alicloud_ecd_desktop_types" "default" {
  instance_type_family = "eds.hf"
  cpu_count            = 4
  memory_size          = 8192
}
resource "alicloud_ecd_bundle" "default" {
  description                 = var.name
  desktop_type                = data.alicloud_ecd_desktop_types.default.ids.0
  bundle_name                 = var.name
  image_id                    = data.alicloud_ecd_images.default.ids.0
  user_disk_size_gib          = [70]
  root_disk_size_gib          = 80
  root_disk_performance_level = "PL1"
  user_disk_performance_level = "PL1"
}
```

## Argument Reference

The following arguments are supported:

* `bundle_name` - (Optional) The name of the bundle.
* `description` - (Optional)  The description of the bundle.
* `desktop_type` - (Required, ForceNew) The desktop type. You can call `alicloud_ecd_desktop_types` to query desktop type.
* `image_id` - (Required) The ID of the image.
* `language` - (Optional) The language. Valid values: `zh-CN`, `zh-HK`, `en-US`, `ja-JP`.
* `root_disk_performance_level` - (Optional, ForceNew) The root disk performance level. Valid values: `PL0`, `PL1`, `PL2`, `PL3`.
* `root_disk_size_gib` - (Required, ForceNew) The root disk size gib.
* `user_disk_performance_level` - (Optional, ForceNew) The user disk performance level. Valid values: `PL0`, `PL1`, `PL2`, `PL3`.
* `user_disk_size_gib` - (Required, ForceNew) The size of the data disk. Currently, only one data disk can be set. Unit: GiB. 
  - The size of the data disk that supports the setting corresponds to the specification. For more information, see [Overview of Desktop Specifications](https://help.aliyun.com/document_detail/188609.htm?spm=a2c4g.11186623.0.0.6406297bE0U5DG).
  - The data disk size (user_disk_size_gib) set in the template must be greater than the data disk size (data_disk_size) in the mirror.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bundle.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Bundle.
* `delete` - (Defaults to 1 mins) Used when delete the Bundle.
* `update` - (Defaults to 1 mins) Used when update the Bundle.

## Import

ECD Bundle can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_bundle.example <id>
```