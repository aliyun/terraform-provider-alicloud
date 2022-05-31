variable "name" {
  default = "snat-entry-example-name"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id     = alicloud_vpc.vpc.id
  cidr_block = "172.16.0.0/21"
  zone_id    = data.alicloud_zones.default.zones[0].id
  name       = var.name
}

resource "alicloud_nat_gateway" "default" {
  vpc_id        = alicloud_vswitch.vswitch.vpc_id
  specification = "Small"
  name          = var.name
}

resource "alicloud_eip_address" "default" {
  count        = 2
  address_name = var.name
}

resource "alicloud_eip_association" "default" {
  count         = 2
  allocation_id = element(alicloud_eip_address.default.*.id, count.index)
  instance_id   = alicloud_nat_gateway.default.id
}

resource "alicloud_common_bandwidth_package" "default" {
  name                 = "tf_cbp"
  bandwidth            = 10
  internet_charge_type = "PayByTraffic"
  ratio                = 100
}

resource "alicloud_common_bandwidth_package_attachment" "default" {
  count                = 2
  bandwidth_package_id = alicloud_common_bandwidth_package.default.id
  instance_id          = element(alicloud_eip_address.default.*.id, count.index)
}

resource "alicloud_snat_entry" "default" {
  depends_on        = [alicloud_eip_association.default]
  snat_table_id     = alicloud_nat_gateway.default.snat_table_ids
  source_vswitch_id = alicloud_vswitch.vswitch.id
  snat_ip           = join(",", alicloud_eip_address.default.*.ip_address)
}
