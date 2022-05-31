data "alicloud_amqp_queues" "ids" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
  ids               = ["my-Queue-1", "my-Queue-2"]
}
output "amqp_queue_id_1" {
  value = data.alicloud_amqp_queues.ids.queues.0.id
}

data "alicloud_amqp_queues" "nameRegex" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
  name_regex        = "^my-Queue"
}
output "amqp_queue_id_2" {
  value = data.alicloud_amqp_queues.nameRegex.queues.0.id
}

