---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_shared_resource"
description: |-
  Provides a Alicloud Resource Manager Shared Resource resource.
---

# alicloud_resource_manager_shared_resource

Provides a Resource Manager Shared Resource resource.



For information about Resource Manager Shared Resource and how to use it, see [What is Shared Resource](https://www.alibabacloud.com/help/en/resource-management/latest/api-resourcesharing-2020-01-10-associateresourceshare).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_shared_resource&exampleId=ded62a21-707b-a871-218b-b2bed5d0e1b274196330&activeTab=example&spm=docs.r.resource_manager_shared_resource.0.ded62a2170&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}-${random_integer.default.result}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "192.168.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  vswitch_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_resource_manager_resource_share" "default" {
  resource_share_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_resource_manager_shared_resource" "default" {
  resource_share_id = alicloud_resource_manager_resource_share.default.id
  resource_id       = alicloud_vswitch.default.id
  resource_type     = "VSwitch"
}
```

## Argument Reference

The following arguments are supported:
* `resource_id` - (Required, ForceNew) The ID of the shared resource.
* `resource_share_id` - (Required, ForceNew) The ID of the resource share.
* `resource_type` - (Required, ForceNew) The type of the shared resource. Valid values:
  - `VSwitch`. 
  - The following types are added after v1.173.0: `ROSTemplate` and `ServiceCatalogPortfolio`. 
  - The following types are added after v1.192.0: `PrefixList` and `Image`.  
  - The following types are added after v1.194.1: `PublicIpAddressPool`.
  - The following types are added after v1.208.0: `KMSInstance`.
  - The following types are added after v1.240.0: `Snapshot`.
  - For more information about the types of resources that can be shared, see [Services that work with Resource Sharing](https://help.aliyun.com/zh/resource-management/resource-sharing/product-overview/services-that-work-with-resource-sharing?spm=api-workbench.API%20Document.0.0.32fff3cdFveEud)

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<resource_share_id>:<resource_id>:<resource_type>`.
* `create_time` - (Available since v1.259.0) The time when the shared resource was associated with the resource share.
* `status` - The status of the Shared Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Shared Resource.
* `delete` - (Defaults to 10 mins) Used when delete the Shared Resource.

## Import

Resource Manager Shared Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_shared_resource.example <resource_share_id>:<resource_id>:<resource_type>
```
