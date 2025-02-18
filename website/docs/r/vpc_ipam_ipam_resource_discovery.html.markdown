---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam_resource_discovery"
description: |-
  Provides a Alicloud Vpc Ipam Ipam Resource Discovery resource.
---

# alicloud_vpc_ipam_ipam_resource_discovery

Provides a Vpc Ipam Ipam Resource Discovery resource.

IP Address Management Resource Discovery.

For information about Vpc Ipam Ipam Resource Discovery and how to use it, see [What is Ipam Resource Discovery](https://next.api.alibabacloud.com/document/VpcIpam/2023-02-28/CreateIpamResourceDiscovery).

-> **NOTE:** Available since v1.243.0.

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


resource "alicloud_vpc_ipam_ipam_resource_discovery" "default" {
  operating_region_list               = ["cn-hangzhou"]
  ipam_resource_discovery_description = "This is a custom IPAM resource discovery."
  ipam_resource_discovery_name        = "example_resource_discovery"
}
```

## Argument Reference

The following arguments are supported:
* `ipam_resource_discovery_description` - (Optional) The description of resource discovery.
* `ipam_resource_discovery_name` - (Optional) The name of the resource
* `operating_region_list` - (Required, Set) The list of operating regions for resource discovery.
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `tags` - (Optional, Map) Label list information.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the resource discovery was created.
* `region_id` - The region ID of the resource
* `status` - The status of the resource discovery instance. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipam Resource Discovery.
* `delete` - (Defaults to 5 mins) Used when delete the Ipam Resource Discovery.
* `update` - (Defaults to 5 mins) Used when update the Ipam Resource Discovery.

## Import

Vpc Ipam Ipam Resource Discovery can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipam_ipam_resource_discovery.example <id>
```