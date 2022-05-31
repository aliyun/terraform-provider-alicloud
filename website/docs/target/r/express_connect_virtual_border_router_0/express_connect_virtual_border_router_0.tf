data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^my-PhysicalConnection"
}

resource "alicloud_express_connect_virtual_border_router" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = "example_value"
  vlan_id                    = 1
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

