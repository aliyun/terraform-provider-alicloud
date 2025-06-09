---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_table_propagations"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-route-table-propagations"
description: |-
  Provides a list of CEN Transit Router Route Table Propagations to the user.
---

# alicloud_cen_transit_router_route_table_propagations

This data source provides the CEN Transit Router Route Table Propagations of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}

resource "random_integer" "default" {
  min = 1
  max = 2999
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.default.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = random_integer.default.id
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_cen_transit_router_vbr_attachment" "default" {
  cen_id                                = alicloud_cen_instance.default.id
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  vbr_id                                = alicloud_express_connect_virtual_border_router.default.id
  auto_publish_route_enabled            = true
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id               = alicloud_cen_transit_router.default.transit_router_id
  transit_router_route_table_name = var.name
}

resource "alicloud_cen_transit_router_route_table_propagation" "default" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
}

data "alicloud_cen_transit_router_route_table_propagations" "ids" {
  transit_router_route_table_id = alicloud_cen_transit_router_route_table_propagation.default.transit_router_route_table_id
  ids                           = [alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id]
}

output "cen_transit_router_route_table_propagation_id_0" {
  value = data.alicloud_cen_transit_router_route_table_propagations.ids.propagations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Transit Router Route Table Propagation IDs.
* `transit_router_route_table_id` - (Required, ForceNew) The ID of the route table of the Enterprise Edition transit router.
* `transit_router_attachment_id` - (Optional, ForceNew) The ID of the network instance connection.
* `status` - (Optional, ForceNew) The status of the route learning correlation. Valid values: `Active`, `Enabling`, `Disabling`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `propagations` - A list of Transit Router Route Table Propagations. Each element contains the following attributes:
  * `id` - The ID of the network instance connection.
  * `transit_router_attachment_id` - The ID of the network instance connection.
  * `transit_router_route_table_id` - The ID of the route table of the Enterprise Edition transit router.
  * `resource_id` - The ID of the network instance.
  * `resource_type` - The type of the network instance.
  * `status` - The status of the route learning correlation.
