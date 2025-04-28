---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_image_import"
sidebar_current: "docs-alicloud-resource-image-import"
description: |-
  Provides a Alicloud ECS Image Import resource.
---

# alicloud_image_import

Provides a ECS Image Import resource.

For information about ECS Image Import and how to use it, see [What is Image Import](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-importimage).

-> **NOTE:** Available since v1.69.0.

-> **NOTE:** You must upload the image file to the object storage OSS in advance.

-> **NOTE:** The region where the image is imported must be the same region as the OSS bucket where the image file is uploaded.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_image_import&exampleId=cbaeaecf-5db3-ee5b-a5d4-0702941f1900a4a7558e&activeTab=example&spm=docs.r.image_import.0.cbaeaecf5d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-image-import-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "fc/hello.zip"
  content = <<EOF
    # -*- coding: utf-8 -*-
    def handler(event, context):
    print "hello world"
    return 'hello world'
    EOF
}

resource "alicloud_image_import" "default" {
  architecture = "x86_64"
  os_type      = "linux"
  platform     = "Ubuntu"
  license_type = "Auto"
  image_name   = var.name
  description  = var.name
  disk_device_mapping {
    oss_bucket      = alicloud_oss_bucket.default.id
    oss_object      = alicloud_oss_bucket_object.default.id
    disk_image_size = 5
  }
}
```

## Argument Reference

The following arguments are supported:

* `architecture` - (Optional, ForceNew) The architecture of the image. Default value: `x86_64`. Valid values: `x86_64`, `i386`.
* `os_type` - (Optional, ForceNew) The type of the operating system. Default value: `linux`. Valid values: `windows`, `linux`.
* `platform` - (Optional, ForceNew) The operating system platform. More valid values refer to [ImportImage OpenAPI](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/importimage).
-> **NOTE:** Before provider version 1.197.0, the default value of `platform` is `Ubuntu`.
* `boot_mode` - (Optional, ForceNew, Available since v1.225.0) The boot mode of the image. Valid values: `BIOS`, `UEFI`.
* `license_type` - (Optional, ForceNew) The type of the license used to activate the operating system after the image is imported. Default value: `Auto`. Valid values: `Auto`, `Aliyun`, `BYOL`.
* `image_name` - (Optional) The name of the image. The `image_name` must be `2` to `128` characters in length. The `image_name` must start with a letter and cannot start with acs: or aliyun. The `image_name` cannot contain http:// or https://. The `image_name` can contain letters, digits, periods (.), colons (:), underscores (_), and hyphens (-).
* `description` - (Optional) The description of the image. The `description` must be 2 to 256 characters in length and cannot start with http:// or https://.
* `disk_device_mapping` - (Required, ForceNew, Set) The information about the custom image. See [`disk_device_mapping`](#disk_device_mapping) below.

### `disk_device_mapping`

The disk_device_mapping supports the following:

* `format` - (Optional, ForceNew) The format of the image. Valid values: `RAW`, `VHD`, `qcow2`.
* `oss_bucket` - (Optional, ForceNew) The OSS bucket where the image file is stored.
* `oss_object` - (Optional, ForceNew) The name (key) of the object that the uploaded image is stored as in the OSS bucket.
* `device` - (Optional, ForceNew) The device name of the disk.
* `disk_image_size` - (Optional, ForceNew, Int) The size of the disk. Default value: `5`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Image Import.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Image Import.
* `delete` - (Defaults to 20 mins) Used when delete the Image Import.

## Import

ECS Image Import can be imported using the id, e.g.

```shell
$ terraform import alicloud_image_import.example <id>
```
