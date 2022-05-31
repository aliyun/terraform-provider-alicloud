resource "alicloud_bastionhost_host" "example" {
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  host_name            = "example_value"
  instance_id          = "bastionhost-cn-tl3xxxxxxx"
  os_type              = "Linux"
  source               = "Local"
}

