# Create a new tr-attachment and use it to attach one transit router to a new CEN
variable "name" {
  default = "tf-testAccCenTransitRouter"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = "terraform01"
}

resource "alicloud_cen_transit_router" "default" {
  transit_router_name = var.name
  cen_id              = alicloud_cen_instance.default.id
}
