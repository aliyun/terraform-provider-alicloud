# create a CA certificate
resource "alicloud_slb_ca_certificate" "foo" {
  ca_certificate_name = "tf-testAccSlbCACertificate"
  ca_certificate      = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJnI******90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}
