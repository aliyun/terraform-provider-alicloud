variable "name" {
  default = "testAccSecurityPolicy"
}

resource "alicloud_alb_security_policy" "default" {
  security_policy_name = var.name
  tls_versions         = ["TLSv1.0"]
  ciphers              = ["ECDHE-ECDSA-AES128-SHA", "AES256-SHA"]
}

