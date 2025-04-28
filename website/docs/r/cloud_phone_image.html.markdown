---
subcategory: "Cloud Phone"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_phone_image"
description: |-
  Provides a Alicloud Cloud Phone Image resource.
---

# alicloud_cloud_phone_image

Provides a Cloud Phone Image resource.

Cloud phone image.

For information about Cloud Phone Image and how to use it, see [What is Image](https://next.api.alibabacloud.com/document/eds-aic/2023-09-30/CreateCustomImage).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_phone_image&exampleId=8c53cb85-4df1-0053-0193-d06f12b4f7c25a1c7f67&activeTab=example&spm=docs.r.cloud_phone_image.0.8c53cb854d&intl_lang=EN_US" target="_blank">
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

resource "alicloud_cloud_phone_policy" "defaultjZ1gi0" {
}

resource "alicloud_cloud_phone_instance_group" "defaultYHMlTO" {
  instance_group_spec = "acp.basic.small"
  policy_group_id     = alicloud_cloud_phone_policy.defaultjZ1gi0.id
  instance_group_name = "AutoCreateGroupName"
  period              = "1"
  number_of_instances = "1"
  charge_type         = "PostPaid"
  image_id            = "imgc-075cllfeuazh03tg9"
  period_unit         = "Hour"
  auto_renew          = false
  amount              = "1"
  auto_pay            = false
  gpu_acceleration    = false
}

resource "alicloud_cloud_phone_instance" "default04hhXk" {
  android_instance_group_id = alicloud_cloud_phone_instance_group.defaultYHMlTO.id
  android_instance_name     = "CreateInstanceName"
}


resource "alicloud_cloud_phone_image" "default" {
  image_name  = "ImageName"
  instance_id = alicloud_cloud_phone_instance.default04hhXk.id
}
```

## Argument Reference

The following arguments are supported:
* `image_name` - (Required) The image name.
* `instance_id` - (Required) The instance ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the mirror.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Image.
* `delete` - (Defaults to 5 mins) Used when delete the Image.
* `update` - (Defaults to 5 mins) Used when update the Image.

## Import

Cloud Phone Image can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_phone_image.example <id>
```