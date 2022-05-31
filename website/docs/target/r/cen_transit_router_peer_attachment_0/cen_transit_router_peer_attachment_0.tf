variable "name" {
  default = "tf-testAcccExample"
}

provider "alicloud" {
  alias  = "us"
  region = "us-east-1"
}

provider "alicloud" {
  alias  = "cn"
  region = "cn-hangzhou"
}

resource "alicloud_cen_instance" "default" {
  provider          = alicloud.cn
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_bandwidth_package" "default" {
  bandwidth                  = 5
  cen_bandwidth_package_name = var.name
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "North-America"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  provider             = alicloud.cn
  instance_id          = alicloud_cen_instance.default.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
}

resource "alicloud_cen_transit_router" "cn" {
  provider   = alicloud.cn
  cen_id     = alicloud_cen_instance.default.id
  depends_on = [alicloud_cen_bandwidth_package_attachment.default]
}

resource "alicloud_cen_transit_router" "us" {
  provider   = alicloud.us
  cen_id     = alicloud_cen_instance.default.id
  depends_on = [alicloud_cen_transit_router.default_0]
}

resource "alicloud_cen_transit_router_peer_attachment" "default" {
  provider                              = alicloud.cn
  cen_id                                = alicloud_cen_instance.default.id
  transit_router_id                     = alicloud_cen_transit_router.cn.transit_router_id
  peer_transit_router_region_id         = "us-east-1"
  peer_transit_router_id                = alicloud_cen_transit_router.us.transit_router_id
  cen_bandwidth_package_id              = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  bandwidth                             = 5
  transit_router_attachment_description = var.name
  transit_router_attachment_name        = var.name
}

