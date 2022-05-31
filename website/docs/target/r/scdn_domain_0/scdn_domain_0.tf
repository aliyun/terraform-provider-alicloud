resource "alicloud_scdn_domain" "example" {
  domain_name = "my-Domain"
  sources {
    content  = "xxx.aliyuncs.com"
    enabled  = "online"
    port     = 80
    priority = "20"
    type     = "oss"
  }
}

