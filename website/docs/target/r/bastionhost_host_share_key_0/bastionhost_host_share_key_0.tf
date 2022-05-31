data "alicloud_bastionhost_instances" "default" {
}

resource "alicloud_bastionhost_host_share_key" "default" {
  host_share_key_name = "example_name"
  instance_id         = data.alicloud_bastionhost_instances.default.instances.0.id
  pass_phrase         = "example_value"
  private_key         = "example_value"
}
