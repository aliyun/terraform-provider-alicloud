resource "alicloud_slb_ca_certificate" "foo-file" {
  ca_certificate_name = "tf-testAccSlbCACertificate"
  ca_certificate      = file("${path.module}/ca_certificate.pem")
}
