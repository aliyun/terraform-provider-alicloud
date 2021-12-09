---
subcategory: "Elastic Container Instance (ECI)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_image_cache"
sidebar_current: "docs-alicloud-eci-image-cache"
description: |-
  Provides an Alicloud ECI Image Cache resource.
---

# alicloud\_eci\_image\_cache

An ECI Image Cache can help user to solve the time-consuming problem of image pull. For information about Alicloud ECI Image Cache and how to use it, see [What is Resource Alicloud ECI Image Cache](https://www.alibabacloud.com/help/doc-detail/146891.htm).

-> **NOTE:** Available in v1.89.0+.

-> **NOTE:** Each image cache corresponds to a snapshot, and the user does not delete the snapshot directly, otherwise the cache will fail.

## Example Usage

Basic Usage

```
resource "alicloud_eci_image_cache" "example" {
  image_cache_name  = "tf-test"
  images            = ["registry.cn-beijing.aliyuncs.com/sceneplatform/sae-image-xxxx:latest"]
  security_group_id = "sg-2zeef68b66fxxxx"
  vswitch_id        = "vsw-2zef9k7ng82xxxx"
  eip_instance_id   = "eip-uf60c7cqb2pcrkgxhxxxx"
}
```
## Argument Reference

The following arguments are supported:

* `image_cache_name` - (Required, ForceNew) The name of the image cache.
* `images` - (Required, ForceNew) The images to be cached. The image name must be versioned.
* `security_group_id` - (Required, ForceNew) The ID of the security group. You do not need to specify the same security group as the container group.
* `vswitch_id` - (Required, ForceNew) The ID of the VSwitch. You do not need to specify the same VSwitch as the container group.
* `eip_instance_id` - (Optional, ForceNew) The instance ID of the Elastic IP Address (EIP). If you want to pull images from the Internet, you must specify an EIP to make sure that the container group can access the Internet. You can also configure the network address translation (NAT) gateway. We recommend that you configure the NAT gateway for the Internet access. Refer to [Public Network Access Method](https://help.aliyun.com/document_detail/99146.html)
* `image_cache_size`   - (Optional, ForceNew) The size of the image cache. Default to `20`. Unit: GiB.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `retention_days` - (Optional, ForceNew) The retention days of the image cache. Once the image cache expires, it will be cleared. By default, the image cache never expires. Note: The image cache that fails to be created is retained for only one day.
* `zone_id` - (Optional, ForceNew) The zone id to cache image.
* `image_registry_credential` - (Optional, ForceNew) The Image Registry parameters about the image to be cached.

### Block image_registry_credential
* `server` - (Optional) The address of Image Registry without `http://` or `https://`.
* `user_name` - (Optional) The user name of Image Registry.
* `password` - (Optional) The password of the Image Registry.

## Attributes Reference

* `id` -The id of the image cache.
* `container_group_id` - The ID of the container group job that is used to create the image cache.
* `status` -The status of the image cache.

## Import

ECI Image Cache can be imported using the id, e.g.

```
$ terraform import alicloud_eci_image_cache.example abc123456
```
