resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  bandwidth           = 5
  desktop_access_type = "Internet"
  office_site_name    = "site_name"
}

