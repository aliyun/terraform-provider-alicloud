---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_bgp_networks"
sidebar_current: "docs-alicloud-datasource-vpc-bgp-networks"
description: |-
  Provides a list of Vpc Bgp Networks to the user.
---

# alicloud\_vpc\_bgp\_networks

This data source provides the Vpc Bgp Networks of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_bgp_networks" "ids" {
  ids = ["example_value"]
}
output "vpc_bgp_network_id_1" {
  value = data.alicloud_vpc_bgp_networks.ids.networks.0.id
}

data "alicloud_vpc_bgp_networks" "routerId" {
  router_id = "example_value"
}
output "vpc_bgp_network_id_2" {
  value = data.alicloud_vpc_bgp_networks.routerId.networks.0.id
}

data "alicloud_vpc_bgp_networks" "status" {
  status = "Available"
}
output "vpc_bgp_network_id_3" {
  value = data.alicloud_vpc_bgp_networks.status.networks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Bgp Network IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `router_id` - (Optional, ForceNew) The ID of the router to which the route table belongs.
* `status` - (Optional, ForceNew) The state of the advertised BGP network. Valid values: `Available`, `Pending`, `Deleting`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `networks` - A list of Vpc Bgp Networks. Each element contains the following attributes:
	* `dst_cidr_block` - Advertised BGP networks.
	* `id` - The ID of the Bgp Network. The value formats as `<router_id>:<dst_cidr_block>`.
	* `router_id` - The ID of the vRouter.
	* `status` - The state of the advertised BGP network.