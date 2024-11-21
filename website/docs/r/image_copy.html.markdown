---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_image_copy"
sidebar_current: "docs-alicloud-resource-image-copy"
description: |-
  Provides an ECS image copy resource.
---

# alicloud\_image\_copy

Copies a custom image from one region to another. You can use copied images to perform operations in the target region, such as creating instances (RunInstances) and replacing system disks (ReplaceSystemDisk).

-> **NOTE:** You can only copy the custom image when it is in the Available state.

-> **NOTE:** You can only copy the image belonging to your Alibaba Cloud account. Images cannot be copied from one account to another.

-> **NOTE:** If the copying is not completed, you cannot call DeleteImage to delete the image but you can call CancelCopyImage to cancel the copying.

-> **NOTE:** Available in 1.66.0+.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_image_copy&exampleId=42d87a67-8b5a-f558-c48e-b4c3970e8fe954386fd7&activeTab=example&spm=docs.r.image_copy.0.42d87a678b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  alias  = "sh"
  region = "cn-shanghai"
}
provider "alicloud" {
  alias  = "hz"
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  provider                    = alicloud.hz
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  provider             = alicloud.hz
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  provider   = alicloud.hz
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "default" {
  provider   = alicloud.hz
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  provider     = alicloud.hz
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  provider = alicloud.hz
  name     = "terraform-example"
  vpc_id   = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  provider                   = alicloud.hz
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_name              = "terraform-example"
  security_groups            = [alicloud_security_group.default.id]
  vswitch_id                 = alicloud_vswitch.default.id
  instance_type              = data.alicloud_instance_types.default.ids[0]
  image_id                   = data.alicloud_images.default.ids[0]
  internet_max_bandwidth_out = 10
}

resource "alicloud_image" "default" {
  provider    = alicloud.hz
  instance_id = alicloud_instance.default.id
  image_name  = "terraform-example"
  description = "terraform-example"
}

resource "alicloud_image_copy" "default" {
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
  
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when copying the image (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the image.
   
   
## Attributes Reference0
 
 The following attributes are exported:
 
* `id` - ID of the image.

## Import
 
image can be imported using the id, e.g.

```shell
$ terraform import alicloud_image_copy.default m-uf66871ape***yg1q***
```
