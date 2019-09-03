resource "alicloud_cdn_domain" "domain" {
  domain_name = "${var.domain_name}"
  cdn_type    = "${var.cdn_type}"
  source_type = "${var.source_type}"
  sources     = "${var.sources}"

  // configs
  optimize_enable      = "${var.enable}"
  page_compress_enable = "${var.enable}"
  range_enable         = "${var.enable}"
  video_seek_enable    = "${var.enable}"
  block_ips            = "${var.block_ips}"

  parameter_filter_config {
      enable        = "${var.enable}"
      hash_key_args = "${var.hash_key_args}"
  }

  page_404_config {
      page_type       = "${var.page_type}"
      custom_page_url = "http://${var.domain_name}/notfound/"
  }

  refer_config {
      refer_type  = "${var.refer_type}"
      refer_list  = "${var.refer_list}"
      allow_empty = "${var.enable}"
    }

  auth_config {
      auth_type  = "${var.auth_type}"
      master_key = "helloworld1"
      slave_key  = "helloworld2"
  }

  http_header_config {
      header_key   = "Content-Type"
      header_value = "text/plain"
  }

  http_header_config {
      header_key   = "Access-Control-Allow-Origin"
      header_value = "*"
  }

  cache_config {
      cache_content = "/hello/world"
      ttl           = 1000
      cache_type    = "path"
  }
  cache_config {
      cache_content = "/hello/world/youyou"
      ttl           = 1000
      cache_type    = "path"
  }
  cache_config {
      cache_content = "txt,jpg,png"
      ttl           = 2000
      cache_type    = "suffix"
  }
}
