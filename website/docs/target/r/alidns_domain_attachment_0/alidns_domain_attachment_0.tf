resource "alicloud_alidns_domain_attachment" "dns" {
  instance_id  = "dns-cn-mp91lyq9xxxx"
  domain_names = ["test111.abc", "test222.abc"]
}
