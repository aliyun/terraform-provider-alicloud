resource "alicloud_bastionhost_user_attachment" "example" {
  instance_id   = "bastionhost-cn-tl3xxxxxxx"
  user_group_id = "10"
  user_id       = "100"
}

