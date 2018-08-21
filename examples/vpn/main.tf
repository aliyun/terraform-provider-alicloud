resource "alicloud_vpn" "main" {
  name = "${var.long_name}"
  vpc_id = "${var.vpc}"
  bandwidth = "${var.bandwidth}"
  enable_ssl = true
  instance_charge_type = "${var.instance_charge_type}"
  auto_pay = "${var.auto_pay}"
}

resource "alicloud_vpn_customer_gateway" "main" {
  name = "${var.cgw_name}"
  ip_address = "${var.cgw_ip_address}"
  description = "${var.cgw_description}"
}

resource "alicloud_vpn_connection" "conn1" {
   name = "${var.vco_name}"
   vpn_gateway_id = "${alicloud_vpn.main.id}"
   customer_gateway_id = "${alicloud_vpn_customer_gateway.main.id}"
   local_subnet = "${var.local_subnet}"
   remote_subnet = "${var.remote_subnet}"
   effect_immediately = "${var.effect_immediately}"
   ike_config = "${var.ike_config}"
   ipsec_config = "${var.ipsec_config}"
}

resource "alicloud_ssl_vpn_server" "ssl_vpn_server_1" {
    name = "${var.ssl_server_name}"
    vpn_gateway_id = "${alicloud_vpn.main.id}"
    client_ip_pool = "${var.ssl_vpn_client_pool_1}"
    local_subnet = "${var.ssl_vpn_local_subnet_1}"
    proto = "${var.ssl_vpn_proto}"
    cipher = "${var.ssl_vpn_cipher}"
    port = "${var.ssl_vpn_port}"
    compress = "${var.ssl_vpn_compress}"
}

resource "alicloud_ssl_vpn_client_cert" "ssl_client_cert_1" {
    ssl_vpn_server_id = "${alicloud_ssl_vpn_server.ssl_vpn_server_1.id}"
    name = "${var.ssl_vpn_client_cert_1}"
}
