data "alicloud_bastionhost_host_share_keys" "ids" {
  instance_id = "example_value"
  ids         = ["example_value-1", "example_value-2"]
}
output "bastionhost_host_share_key_id_1" {
  value = data.alicloud_bastionhost_host_share_keys.ids.keys.0.id
}

data "alicloud_bastionhost_host_share_keys" "nameRegex" {
  instance_id = "example_value"
  name_regex  = "^my-HostShareKey"
}
output "bastionhost_host_share_key_id_2" {
  value = data.alicloud_bastionhost_host_share_keys.nameRegex.keys.0.id
}
