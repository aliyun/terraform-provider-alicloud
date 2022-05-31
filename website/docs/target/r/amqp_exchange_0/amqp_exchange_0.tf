resource "alicloud_amqp_virtual_host" "example" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
}
resource "alicloud_amqp_exchange" "example" {
  auto_delete_state = false
  exchange_name     = "my-Exchange"
  exchange_type     = "DIRECT"
  instance_id       = alicloud_amqp_virtual_host.example.instance_id
  internal          = false
  virtual_host_name = alicloud_amqp_virtual_host.example.virtual_host_name
}

