resource "alicloud_vpc_dhcp_options_set" "example" {
  dhcp_options_set_name        = "example_value"
  dhcp_options_set_description = "example_value"
  domain_name                  = "example.com"
  domain_name_servers          = "100.100.2.136"
}

