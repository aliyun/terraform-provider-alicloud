---
subcategory: "NAT Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_forward_entry"
description: |-
  Provides a Alicloud Nat Gateway Forward Entry resource.
---

# alicloud_forward_entry

Provides a Nat Gateway Forward Entry resource.

DNAT route table entry.

For information about Nat Gateway Forward Entry and how to use it, see [What is Forward Entry](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/CreateForwardEntry).

-> **NOTE:** Available since v1.40.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_forward_entry&exampleId=fa1768e7-9ee0-7971-14bb-abf6895b60b74f889671&activeTab=example&spm=docs.r.forward_entry.0.fa1768e79e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_nat_gateway" "default" {
  vpc_id           = alicloud_vpc.default.id
  nat_gateway_name = var.name
  nat_type         = "Enhanced"
  vswitch_id       = alicloud_vswitch.default.id
  network_type     = "intranet"
}

resource "alicloud_vpc_nat_ip" "default" {
  nat_ip         = "172.16.0.66"
  nat_ip_name    = var.name
  nat_gateway_id = alicloud_nat_gateway.default.id
  nat_ip_cidr    = alicloud_vswitch.default.cidr_block
}

resource "alicloud_forward_entry" "default" {
  forward_table_id   = alicloud_nat_gateway.default.forward_table_ids
  external_ip        = alicloud_vpc_nat_ip.default.nat_ip
  external_port      = "80"
  ip_protocol        = "tcp"
  internal_ip        = "172.16.0.115"
  internal_port      = "8080"
  forward_entry_name = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_forward_entry&spm=docs.r.forward_entry.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `external_ip` - (Required) - When querying DNAT entries of an Internet NAT gateway, this parameter indicates the Elastic IP address used in the DNAT entry to provide public network access.
  - When querying DNAT entries of a VPC NAT gateway, this parameter indicates the NAT IP address used for access from external networks.
* `external_port` - (Required) - The external port or port range that is used for port forwarding when you query DNAT entries of Internet NAT gateways. Valid values: `1` to `65535`.
  - If you want to query a port range, separate the first port and last port with a forward slash (/), such as 10/20.
  - If you set ExternalPort to a port range, you must also set InternalPort to a port range, and the number of ports specified by these parameters must be the same. For example, if you set ExternalPort to 10/20, you can set InternalPort to 80/90.
* `forward_entry_name` - (Optional) The name of the DNAT entry.
* `forward_table_id` - (Required, ForceNew) The ID of the DNAT table to which the DNAT entry belongs.
* `internal_ip` - (Required) - The private IP address.
  - The private IP address of the ECS instance that uses DNAT entries to communicate with the Internet when you query DNAT entries of Internet NAT gateways.
  - The private IP address that uses DNAT entries when you query DNAT entries of VPC NAT gateways.
* `internal_port` - (Required) - When you configure a DNAT entry for an Internet NAT gateway, this parameter specifies the internal port or port range that requires port forwarding. Valid values: `1` to `65535`.
  - When you configure a DNAT entry for a VPC NAT gateway, this parameter specifies the destination ECS instance port to be mapped. Valid values: `1` to `65535`.
* `ip_protocol` - (Required) The protocol type. Valid values:
  - `tcp`: forwards TCP packets.
  - `udp`: forwards UDP packets.
  - `any`: forwards packets of all protocols. If `IpProtocol` is set to `Any`, both `ExternalPort` and `InternalPort` must also be set to `Any` to implement DNAT IP mapping.
* `port_break` - (Optional) Specifies whether to enable port break. Valid values:
  - `true`: Enables port break.
  - `false` (default): Disables port break.
-> **NOTE:**  If a DNAT entry and an SNAT entry share the same public IP address and you need to configure a port number greater than 1024, you must set `port_break` to `true`. `port_break` is only evaluated during resource creation and update. Modifying it in isolation will not trigger any action.
* `name` - (Optional, Deprecated since v1.119.1) Field `name` has been deprecated from provider version 1.119.1. New field `forward_entry_name` instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<forward_table_id>:<forward_entry_id>`.
* `forward_entry_id` - (Available since v1.43.0) The id of the forward entry on the server.
* `status` - (Available since v1.119.1) The status of forward entry.

## Timeouts

-> **NOTE:** Available since v1.119.1.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Forward Entry.
* `delete` - (Defaults to 5 mins) Used when delete the Forward Entry.
* `update` - (Defaults to 5 mins) Used when update the Forward Entry.

## Import

Nat Gateway Forward Entry can be imported using the id, e.g.

```shell
$ terraform import alicloud_forward_entry.example <forward_table_id>:<forward_entry_id>
```
