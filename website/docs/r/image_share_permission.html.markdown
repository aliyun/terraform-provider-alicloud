---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_image_share_permission"
sidebar_current: "docs-alicloud-resource-image-share-permission"
description: |-
  Provides an ECS image share permission resource.
---

# alicloud\_image\_share\_permission

Manage image sharing permissions. You can share your custom image to other Alibaba Cloud users. The user can use the shared custom image to create ECS instances or replace the system disk of the instance.

-> **NOTE:** You can only share your own custom images to other Alibaba Cloud users.

-> **NOTE:** Each custom image can be shared with up to 50 Alibaba Cloud accounts. You can submit a ticket to share with more users.

-> **NOTE:** After creating an ECS instance using a shared image, once the custom image owner releases the image sharing relationship or deletes the custom image, the instance cannot initialize the system disk.

-> **NOTE:** Available in 1.68.0+.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_image_share_permission&exampleId=b9baaeb7-928e-38e8-47e0-1ce33e5f0c4e21991e94&activeTab=example&spm=docs.r.image_share_permission.0.b9baaeb792&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_name              = "terraform-example"
  security_groups            = [alicloud_security_group.default.id]
  vswitch_id                 = alicloud_vswitch.default.id
  instance_type              = data.alicloud_instance_types.default.ids[0]
  image_id                   = data.alicloud_images.default.ids[0]
  internet_max_bandwidth_out = 10
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_image" "default" {
  instance_id = alicloud_instance.default.id
  image_name  = "terraform-example-${random_integer.default.result}"
  description = "terraform-example"
}
variable "another_uid" {
  default = "123456789"
}
resource "alicloud_image_share_permission" "default" {
  image_id   = alicloud_image.default.id
  account_id = var.another_uid
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required, ForceNew) The source image ID.
* `account_id` - (Required, ForceNew) Alibaba Cloud Account ID. It is used to share images.
   
   
## Attributes Reference0
 
 The following attributes are exported:
 
* `id` - ID of the image. It formats as `<image_id>:<account_id>`

## Import
 
image can be imported using the id, e.g.

```shell
$ terraform import alicloud_image_share_permission.default m-uf66yg1q:123456789
```
