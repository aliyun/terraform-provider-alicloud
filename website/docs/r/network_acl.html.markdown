---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_network_acl"
sidebar_current: "docs-alicloud-resource-network-acl"
description: |-
  Provides a Alicloud VPC Network Acl resource.
---

# alicloud_network_acl

Provides a VPC Network Acl resource. 

For information about VPC Network Acl and how to use it, see [What is Network Acl](https://www.alibabacloud.com/help/en/ens/latest/createnetworkacl).

-> **NOTE:** Available since v1.43.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_network_acl" "example" {
  vpc_id           = alicloud_vpc.example.id
  network_acl_name = var.name
  description      = var.name
  ingress_acl_entries {
    description            = "${var.name}-ingress"
    network_acl_entry_name = "${var.name}-ingress"
    source_cidr_ip         = "10.0.0.0/24"
    policy                 = "accept"
    port                   = "20/80"
    protocol               = "tcp"
  }
  egress_acl_entries {
    description            = "${var.name}-egress"
    network_acl_entry_name = "${var.name}-egress"
    destination_cidr_ip    = "10.0.0.0/24"
    policy                 = "accept"
    port                   = "20/80"
    protocol               = "tcp"
  }
  resources {
    resource_id   = alicloud_vswitch.example.id
    resource_type = "VSwitch"
  }
}
```


## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the network ACL.The description must be 1 to 256 characters in length and cannot start with http:// or https.
* `egress_acl_entries` - (Optional) Out direction rule information. See [`egress_acl_entries`](#egress_acl_entries) below.
* `ingress_acl_entries` - (Optional) Inward direction rule information. See [`ingress_acl_entries`](#ingress_acl_entries) below.
* `network_acl_name` - (Optional) The name of the network ACL.The name must be 1 to 128 characters in length and cannot start with http:// or https.
* `resources` - (Optional) The associated resource. See [`resources`](#resources) below.
* `tags` - (Optional, Map, Available since v1.206.0) The tags of this resource.
* `vpc_id` - (Required, ForceNew) The ID of the associated VPC.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated from v1.122.0+) Field 'name' has been deprecated from provider version 1.122.0. New field 'network_acl_name' instead.

### `egress_acl_entries`

The egress_acl_entries supports the following:

* `description` - (Optional) The description of the outbound rule.The description must be 1 to 256 characters in length and cannot start with http:// or https.
* `destination_cidr_ip` - (Optional) The network of the destination address.
* `network_acl_entry_name` - (Optional) Name of the outbound rule entry.The name must be 1 to 128 characters in length and cannot start with http:// or https.
* `policy` - (Optional) Authorization policy. Value:
  - accept: Allow.
  - drop: Refused.
* `port` - (Optional) The destination port range of the outbound rule.When the Protocol type of the outbound rule is all, icmp, or gre, the port range is - 1/-1, indicating that the port is not restricted.When the Protocol type of the outbound rule is tcp or udp, the port range is 1 to 65535, and the format is 1/200 or 80/80, indicating port 1 to port 200 or port 80.
* `protocol` - (Optional) The protocol type. Value:
  - icmp: Network Control Message Protocol.
  - gre: Generic Routing Encapsulation Protocol.
  - tcp: Transmission Control Protocol.
  - udp: User Datagram Protocol.
  - all: Supports all protocols.

### `ingress_acl_entries`

The ingress_acl_entries supports the following:

* `description` - (Optional) Description of the inbound rule.The description must be 1 to 256 characters in length and cannot start with http:// or https.
* `network_acl_entry_name` - (Optional) The name of the inbound rule entry.The name must be 1 to 128 characters in length and cannot start with http:// or https.
* `policy` - (Optional) Authorization policy. Value:
  - accept: Allow.
  - drop: Refused.
* `port` - (Optional) The source port range of the inbound rule.When the Protocol type of the inbound rule is all, icmp, or gre, the port range is - 1/-1, indicating that the port is not restricted.When the Protocol type of the inbound rule is tcp or udp, the port range is 1 to 65535, and the format is 1/200 or 80/80, indicating port 1 to port 200 or port 80.
* `protocol` - (Optional) The protocol type. Value:
  - icmp: Network Control Message Protocol.
  - gre: Generic Routing Encapsulation Protocol.
  - tcp: Transmission Control Protocol.
  - udp: User Datagram Protocol.
  - all: Supports all protocols.
* `source_cidr_ip` - (Optional) Source address network segment.

### `resources`

The resources supports the following:

* `resource_id` - (Required) The ID of the associated resource.
* `resource_type` - (Required) The type of the associated resource.
* `status` - (Required) The status of the associated resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `status` - The state of the network ACL.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Acl.
* `delete` - (Defaults to 5 mins) Used when delete the Network Acl.
* `update` - (Defaults to 5 mins) Used when update the Network Acl.

## Import

VPC Network Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_network_acl.example <id>
```