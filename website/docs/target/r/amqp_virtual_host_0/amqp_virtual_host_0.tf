resource "alicloud_amqp_virtual_host" "example" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
}

