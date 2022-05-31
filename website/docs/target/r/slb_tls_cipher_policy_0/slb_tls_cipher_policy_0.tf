resource "alicloud_slb_tls_cipher_policy" "example" {
  tls_cipher_policy_name = "Test-example_value"
  tls_versions           = ["TLSv1.2"]
  ciphers                = ["AES256-SHA256", "AES128-GCM-SHA256"]
}
