data "alicloud_amqp_exchanges" "ids" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
  ids               = ["my-Exchange-1", "my-Exchange-2"]
}
output "amqp_exchange_id_1" {
  value = data.alicloud_amqp_exchanges.ids.exchanges.0.id
}

data "alicloud_amqp_exchanges" "nameRegex" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
  name_regex        = "^my-Exchange"
}
output "amqp_exchange_id_2" {
  value = data.alicloud_amqp_exchanges.nameRegex.exchanges.0.id
}

