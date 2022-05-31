variable "name" {
  default = "tf-test"
}

data "alicloud_vpcs" "example" {
  is_default = true
}

resource "alicloud_cen_instance" "example" {
  name = var.name
}

resource "alicloud_cen_instance_attachment" "vpc" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = data.alicloud_vpcs.example.vpcs.0.id
  child_instance_type      = "VPC"
  child_instance_region_id = data.alicloud_vpcs.example.vpcs.0.region_id
}

resource "alicloud_cen_route_service" "this" {
  access_region_id = data.alicloud_vpcs.example.vpcs.0.region_id
  host_region_id   = data.alicloud_vpcs.example.vpcs.0.region_id
  host_vpc_id      = data.alicloud_vpcs.example.vpcs.0.id
  cen_id           = alicloud_cen_instance_attachment.vpc.instance_id
  host             = "100.118.28.52/32"
}
