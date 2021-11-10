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

```
resource "alicloud_image_share_permission" "default" {
  image_id           = "m-bp1gxyh***"
  account_id         = "1234567890"
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

```
$ terraform import alicloud_image_share_permission.default m-uf66yg1q:123456789
```
