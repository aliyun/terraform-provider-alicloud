---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam"
description: |-
  Provides a Alicloud Vpc Ipam Ipam resource.
---

# alicloud_vpc_ipam_ipam

Provides a Vpc Ipam Ipam resource.

IP Address Management.

For information about Vpc Ipam Ipam and how to use it, see [What is Ipam](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.232.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc_ipam_ipam" "default" {
  ipam_description      = "This is my first Ipam."
  ipam_name             = var.name
  resource_group_id     = alicloud_resource_manager_resource_group.defaultResourceGroup.id
  operating_region_list = ["cn-hangzhou"]
}
```

## Argument Reference

The following arguments are supported:
* `ipam_description` - (Optional) The description of IPAM.

  It must be 2 to 256 characters in length and must start with an uppercase letter or a Chinese character, but cannot start with 'http: // 'or 'https. If the description is not filled in, it is blank. The default value is blank.
* `ipam_name` - (Optional) The name of the resource.
* `operating_region_list` - (Required, Set) List of IPAM effective regions.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipam.
* `delete` - (Defaults to 5 mins) Used when delete the Ipam.
* `update` - (Defaults to 5 mins) Used when update the Ipam.

## Import

Vpc Ipam Ipam can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipam_ipam.example <id>
```