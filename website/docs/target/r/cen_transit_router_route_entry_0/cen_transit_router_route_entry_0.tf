# Create a new tr-attachment and use it to attach one transit router to a new CEN
variable "name" {
  default = "tf-testAccCenTransitRouter"
}

variable "transit_router_route_entry_destination_cidr_block_attachment" {
  default = "192.168.0.0/24"
}

variable "transit_router_route_entry_name" {
  default = "sdk_rebot_cen_tr_yaochi"
}

variable "transit_router_route_entry_description" {
  default = "sdk_rebot_cen_tr_yaochi"
}

resource "alicloud_cen_instance" "cen" {
  name        = var.name
  description = "terraform01"
}

resource "alicloud_cen_transit_router" "default" {
  name   = var.name
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}

resource "alicloud_cen_transit_router_route_entry" "default" {
  transit_router_route_table_id                     = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
  transit_router_route_entry_destination_cidr_block = var.transit_router_route_entry_destination_cidr_block_attachment
  transit_router_route_entry_next_hop_type          = "Attachment"
  transit_router_route_entry_name                   = var.transit_router_route_entry_name
  transit_router_route_entry_description            = var.transit_router_route_entry_description
  transit_router_route_entry_next_hop_id            = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
}
