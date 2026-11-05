---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_copy_image"
sidebar_current: "docs-alicloud-resource-copy-image"
description: |-
  Deprecated alias of alicloud_image_copy. Copies a custom image from one region to another.
---

# alicloud\_copy_image

-> **DEPRECATED:** `alicloud_copy_image` is a deprecated alias of [alicloud_image_copy](https://www.terraform.io/docs/providers/alicloud/r/image_copy). Use `alicloud_image_copy` instead. The two resources share the same schema and backend `CopyImage` API; this page is retained only for backward compatibility and will be removed in a future release.

Copies a custom image from one region to another. You can use copied images to perform operations in the target region, such as creating instances (RunInstances) and replacing system disks (ReplaceSystemDisk).

For information about image copy and how to use it, see [CopyImage](https://www.alibabacloud.com/help/doc-detail/25538.html).

-> **NOTE:** You can only copy the custom image when it is in the Available state.

-> **NOTE:** You can only copy the image belonging to your Alibaba Cloud account. Images cannot be copied from one account to another.

-> **NOTE:** If the copying is not completed, you cannot call DeleteImage to delete the image but you can call CancelCopyImage to cancel the copying.

-> **NOTE:** Available since v1.66.0.

## Example Usage

```terraform
provider "alicloud" {
  alias  = "sh"
  region = "cn-shanghai"
}
provider "alicloud" {
  alias  = "hz"
  region = "cn-hangzhou"
}

resource "alicloud_image" "default" {
  provider    = alicloud.hz
  instance_id = alicloud_instance.default.id
  image_name  = "terraform-example"
  description = "terraform-example"
}

resource "alicloud_copy_image" "default" {
  provider         = alicloud.sh
  source_image_id  = alicloud_image.default.id
  source_region_id = "cn-hangzhou"
  image_name       = "terraform-example"
  description      = "terraform-example"
  tags = {
    FinanceDept = "FinanceDeptJoshua"
  }
}
```

## Argument Reference

The following arguments are supported:

* `source_image_id` - (Required, ForceNew) The source image ID.
* `source_region_id` - (Required, ForceNew) The ID of the region to which the source custom image belongs. You can call [DescribeRegions](https://www.alibabacloud.com/help/doc-detail/25609.htm) to view the latest regions of Alibaba Cloud.
* `image_name` - (Optional) The image name. It must be 2 to 128 characters in length, and must begin with a letter or Chinese character (beginning with http:// or https:// is not allowed). It can contain digits, colons (:), underscores (_), or hyphens (-). Default value: null.
* `description` - (Optional) The description of the image. It must be 2 to 256 characters in length and must not start with http:// or https://. Default value: null.
* `encrypted` - (Optional, ForceNew) Indicates whether to encrypt the image.
* `kms_key_id` - (Optional, ForceNew) Key ID used to encrypt the image.
* `tags` - (Optional) The tag value of an image. The value of N ranges from 1 to 20.
* `force` - (Optional) Indicates whether to force delete the custom image, Default is `false`.
  - true：Force deletes the custom image, regardless of whether the image is currently being used by other instances.
  - false：Verifies that the image is not currently in use by any other instances before deleting the image.
* `delete_auto_snapshot` - (Optional) Not the public attribute and it used to automatically delete dependence snapshots while deleting the image.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Optional, Deprecated) Attribute 'name' has been deprecated, please use `image_name` instead.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when copying the image (until it reaches the initial `Available` status).
* `delete` - (Defaults to 10 mins) Used when terminating the image.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the image.

## Import

Image can be imported using the id, e.g.

```shell
$ terraform import alicloud_copy_image.default m-uf66871ape***yg1q***
```
