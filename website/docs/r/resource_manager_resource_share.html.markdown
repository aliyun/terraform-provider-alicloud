---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_share"
sidebar_current: "docs-alicloud-resource-resource-manager-resource-share"
description: |-
  Provides a Alicloud Resource Manager Resource Share resource.
---

# alicloud_resource_manager_resource_share

Provides a Resource Manager Resource Share resource.

For information about Resource Manager Resource Share and how to use it, see [What is Resource Share](https://www.alibabacloud.com/help/en/doc-detail/94475.htm).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_resource_manager_resource_share" "example" {
  resource_share_name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `resource_share_name` - (Required) The name of resource share.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Resource Share.
* `resource_share_owner` - The owner of the Resource Share.
* `status` - The status of the Resource Share.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Resource Share.
* `update` - (Defaults to 5 mins) Used when update the Resource Share.
* `delete` - (Defaults to 15 mins) Used when delete the Resource Share.

## Import

Resource Manager Resource Share can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_resource_share.example <id>
```
