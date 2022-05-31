data "alicloud_waf_instances" "default" {
}

resource "alicloud_waf_domain" "default" {
  domain_name       = "you domain"
  instance_id       = data.alicloud_waf_instances.default.ids.0
  is_access_product = "On"
  source_ips        = ["1.1.1.1"]
  cluster_type      = "PhysicalCluster"
  http2_port        = [443]
  http_port         = [80]
  https_port        = [443]
  http_to_user_ip   = "Off"
  https_redirect    = "Off"
  load_balancing    = "IpHash"
  log_headers {
    key   = "foo"
    value = "http"
  }
}

resource "alicloud_waf_protection_module" "default" {
  instance_id  = data.alicloud_waf_instances.default.ids.0
  domain       = alicloud_waf_domain.default.domain_name
  defense_type = "ac_cc"
  mode         = 0
  status       = 0
}
