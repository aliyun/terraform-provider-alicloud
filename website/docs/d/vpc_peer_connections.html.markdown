---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_peer_connections"
sidebar_current: "docs-alicloud-datasource-vpc-peer-connections"
description: |-
  Provides a list of Vpc Peer Connections to the user.
---

# alicloud_vpc_peer_connections

This data source provides the Vpc Peer Connections of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.186.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_peer_connections" "ids" {}
output "vpc_peer_connection_id_1" {
  value = data.alicloud_vpc_peer_connections.ids.connections.0.id
}

data "alicloud_vpc_peer_connections" "nameRegex" {
  name_regex = "^my-PeerConnection"
}
output "vpc_peer_connection_id_2" {
  value = data.alicloud_vpc_peer_connections.nameRegex.connections.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of PeerConnection IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by PeerConnection name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `peer_connection_name` - (Optional, ForceNew) The name of the resource.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Accepting`, `Activated`, `Creating`, `Deleted`, `Deleting`, `Expired`, `Rejected`, `Updating`.
* `vpc_id` - (Optional, ForceNew) The ID of the requester VPC.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of PeerConnection names.
* `connections` - A list of Vpc Peer Connections. Each element contains the following attributes:
  * `accepting_ali_uid` - The ID of the Alibaba Cloud account (primary account) of the receiving end of the VPC peering connection to be created.
  * `accepting_region_id` - The region ID of the recipient of the VPC peering connection to be created.
  * `accepting_vpc_id` - The VPC ID of the receiving end of the VPC peer connection.
  * `bandwidth` - The bandwidth of the VPC peering connection to be modified. Unit: Mbps.
  * `create_time` - The creation time of the resource.
  * `description` - The description of the VPC peer connection to be created.
  * `id` - The ID of the PeerConnection.
  * `peer_connection_id` - The first ID of the resource.
  * `peer_connection_name` - The name of the resource.
  * `status` - The status of the resource.
  * `vpc_id` - The ID of the requester VPC.