variable "name" {
  default = "tf-testAccCenTransitRouter"
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
