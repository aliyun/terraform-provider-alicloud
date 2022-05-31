resource "alicloud_vpc_traffic_mirror_filter" "example" {
  traffic_mirror_filter_name = "example_value"
}

resource "alicloud_vpc_traffic_mirror_filter_egress_rule" "example" {
  traffic_mirror_filter_id = alicloud_vpc_traffic_mirror_filter.example.id
  priority                 = "1"
  rule_action              = "accept"
  protocol                 = "UDP"
  destination_cidr_block   = "10.0.0.0/24"
  source_cidr_block        = "10.0.0.0/24"
  destination_port_range   = "1/120"
  source_port_range        = "1/120"
}

