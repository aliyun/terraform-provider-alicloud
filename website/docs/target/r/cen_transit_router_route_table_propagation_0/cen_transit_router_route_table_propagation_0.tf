variable "transit_router_attachment_name" {
  default = "sdk_rebot_cen_tr_yaochi"
}

variable "transit_router_attachment_description" {
  default = "sdk_rebot_cen_tr_yaochi"
}

data "alicloud_cen_transit_router_available_resource" "default" {
}

resource "alicloud_vpc" "default" {
  name       = "sdk_rebot_cen_tr_yaochi"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  name              = "sdk_rebot_cen_tr_yaochi"
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "192.168.1.0/24"
  availability_zone = data.alicloud_cen_transit_router_available_resource.default.zones.0.master_zones.0
}

resource "alicloud_vswitch" "default_slave" {
  name              = "sdk_rebot_cen_tr_yaochi"
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "192.168.2.0/24"
  availability_zone = data.alicloud_cen_transit_router_available_resource.default.zones.0.slave_zones.0
}

resource "alicloud_cen_instance" "default" {
  name             = "sdk_rebot_cen_tr_yaochi"
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vpc_id            = alicloud_vpc.default.id
  zone_mapping {
    zone_id    = data.alicloud_cen_transit_router_available_resource.default.zones.0.master_zones.0
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mapping {
    zone_id    = data.alicloud_cen_transit_router_available_resource.default.zones.0.slave_zones.0
    vswitch_id = alicloud_vswitch.default_slave.id
  }
  transit_router_attachment_name        = var.transit_router_attachment_name
  transit_router_attachment_description = var.transit_router_attachment_description
}

resource "alicloud_cen_transit_router_route_table_propagation" "default" {
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachmentid
}
