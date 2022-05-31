resource "alicloud_amqp_virtual_host" "example" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
}
resource "alicloud_amqp_queue" "example" {
  instance_id       = alicloud_amqp_virtual_host.example.instance_id
  queue_name        = "my-Queue"
  virtual_host_name = alicloud_amqp_virtual_host.example.virtual_host_name
}

