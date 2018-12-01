output "slb_server_certificate_id" {
  value = "${alicloud_slb_server_certificate.foo.id}"
}

output "slb_server_certificate_name" {
  value = "${alicloud_slb_server_certificate.foo.name}"
}

output "slb_server_certificate" {
  value = "${alicloud_slb_server_certificate.foo.server_certificate}"
}

output "slb_server_certificate_private_key" {
  value = "${alicloud_slb_server_certificate.foo.private_key}"
}

output "slb_server_certificate_file" {
  value = "${path.module}/server_certificate.pem"
}

output "slb_server_certificate_private_key_file" {
  value = "${path.module}/private_key.pem"
}
