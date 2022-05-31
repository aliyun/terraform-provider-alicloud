resource "alicloud_amqp_virtual_host" "example" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
}
resource "alicloud_amqp_exchange" "example" {
  auto_delete_state = false
  exchange_name     = "my-Exchange"
  exchange_type     = "HEADERS"
  instance_id       = alicloud_amqp_virtual_host.example.instance_id
  internal          = false
  virtual_host_name = alicloud_amqp_virtual_host.example.virtual_host_name
}
resource "alicloud_amqp_queue" "example" {
  instance_id       = alicloud_amqp_virtual_host.example.instance_id
  queue_name        = "my-Queue"
  virtual_host_name = alicloud_amqp_virtual_host.example.virtual_host_name
}
resource "alicloud_amqp_binding" "example" {
  argument          = "x-match:all"
  binding_key       = alicloud_amqp_queue.example.queue_name
  binding_type      = "QUEUE"
  destination_name  = "binding-queue"
  instance_id       = alicloud_amqp_exchange.example.instance_id
  source_exchange   = alicloud_amqp_exchange.example.exchange_name
  virtual_host_name = alicloud_amqp_exchange.example.virtual_host_name
}
