# Create a new Domain Record
resource "alicloud_alidns_record" "record" {
  domain_name = "domainname"
  rr          = "@"
  type        = "A"
  value       = "192.168.99.99"
  remark      = "Test new alidns record."
  status      = "ENABLE"
}
