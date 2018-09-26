
data "alicloud_mns_queues" "queues" {
  name_prefix = "${alicloud_mns_queue.queue.name}"
}

resource "alicloud_mns_queue" "queue"{
	name="${var.name}"
	delay_seconds=${var.delay_seconds}
	maximum_message_size=${var.maximum_message_size}
	message_retention_period=${var.message_retention_period}
	visibility_timeout=${var.visibility_timeout}
	polling_wait_seconds=${var.polling_wait_seconds}
}