data "alicloud_bastionhost_hosts" "ids" {
  instance_id = "example_value"
  ids         = ["1", "2"]
}
output "bastionhost_host_id_1" {
  value = data.alicloud_bastionhost_hosts.ids.hosts.0.id
}

data "alicloud_bastionhost_hosts" "nameRegex" {
  instance_id = "example_value"
  name_regex  = "^my-Host"
}
output "bastionhost_host_id_2" {
  value = data.alicloud_bastionhost_hosts.nameRegex.hosts.0.id
}

