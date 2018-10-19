output "slb_ca_certificate_id" {
  value = "${alicloud_slb_ca_certificate.foo.id}"
}

output "slb_ca_certificate_name" {
  value = "${alicloud_slb_ca_certificate.foo.name}"
}

output "slb_ca_certificate" {
  value = "${alicloud_slb_ca_certificate.foo.ca_certificate}"
}
