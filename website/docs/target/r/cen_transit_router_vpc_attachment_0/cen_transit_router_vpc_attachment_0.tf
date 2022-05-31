variable "transit_router_attachment_name" {
  default = "sdk_rebot_cen_tr_yaochi"
}

variable "transit_router_attachment_description" {
  default = "sdk_rebot_cen_tr_yaochi"
}

data "alicloud_cen_transit_router_available_resources" "default" {

}

resource "alicloud_vpc" "default" {
  vpc_name   = "sdk_rebot_cen_tr_yaochi"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = "sdk_rebot_cen_tr_yaochi"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[0]
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = "sdk_rebot_cen_tr_yaochi"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.2.0/24"
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].slave_zones[0]
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = "sdk_rebot_cen_tr_yaochi"
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.id
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
