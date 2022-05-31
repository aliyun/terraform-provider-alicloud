
resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block             = "172.16.0.0/12"
  desktop_access_type    = "Internet"
  office_site_name       = "your_office_site_name"
  enable_internet_access = false
}

resource "alicloud_ecd_nas_file_system" "example" {
  nas_file_system_name = "example_value"
  office_site_id       = alicloud_ecd_simple_office_site.default.id
  description          = "example_value"
}

