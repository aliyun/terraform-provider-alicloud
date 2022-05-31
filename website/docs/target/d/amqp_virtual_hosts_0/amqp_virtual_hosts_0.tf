data "alicloud_amqp_virtual_hosts" "ids" {
  instance_id = "amqp-abc12345"
  ids         = ["my-VirtualHost-1", "my-VirtualHost-2"]
}
output "amqp_virtual_host_id_1" {
  value = data.alicloud_amqp_virtual_hosts.ids.hosts.0.id
}

data "alicloud_amqp_virtual_hosts" "nameRegex" {
  instance_id = "amqp-abc12345"
  name_regex  = "^my-VirtualHost"
}
output "amqp_virtual_host_id_2" {
  value = data.alicloud_amqp_virtual_hosts.nameRegex.hosts.0.id
}

