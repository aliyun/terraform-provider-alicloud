resource "alicloud_ddoscoo_domain_resource" "example" {
  domain       = "tftestacc1234.abc"
  rs_type      = 0
  instance_ids = ["ddoscoo-cn-6ja1rl4j****"]
  real_servers = ["177.167.32.11"]
  https_ext    = "{\"Http2\":1,\"Http2https\":0，\"Https2http\":0}"
  proxy_types {
    proxy_ports = [443]
    proxy_type  = "https"
  }
}

