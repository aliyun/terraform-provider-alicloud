---
subcategory: "ECI"
layout: "alicloud"
page_title: "Alicloud: alicloud_eciopenapi_image_cache"
sidebar_current: "docs-alicloud-resource-eciopenapi-image-cache"
description: |-
   Create ECI image cache. 
   
# alicloud\_eciopenapi\_image\_cache

The `alicloud_eciopenapi_image_cache` resources provide the creation of eciopenapi image caches that are available in Alibaba Cloud accounts.

-> **NOTE:** Available in 1.77.0+

## Example Usage

```
variable "name" {
	default = "YourECIName"
}

data "alicloud_vpcs" "default" {
  is_default = true
}

data "alicloud_vswitches" "default" {
  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}

resource"alicloud_eciopenapi_image_cache" "default" {
  image_cache_name  = "name"
  images            = ["centos_6_10_x64_20G_alibase_20200xxx.vhd", "ubuntu_18_04_x64_20G_alibase_20200xxx.vhd"]
  vswitch_id        = "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"
  security_group_id = "${alicloud_security_group.default.id}"
}

```

## Argument Reference

The following arguments are supported:

* `image_cache_name` - (Required) The image cache name.
* `security_group_id` - (Required) The security group ID does not need to be consistent with the ECI created.
* `vswitch_id` - (Required) The switch ID is not required to be the same as the ECI created.
* `images` - (Required) The image to cache.
* `eip_instance_id` - (Optional) Flexible public network IP. If you want to pull the public network image, you need to ensure that ECI can access the public network, and users need to configure the public network IP. Users can also configure the switch NAT gateway. The latter is recommended.
* `image_cache_size` - (Optional) The size of the cache. Default is 20 GB.
* `resource_group_id` - (Optional) The resource group ID.
* `retention_days` - (Optional) Image cache retention time, expired will be cleaned up, never expires by default. Note: The failed image cache is only retained for one day.
* `zone_id` - (Optional) The instance zone ID.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `image_cache_id` - Image cache ID.
 
* `container_group_id` - Generate Cache Job (an ECI instance).
  