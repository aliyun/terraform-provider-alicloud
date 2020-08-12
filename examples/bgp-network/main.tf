resource "alicloud_bgp_network" "foo" {
  dst_cidr_block = "1.1.1.0/24"
  router_id = var.route_id
}
