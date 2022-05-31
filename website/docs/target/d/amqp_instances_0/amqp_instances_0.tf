data "alicloud_amqp_instances" "ids" {
  ids = ["amqp-abc12345", "amqp-abc34567"]
}
output "amqp_instance_id_1" {
  value = data.alicloud_amqp_instances.ids.instances.0.id
}

data "alicloud_amqp_instances" "nameRegex" {
  name_regex = "^my-Instance"
}
output "amqp_instance_id_2" {
  value = data.alicloud_amqp_instances.nameRegex.instances.0.id
}

