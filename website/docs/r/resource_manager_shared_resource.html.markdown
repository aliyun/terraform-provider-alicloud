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
* `permission_name` - (Optional, Available since v1.268.0) The name of a permission. If you do not configure this parameter, the system automatically associates the default permission for the specified resource type with the resource share.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `resource_arn` - (Optional, ForceNew, Available since v1.268.0) Associated resource ARN.

-> **NOTE:**  This parameter is not available when the association type 'AssociationType' is the resource consumer 'Target'.

* `resource_id` - (Optional, ForceNew, Computed) The ID of the shared resource.
* `resource_share_id` - (Required, ForceNew) The ID of the resource share.
* `resource_type` - (Optional, ForceNew, Computed) The type of the shared resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<resource_share_id>:<resource_id>:<resource_type>`.
* `create_time` - The time when the shared resource was associated with the resource share.
* `status` - The association status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Shared Resource.
* `delete` - (Defaults to 10 mins) Used when delete the Shared Resource.

## Import

Resource Manager Shared Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_shared_resource.example <resource_share_id>:<resource_id>:<resource_type>
```