resource "alicloud_ram_security_preference" "example" {
  enable_save_mfa_ticket        = false
  allow_user_to_change_password = true
}
