# Create a new instance-attachment and use it to attach one child instance to a new CEN
variable "name" {
  default = "tf-testAccCenTransitRouterVbrAttachment"
}

variable "vbr_id" {
  default = "vbr-xxxxxxxxxx"
}

variable "transit_router_attachment_name" {
  default = "tf-test"
}

variable "transit_router_attachment_description" {
  default = "tf-test"
}

resource "alicloud_cen_instance" "cen" {
  instance_name = var.name
  description   = "terraform01"
}

resource "alicloud_transit_router" "tr" {
  name   = var.name
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_cen_transit_router_vbr_attachment" "foo" {
  vbr_id                                = var.vbr_id
  cen_id                                = alicloud_cen_instance.cen.id
  transit_router_id                     = alicloud_transit_router.tr.transit_router_id
  auto_publish_route_enabled            = true
  transit_router_attachment_name        = var.transit_router_attachment_name
  transit_router_attachment_description = var.transit_router_attachment_description
}
