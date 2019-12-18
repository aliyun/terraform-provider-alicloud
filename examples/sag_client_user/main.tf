resource "alicloud_sag_client_user" "default" {
  sag_id    = "sag-xxxxx"
  bandwidth = "20"
  user_mail = "tftest-xxxxx@test.com"
  user_name = "th-username-xxxxx"
  password  = "xxxxxxx"
  client_ip = "192.1.10.0"
}