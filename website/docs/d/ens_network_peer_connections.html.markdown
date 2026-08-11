---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network_peer_connections"
sidebar_current: "docs-alicloud-datasource-ens-network-peer-connections"
description: |-
  Provides a list of Ens Network Peer Connections to the user.
---

# alicloud\_ens\_network\_peer\_connections

This data source provides the Ens Network Peer Connections of current Alibaba Cloud user.

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ens_network_peer_connections" "default" {
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  name_regex    = "^my-peer-connection"
}

output "ens_network_peer_connection_id_1" {
  value = data.alicloud_ens_network_peer_connections.default.connections.0.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, ForceNew) The instance ID of the network peer connection.
* `ens_region_id` - (Optional, ForceNew) The ID of the ENS node.
* `network_id` - (Optional, ForceNew) The network ID of the peer-to-peer connection initiator.
* `network_peer_connection_name` - (Optional, ForceNew) The name of the network peering connection instance.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by the network peer connection name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Network Peer Connection IDs.
* `names` - A list of Network Peer Connection names.
* `connections` - A list of Ens Network Peer Connections. Each element contains the following attributes:
  * `instance_id` - The instance ID of the network peer connection.
  * `network_id` - The network ID of the peer-to-peer connection initiator.
  * `accepting_network_id` - The network ID of the peer-to-peer connection receiver.
  * `ens_region_id` - The ID of the ENS node.
  * `status` - The status of the network peer connection. Valid values: `Creating`, `Activated`, `Deleting`.
  * `create_time` - The creation time of the resource.
  * `description` - The description of the network peering connection instance.
  * `network_peer_connection_name` - The name of the network peering connection instance.
