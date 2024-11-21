---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_entry"
sidebar_current: "docs-alicloud-resource-cen-transit_router_route_entry"
description: |-
  Provides a Alicloud CEN transit router route entry resource.
---

# alicloud_cen_transit_router_route_entry

Provides a CEN transit router route entry resource.[What is Cen Transit Router Route Entry](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-cbn-2017-09-12-createtransitrouterrouteentry)

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_route_entry&exampleId=1903c1c5-f6e2-90bf-7d28-3d44b13102b1bd7533a1&activeTab=example&spm=docs.r.cen_transit_router_route_entry.0.1903c1c5f6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "tf_example"
}
resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
  description       = "an example for cen"
}

resource "alicloud_cen_transit_router" "example" {
  transit_router_name = var.name
  cen_id              = alicloud_cen_instance.example.id
}

resource "alicloud_cen_transit_router_route_table" "example" {
  transit_router_id = alicloud_cen_transit_router.example.transit_router_id
}

data "alicloud_express_connect_physical_connections" "example" {
  name_regex = "^preserved-NODELETING"
}
resource "random_integer" "vlan_id" {
  max = 2999
  min = 1
}
resource "alicloud_express_connect_virtual_border_router" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.example.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = random_integer.vlan_id.id
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}
resource "alicloud_cen_transit_router_vbr_attachment" "example" {
  vbr_id                                = alicloud_express_connect_virtual_border_router.example.id
  cen_id                                = alicloud_cen_instance.example.id
  transit_router_id                     = alicloud_cen_transit_router.example.transit_router_id
  auto_publish_route_enabled            = true
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
}

resource "alicloud_cen_transit_router_route_entry" "example" {
  transit_router_route_table_id                     = alicloud_cen_transit_router_route_table.example.transit_router_route_table_id
  transit_router_route_entry_destination_cidr_block = "192.168.0.0/24"
  transit_router_route_entry_next_hop_type          = "Attachment"
  transit_router_route_entry_name                   = var.name
  transit_router_route_entry_description            = var.name
  transit_router_route_entry_next_hop_id            = alicloud_cen_transit_router_vbr_attachment.example.transit_router_attachment_id
}
```
## Argument Reference

The following arguments are supported:

* `transit_router_route_table_id` - (Required, ForceNew) The ID of the transit router route table.
* `transit_router_route_entry_destination_cidr_block` - (Required, ForceNew) The CIDR of the transit router route entry.
* `transit_router_route_entry_next_hop_type` - (Required, ForceNew) The Type of the transit router route entry next hop,Valid values `Attachment` and `BlackHole`.
* `transit_router_route_entry_name` - (Optional) The name of the transit router route entry.
* `transit_router_route_entry_description` - (Optional) The description of the transit router route entry.
* `transit_router_route_entry_next_hop_id` - (Optional, ForceNew) The ID of the transit router route entry next hop.
* `dry_run` - (Optional) The dry run.

-> **NOTE:** If transit_router_route_entry_next_hop_type is `Attachment`, transit_router_route_entry_next_hop_id is required.
             If transit_router_route_entry_next_hop_type is `BlackHole`, transit_router_route_entry_next_hop_id cannot be filled.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_route_table_id>:<transit_router_route_entry_id>`.
* `transit_router_route_entry_id` - The ID of the route entry.
* `status` - The associating status of the Transit Router.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when creating the cen transit router route entry (until it reaches the initial `Active` status).
* `update` - (Defaults to 6 mins) Used when update the cen transit router route entry.
* `delete` - (Defaults to 6 mins) Used when delete the cen transit router route entry.

## Import

CEN instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_route_entry.default vtb-*********:rte-*******
```
