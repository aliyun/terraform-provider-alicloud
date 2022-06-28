---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_bgp_peers"
sidebar_current: "docs-alicloud-datasource-vpc-bgp-peers"
description: |-
  Provides a list of Vpc Bgp Peers to the user.
---

# alicloud\_vpc\_bgp\_peers

This data source provides the Vpc Bgp Peers of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_bgp_peers" "ids" {
  ids = ["example_value-1", "example_value-2"]
}
output "vpc_bgp_peer_id_1" {
  value = data.alicloud_vpc_bgp_peers.ids.peers.0.id
}

data "alicloud_vpc_bgp_peers" "bgpGroupId" {
  bgp_group_id = "example_value"
}
output "vpc_bgp_peer_id_2" {
  value = data.alicloud_vpc_bgp_peers.bgpGroupId.peers.0.id
}

data "alicloud_vpc_bgp_peers" "routerId" {
  router_id = "example_value"
}
output "vpc_bgp_peer_id_3" {
  value = data.alicloud_vpc_bgp_peers.routerId.peers.0.id
}

data "alicloud_vpc_bgp_peers" "status" {
  status = "Available"
}
output "vpc_bgp_peer_id_4" {
  value = data.alicloud_vpc_bgp_peers.status.peers.0.id
}

```

## Argument Reference

The following arguments are supported:

* `bgp_group_id` - (Optional, ForceNew) The ID of the BGP group to which the BGP peer that you want to query belongs.
* `ids` - (Optional, ForceNew, Computed)  A list of Bgp Peer IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `router_id` - (Optional, ForceNew) The ID of the virtual border router (VBR) that is associated with the BGP peer that you want to query.
* `status` - (Optional, ForceNew) The status of the BGP peer. Valid values: `Available`, `Deleted`, `Deleting`, `Modifying`, `Pending`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `peers` - A list of Vpc Bgp Peers. Each element contains the following attributes:
	* `auth_key` - The authentication key of the BGP group.
	* `bfd_multi_hop` - The BFD hop count.
	* `bgp_group_id` - The ID of the BGP group.
	* `bgp_peer_id` - The ID of the BGP neighbor.
	* `bgp_peer_name` - The name of the BGP neighbor.
	* `bgp_status` - The status of the BGP connection.
	* `description` - The description of the BGP group.
	* `enable_bfd` - Indicates whether the Bidirectional Forwarding Detection (BFD) protocol is enabled.
	* `hold` - The hold time.
	* `id` - The ID of the Bgp Peer.
	* `ip_version` - The IP version.
	* `is_fake` - Indicates whether a fake AS number is used.
	* `keepalive` - The keepalive time.
	* `local_asn` - The AS number of the device on the Alibaba Cloud side.
	* `peer_asn` - The autonomous system (AS) number of the BGP peer.
	* `peer_ip_address` - The IP address of the BGP neighbor.
	* `route_limit` - The limit on routes.
	* `router_id` - The ID of the router.
	* `status` - The status of the BGP peer.