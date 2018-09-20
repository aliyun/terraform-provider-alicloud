resource "alicloud_mns_topic" "topic"{
	name="${var.name}"
	maximum_message_size=${var.maximum_message_size}
	logging_enabled=${var.logging_enabled} 
}