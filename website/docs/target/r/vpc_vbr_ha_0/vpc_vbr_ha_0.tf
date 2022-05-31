data "alicloud_express_connect_physical_connections" "example" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "example" {
  count                      = 2
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.example.connections[count.index].id
  virtual_border_router_name = var.name
  vlan_id                    = 100
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = "example_value"
  description       = "example_value"
}

resource "alicloud_cen_instance_attachment" "example" {
  count                    = 2
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_express_connect_virtual_border_router.example[count.index].id
  child_instance_type      = "VBR"
  child_instance_region_id = "cn-hangzhou"
}

resource "alicloud_vpc_vbr_ha" "example" {
  depends_on  = ["alicloud_cen_instance_attachment.example"]
  vbr_id      = alicloud_cen_instance_attachment.example[0].id
  peer_vbr_id = alicloud_cen_instance_attachment.example[1].id
  vbr_ha_name = "example_value"
  description = "example_value"
}

