resource "alicloud_waf_certificate" "default" {
  certificate_name = "your_certificate_name"
  instance_id      = "your_instance_id"
  domain           = "your_domain_name"
  private_key      = "your_private_key"
  certificate      = "your_certificate"
}
resource "alicloud_waf_certificate" "default2" {
  instance_id    = "your_instance_id"
  domain         = "your_domain_name"
  certificate_id = "your_certificate_id"
}
