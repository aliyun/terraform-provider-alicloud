resource "alicloud_vod_domain" "default" {
  domain_name = "your_domain_name"
  scope       = "domestic"
  sources {
    source_type    = "domain"
    source_content = "your_source_content"
    source_port    = "80"
  }
  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}

