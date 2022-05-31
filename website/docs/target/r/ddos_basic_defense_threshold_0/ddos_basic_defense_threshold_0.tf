resource "alicloud_ddos_basic_defense_threshold" "example" {
  instance_id   = alicloud_eip_address.default.id
  ddos_type     = "defense"
  instance_type = "eip"
  bps           = 390
  pps           = 90000
}

resource "alicloud_eip_address" "default" {
  address_name         = var.name
  isp                  = "BGP"
  internet_charge_type = "PayByBandwidth"
  payment_type         = "PayAsYouGo"
}
