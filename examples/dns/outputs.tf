output "domain" {
  value = "${alicloud_dns.dns.id}"
}

output "group" {
  value = "${alicloud_dns_group.group.id}"
}

output "record" {
  value = "${alicloud_dns_record.record.id}"
}