---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_image"
description: |-
  Provides a Alicloud ENS Image resource.
---

# alicloud_ens_image

Provides a ENS Image resource. 

For information about ENS Image and how to use it, see [What is Image](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.216.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_image&exampleId=1841648a-a784-d4be-baa1-166241099c9d86003eca&activeTab=example&spm=docs.r.ens_image.0.1841648aa7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ens_instance" "default" {
  system_disk {
    size = "20"
  }
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "PayAsYouGo"
  password                   = "12345678ABCabc"
  amount                     = "1"
  internet_max_bandwidth_out = "10"
  public_ip_identification   = true
  ens_region_id              = "cn-chenzhou-telecom_unicom_cmcc"
  period_unit                = "Month"
  instance_type              = "ens.sn1.stiny"
  status                     = "Stopped"
}


resource "alicloud_ens_image" "default" {
  image_name = var.name

  instance_id               = alicloud_ens_instance.default.id
  delete_after_image_upload = "false"
}
```

## Argument Reference

The following arguments are supported:
* `delete_after_image_upload` - (Optional) Whether the instance is automatically released after the image is packaged and uploaded successfully, only the build machine is supported.  Optional values: true: When the instance is released, the image is released together with the instance. false: When the instance is released, the image is retained and is not released together with the instance. Empty means false by default.
* `image_name` - (Required) Image Name.
* `instance_id` - (Optional, ForceNew) The ID of the instance corresponding to the image.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Image creation time.
* `status` - Mirror Status  Optional values: Creating: Creating Packing: Packing Uploading: Uploading Pack_failed: Packing failed Upload_failed: Upload failed Available: Only images in the Available state can be used and operated. Unavailable: Not available Copying: Copying.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Image.
* `delete` - (Defaults to 5 mins) Used when delete the Image.
* `update` - (Defaults to 5 mins) Used when update the Image.

## Import

ENS Image can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_image.example <id>
```