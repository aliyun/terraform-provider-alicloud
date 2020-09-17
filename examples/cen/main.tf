provider "alicloud" {
  alias  = "hz"
  region = "cn-hangzhou"
}

provider "alicloud" {
  alias  = "bj"
  region = "cn-beijing"
}

data "alicloud_zones" "default" {
  provider                    = alicloud.hz
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  provider          = alicloud.hz
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  provider    = alicloud.hz
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_cen_instance" "cen" {
  name        = var.name
  description = var.description
}

resource "alicloud_vpc" "vpc_1" {
  provider   = alicloud.hz
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "vpc_2" {
  provider   = alicloud.bj
  name       = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_cen_bandwidth_package" "bwp" {
  bandwidth = 2

  geographic_region_ids = [
    "China",
    "China",
  ]
}

resource "alicloud_vswitch" "default" {
  provider          = alicloud.hz
  vpc_id            = alicloud_vpc.vpc_1.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = data.alicloud_zones.default.zones[0].id
  name              = var.name
}

resource "alicloud_security_group" "default" {
  provider    = alicloud.hz
  name        = var.name
  description = "foo"
  vpc_id      = alicloud_vpc.vpc_1.id
}

resource "alicloud_instance" "default" {
  provider                   = alicloud.hz
  vswitch_id                 = alicloud_vswitch.default.id
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = [alicloud_security_group.default.id]
  instance_name              = var.name
}

resource "alicloud_cen_bandwidth_package_attachment" "bwp_attach" {
  instance_id          = alicloud_cen_instance.cen.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.bwp.id
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
  instance_id              = alicloud_cen_instance.cen.id
  child_instance_id        = alicloud_vpc.vpc_1.id
  child_instance_type      = "VPC"
  child_instance_region_id = "cn-hangzhou"

  depends_on = [alicloud_vswitch.default]
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
  instance_id              = alicloud_cen_instance.cen.id
  child_instance_id        = alicloud_vpc.vpc_2.id
  child_instance_type      = "VPC"
  child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_bandwidth_limit" "limit" {
  instance_id = alicloud_cen_instance.cen.id

  region_ids = [
    "cn-hangzhou",
    "cn-beijing",
  ]

  bandwidth_limit = 2

  depends_on = [
    alicloud_cen_bandwidth_package_attachment.bwp_attach,
    alicloud_cen_instance_attachment.vpc_attach_1,
    alicloud_cen_instance_attachment.vpc_attach_2,
  ]
}

resource "alicloud_route_entry" "route_entry" {
  provider              = alicloud.hz
  route_table_id        = alicloud_vpc.vpc_1.route_table_id
  destination_cidrblock = "11.0.0.0/16"
  nexthop_type          = "Instance"
  nexthop_id            = alicloud_instance.default.id
}

resource "alicloud_cen_route_entry" "cen_route_entry" {
  provider       = alicloud.hz
  instance_id    = alicloud_cen_instance.cen.id
  route_table_id = alicloud_vpc.vpc_1.route_table_id
  cidr_block     = alicloud_route_entry.route_entry.destination_cidrblock

  depends_on = [alicloud_cen_instance_attachment.vpc_attach_1]
}

