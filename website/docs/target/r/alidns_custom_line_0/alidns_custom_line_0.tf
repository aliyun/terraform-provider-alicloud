resource "alicloud_alidns_custom_line" "default" {
  custom_line_name = "tf-testacc"
  domain_name      = "your_domain_name"
  ip_segment_list {
    start_ip = "192.0.2.123"
    end_ip   = "192.0.2.125"
  }
}
