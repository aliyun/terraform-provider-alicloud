resource "alicloud_vpc" "default" {
  name       = "tf_vpc_test"
  cidr_block = "${var.vpc_cidr}"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "${var.cidr_blocks}"
  availability_zone = "${var.availability_zones}"
}

resource "alicloud_vpn_gateway" "default" {
  name                 = "tf_vpn_gateway_test"
  vpc_id               = "${alicloud_vpc.default.id}"
  bandwidth            = "${var.bandwidth}"
  instance_charge_type = "${var.instance_charge_type}"
  enable_ssl           = false
}


resource "alicloud_vpn_connection" "default" {
  name                = "tf_vpn_connection_test"
  customer_gateway_id = "${alicloud_vpn_customer_gateway.default.id}"
  vpn_gateway_id      = "${alicloud_vpn_gateway.default.id}"
  local_subnet        = ["192.168.2.0/24"]
  remote_subnet       = ["192.168.3.0/24"]
}

resource "alicloud_vpn_customer_gateway" "default" {
  name       = "tf_customer_gateway_test"
  ip_address = "192.168.1.1"
}

resource "alicloud_vpn_route_entry" "default" {
  description    = "tf_vpn_route_entry_test"
  vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
  route_dest     = "${var.route_dest}"
  next_hop       = "${alicloud_vpn_connection.default.id}"
  weight         = "${var.weight}"
  publish_vpc    = "${var.publish_vpc}"
}
