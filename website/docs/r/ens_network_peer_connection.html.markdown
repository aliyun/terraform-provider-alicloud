---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network_peer_connection"
description: |-
  Provides a Alicloud ENS Network Peer Connection resource.
---

# alicloud\_ens\_network\_peer\_connection

Provides a ENS Network Peer Connection resource.

A network peer connection connects two ENS networks so that resources in the two networks can communicate with each other.

For information about ENS Network Peer Connection and how to use it, see [CreateNetworkPeerConnection](https://www.alibabacloud.com/help/en/ens/developer-reference/api-ens-2017-11-10-createnetworkpeerconnection).

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_network" "default" {
  network_name  = "${var.name}-net1"
  description   = var.name
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "accepting" {
  network_name  = "${var.name}-net2"
  description   = var.name
  cidr_block    = "192.168.3.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network_peer_connection" "default" {
  network_peer_connection_name = var.name
  description                  = var.name
  network_id                   = alicloud_ens_network.default.id
  accepting_network_id         = alicloud_ens_network.accepting.id
}
```

## Argument Reference

The following arguments are supported:

* `network_id` - (Required, ForceNew) The network ID of the peer-to-peer connection initiator.
* `accepting_network_id` - (Required, ForceNew) The network ID of the peer-to-peer connection receiver.
* `network_peer_connection_name` - (Optional) The name of the network peering connection instance. The length is 1 to 128 characters, and cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the network peering connection instance. The length is 1 to 256 characters, and cannot start with `http://` or `https://`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. It is the instance ID of the network peer connection.
* `instance_id` - The instance ID of the network peer connection.
* `ens_region_id` - The ID of the ENS node.
* `status` - The status of the network peer connection. Valid values: `Creating`, `Activated`, `Deleting`.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Network Peer Connection.
* `delete` - (Defaults to 5 mins) Used when delete the Network Peer Connection.
* `update` - (Defaults to 5 mins) Used when update the Network Peer Connection.

## Import

ENS Network Peer Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_network_peer_connection.example <instance_id>
```
