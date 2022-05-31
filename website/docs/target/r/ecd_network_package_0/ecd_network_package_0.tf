
resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = "your_office_site_name"
}

resource "alicloud_ecd_network_package" "example" {
  bandwidth      = 10
  office_site_id = alicloud_ecd_simple_office_site.default.id
}

