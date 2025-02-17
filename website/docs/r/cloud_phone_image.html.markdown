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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Image.
* `delete` - (Defaults to 5 mins) Used when delete the Image.
* `update` - (Defaults to 5 mins) Used when update the Image.

## Import

Cloud Phone Image can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_phone_image.example <id>
```