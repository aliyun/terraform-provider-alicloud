---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_share"
description: |-
  Provides a Alicloud Resource Manager Resource Share resource.
---

# alicloud_resource_manager_resource_share

Provides a Resource Manager Resource Share resource. RS resource sharing.

For information about Resource Manager Resource Share and how to use it, see [What is Resource Share](https://www.alibabacloud.com/help/en/doc-detail/94475.htm).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vs" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_resource_manager_resource_share" "default" {
  resource_share_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `allow_external_targets` - (Optional, Available since v1.210.0) Whether to allow sharing to accounts outside the resource directory. Value:
  - false (default): Only sharing within the resource directory is allowed.
  - true: Allow sharing to any account.
* `permission_names` - (Optional, Available since v1.210.0) Share permission name. When it is empty, the system automatically binds the default permissions associated with the resource type. For more information, see [Permission Library](~~ 465474 ~~).
* `resource_group_id` - (Optional, Computed, Available since v1.210.0) The ID of the resource group.
* `resource_share_name` - (Required) The name of resource share.
* `resources` - (Optional, Available since v1.210.0) List of shared resources. See [`resources`](#resources) below.
* `targets` - (Optional, Available since v1.210.0) Resource user.

### `resources`

The resources supports the following:
* `resource_id` - (Optional) The ID of the shared resource.
The value range of N: 1 to 5, that is, a maximum of 5 shared resources are added at a time.
-> **NOTE:**  'Resources.N.ResourceId' and'resources. N.ResourceType' appear in pairs and need to be set at the same time.
* `resource_type` - (Optional) Shared resource type.
The value range of N: 1 to 5, that is, a maximum of 5 shared resources are added at a time.
For the types of resources that support sharing, see [Cloud services that support sharing](~~ 450526 ~~).
-> **NOTE:**  'Resources.N.ResourceId' and'resources. N.ResourceType' appear in pairs and need to be set at the same time.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The create time of resource share.
* `resource_share_owner` - The owner of resource share,  `Self` and `OtherAccounts`.
* `status` - The status of resource share.  `Active`,`Deleted` and `Deleting`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Resource Share.
* `delete` - (Defaults to 5 mins) Used when delete the Resource Share.
* `update` - (Defaults to 5 mins) Used when update the Resource Share.

## Import

Resource Manager Resource Share can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_resource_share.example <id>
```