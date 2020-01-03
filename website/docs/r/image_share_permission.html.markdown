---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_image_share_permission"
sidebar_current: "docs-alicloud-resource-image-share-permission"
description: |-
  Provides an ECS image share permission resource.
---

# alicloud\_image\_share\_permission

Manage image sharing permissions. You can share your custom image to other Alibaba Cloud users. The user can use the shared custom image to create ECS instances (RunInstances) or replace the system disk (ReplaceSystemDisk) of the instance.

-> **NOTE:** You can only share your own custom images to other Alibaba Cloud users.

-> **NOTE:** Each custom image is shared to a maximum of 10 Alibaba Cloud accounts at one time. Therefore, the parameter AddAccount.n or the parameter RemoveAccount.n can pass in up to 10 Alibaba Cloud accounts at a time, and the system will ignore this parameter for more than 10 accounts.

-> **NOTE:** Each custom image can be shared with up to 50 Alibaba Cloud accounts. You can submit a ticket to share with more users.

-> **NOTE:** After creating an ECS instance (RunInstances) using a shared image, once the custom image owner releases the image sharing relationship or deletes the custom image (DeleteImage), the instance cannot initialize the system disk (ReInitDisk).

-> **NOTE:** Available in 1.68.0+.

## Example Usage

```
resource "alicloud_image_share_permission" "default" {
  image_id           = "m-bp1gxyhdswlsn18tu***"
  account_id         = "account1"
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required, ForceNew) The source image ID.
* `account_id` - (Required, ForceNew) Alibaba Cloud Account ID. It is used to share images.
   
   
 ## Attributes Reference0
 
 The following attributes are exported:
 
* `id` - ID of the image.

 ## Import
 
image can be imported using the id, e.g.

```
$ terraform import alicloud_image_share_permission.default m-uf66871ape***yg1q***
```