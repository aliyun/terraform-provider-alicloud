variable "name" {
  default = "tfacc_host_account_share_key_attachment"
}

data "alicloud_bastionhost_instances" "default" {
}

resource "alicloud_bastionhost_host_share_key" "default" {
  host_share_key_name = "example_name"
  instance_id         = data.alicloud_bastionhost_instances.default.instances.0.id
  pass_phrase         = "example_value"
  private_key         = "example_value"
}

resource "alicloud_bastionhost_host" "default" {
  instance_id          = data.alicloud_bastionhost_instances.default.ids.0
  host_name            = var.name
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  os_type              = "Linux"
  source               = "Local"
}

resource "alicloud_bastionhost_host_account" "default" {
  instance_id       = data.alicloud_bastionhost_instances.default.ids.0
  host_account_name = var.name
  host_id           = alicloud_bastionhost_host.default.host_id
  protocol_name     = "SSH"
  password          = "YourPassword12345"
}

resource "alicloud_bastionhost_host_account_share_key_attachment" "default" {
  instance_id       = data.alicloud_bastionhost_instances.default.instances.0.id
  host_share_key_id = alicloud_bastionhost_host_share_key.default.host_share_key_id
  host_account_id   = alicloud_bastionhost_host_account.default.host_account_id
}
