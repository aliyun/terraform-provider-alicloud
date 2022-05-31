
resource "alicloud_vpc" "example" {
  vpc_name   = "test"
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_vpc_dhcp_options_set" "example" {
  dhcp_options_set_name        = "example_value"
  dhcp_options_set_description = "example_value"
  domain_name                  = "example.com"
  domain_name_servers          = "100.100.2.136"
}
resource "alicloud_vpc_dhcp_options_set_attachment" "example" {
  vpc_id              = alicloud_vpc.example.id
  dhcp_options_set_id = alicloud_vpc_dhcp_options_set.example.id
}

