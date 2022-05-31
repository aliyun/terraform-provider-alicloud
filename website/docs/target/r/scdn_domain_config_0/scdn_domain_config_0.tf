# Create a new Domain config.
resource "alicloud_scdn_domain" "domain" {
  domain_name = "mydomain.xiaozhu.com"
  cdn_type    = "web"
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    type     = "ipaddr"
    priority = "20"
    port     = 80
  }
}

resource "alicloud_scdn_domain_config" "config" {
  domain_name   = alicloud_scdn_domain.domain.domain_name
  function_name = "ip_allow_list_set"
  function_args {
    arg_name  = "ip_list"
    arg_value = "110.110.110.110"
  }
}
