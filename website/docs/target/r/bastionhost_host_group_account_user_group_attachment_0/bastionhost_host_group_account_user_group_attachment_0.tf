resource "alicloud_bastionhost_host" "default" {
  instance_id          = "bastionhost-cn-tl3xxxxxxx"
  host_name            = var.name
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  os_type              = "Linux"
  source               = "Local"
}
resource "alicloud_bastionhost_host_account" "default" {
  count             = 3
  instance_id       = alicloud_bastionhost_host.default.instance_id
  host_account_name = "example_value-${count.index}"
  host_id           = alicloud_bastionhost_host.default.host_id
  protocol_name     = "SSH"
  password          = "YourPassword12345"
}
resource "alicloud_bastionhost_user_group" "default" {
  instance_id     = alicloud_bastionhost_host.default.instance_id
  user_group_name = "my-local-user"
}
resource "alicloud_bastionhost_host_group" "default" {
  host_group_name = "example_value"
  instance_id     = "bastionhost-cn-tl3xxxxxxx"
}
resource "alicloud_bastionhost_host_group_account_user_group_attachment" "default" {
  instance_id        = alicloud_bastionhost_host.default.instance_id
  user_group_id      = alicloud_bastionhost_user_group.default.user_group_id
  host_group_id      = alicloud_bastionhost_host_group.default.host_group_id
  host_account_names = alicloud_bastionhost_host_account.default.*.host_account_name
}
