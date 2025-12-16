---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_bundle"
sidebar_current: "docs-alicloud-resource-ecd-bundle"
description: |-
  Provides a Alicloud ECD Bundle resource.
---

# alicloud_ecd_bundle

Provides a ECD Bundle resource.

For information about ECD Bundle and how to use it, see [What is Bundle](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-ecd-2020-09-30-createbundle).

-> **NOTE:** Available since v1.170.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecd_bundle&exampleId=39e0cefd-8564-7f4a-156d-465eb181cfbf66ace291&activeTab=example&spm=docs.r.ecd_bundle.0.39e0cefd85&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecd_bundle&spm=docs.r.ecd_bundle.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The ID of the image.
* `desktop_type` - (Required, ForceNew) The desktop type. You can call `alicloud_ecd_desktop_types` to query desktop type.
* `root_disk_size_gib` - (Required, ForceNew) The root disk size gib.
* `user_disk_size_gib` - (Required, ForceNew) The size of the data disk. Currently, only one data disk can be set. Unit: GiB.
  - The size of the data disk that supports the setting corresponds to the specification. For more information, see [Overview of Desktop Specifications](https://help.aliyun.com/document_detail/188609.htm?spm=a2c4g.11186623.0.0.6406297bE0U5DG).
  - The data disk size (user_disk_size_gib) set in the template must be greater than the data disk size (data_disk_size) in the mirror.
* `root_disk_performance_level` - (Optional, ForceNew) The root disk performance level. Valid values: `PL0`, `PL1`, `PL2`, `PL3`.
* `user_disk_performance_level` - (Optional, ForceNew) The user disk performance level. Valid values: `PL0`, `PL1`, `PL2`, `PL3`.
* `language` - (Optional) The language. Valid values: `zh-CN`, `zh-HK`, `en-US`, `ja-JP`.
* `bundle_name` - (Optional) The name of the bundle.
* `description` - (Optional)  The description of the bundle.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bundle.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Bundle.
* `update` - (Defaults to 1 mins) Used when update the Bundle.
* `delete` - (Defaults to 1 mins) Used when delete the Bundle.

## Import

ECD Bundle can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_bundle.example <id>
```
