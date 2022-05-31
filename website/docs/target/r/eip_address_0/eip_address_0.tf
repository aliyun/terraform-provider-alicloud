resource "alicloud_eip_address" "example" {
  address_name         = "tf-testAcc1234"
  isp                  = "BGP"
  internet_charge_type = "PayByBandwidth"
  payment_type         = "PayAsYouGo"
}

