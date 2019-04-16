output "cert" {
  value = "${alicloud_cas_certificate.cert.*.id}"
}
