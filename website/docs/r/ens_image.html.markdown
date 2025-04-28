---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_image"
description: |-
  Provides a Alicloud ENS Image resource.
---

# alicloud_ens_image

Provides a ENS Image resource.



For information about ENS Image and how to use it, see [What is Image](https://www.alibabacloud.com/help/en/ens/developer-reference/api-ens-2017-11-10-createimage).

-> **NOTE:** Available since v1.216.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_image&exampleId=a4a5c27d-c610-17ef-f6cf-5ca936426d306a5c5b36&activeTab=example&spm=docs.r.ens_image.0.a4a5c27dc6&intl_lang=EN_US" target="_blank">
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
  image_name                = var.name
  instance_id               = alicloud_ens_instance.default.id
  delete_after_image_upload = "false"
}
```

## Argument Reference

The following arguments are supported:
* `delete_after_image_upload` - (Optional) Specifies whether to automatically release the instance after the image is packaged and uploaded. Only image builders are supported. Default value: `false`. Valid values:
  - `true`: When the instance is released, the image is released together with the instance.
  - `false`: When the instance is released, the image is retained and is not released together with the instance.
  Empty means false by default.
* `image_name` - (Required) The name of the image. The name must be 2 to 128 characters in length. The name can contain letters, digits, colons (:), underscores (_), and hyphens (-). It must start with a letter but cannot start with http:// or https://. The name can contain letters, digits, colons (:), underscores (_), and hyphens (-).
* `instance_id` - (Optional, ForceNew) The ID of the instance.
* `target_oss_region_id` - (Optional, ForceNew, Available since v1.247.0) The region of the target OSS where the image is to be stored.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Image.
* `create_time` - The image creation time.
* `status` - The state of the image.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 120 mins) Used when create the Image.
* `delete` - (Defaults to 5 mins) Used when delete the Image.
* `update` - (Defaults to 5 mins) Used when update the Image.

## Import

ENS Image can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_image.example <id>
```
