---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_shared_resource"
sidebar_current: "docs-alicloud-resource-resource-manager-shared-resource"
description: |-
  Provides a Alicloud Resource Manager Shared Resource resource.
---

# alicloud_resource_manager_shared_resource

Provides a Resource Manager Shared Resource resource.

For information about Resource Manager Shared Resource and how to use it, see [What is Shared Resource](https://www.alibabacloud.com/help/en/resource-management/latest/api-resourcesharing-2020-01-10-associateresourceshare).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc" "example" {
  name       = "tf-accvpc"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "example" {
  availability_zone = "cn-hangzhou-g"
  cidr_block        = "192.168.0.0/16"
  vpc_id            = alicloud_vpc.example.id
}

resource "alicloud_resource_manager_resource_share" "example" {
  resource_share_name = "example_value"
}

resource "alicloud_resource_manager_shared_resource" "example" {
  resource_id       = alicloud_vswitch.example.id
  resource_share_id = alicloud_resource_manager_resource_share.example.resource_share_id
  resource_type     = "VSwitch"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, ForceNew) The resource ID need shared.
* `resource_share_id` - (Required, ForceNew) The resource share ID of resource manager.
* `resource_type` - (Required, ForceNew) The resource type of should shared. Valid values:
  - `VSwitch`. 
  - The following types are added after v1.173.0: `ROSTemplate` and `ServiceCatalogPortfolio`. 
  - The following types are added after v1.192.0: `PrefixList` and `Image`.  
  - The following types are added after v1.194.1: `PublicIpAddressPool`.
  - The following types are added after v1.208.0: `KMSInstance`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Shared Resource. The value is formatted `<resource_share_id>:<resource_id>:<resource_type>`.
* `status` - The status of the Shared Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the Shared Resource.
* `delete` - (Defaults to 11 mins) Used when delete the Shared Resource.

## Import

Resource Manager Shared Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_shared_resource.example <resource_share_id>:<resource_id>:<resource_type>
```
