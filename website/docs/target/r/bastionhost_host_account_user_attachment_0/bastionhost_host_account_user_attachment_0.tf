resource "alicloud_bastionhost_host" "default" {
  instance_id          = "bastionhost-cn-tl32bh0no30"
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
resource "alicloud_bastionhost_user" "default" {
  instance_id         = alicloud_bastionhost_host.default.instance_id
  mobile_country_code = "CN"
  mobile              = "13312345678"
  password            = "YourPassword-123"
  source              = "Local"
  user_name           = "my-local-user"
}

resource "alicloud_bastionhost_host_account_user_attachment" "default" {
  instance_id      = alicloud_bastionhost_host.default.instance_id
  user_id          = alicloud_bastionhost_user.default.user_id
  host_id          = alicloud_bastionhost_host.default.host_id
  host_account_ids = alicloud_bastionhost_host_account.default.*.host_account_id
}
