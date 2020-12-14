output "instance_id" {
  description = "Id of your Kafka instance."
  value       = "${alicloud_alikafka_instance.default.id}"
}

output "name" {
  description = "Name of your Kafka instance."
  value       = "${alicloud_alikafka_instance.default.name}"
}

output "topic_quota" {
  description = "The max num of topic can be create of the instance."
  value       = "${alicloud_alikafka_instance.default.topic_quota}"
}

output "disk_type" {
  description = "The disk type of the instance."
  value       = "${alicloud_alikafka_instance.default.disk_type}"
}

output "disk_size" {
  description = "The disk size of the instance."
  value       = "${alicloud_alikafka_instance.default.disk_size}"
}

output "deploy_type" {
  description = "The deploy type of the instance."
  value       = "${alicloud_alikafka_instance.default.deploy_type}"
}

output "io_max" {
  description = "The peak value of io of the instance."
  value       = "${alicloud_alikafka_instance.default.io_max}"
}

output "eip_max" {
  description = "The peak bandwidth of the instance."
  value       = "${alicloud_alikafka_instance.default.eip_max}"
}

output "paid_type" {
  description = "The paid type of the instance."
  value       = "${alicloud_alikafka_instance.default.paid_type}"
}

output "spec_type" {
  description = "The spec of the instance."
  value       = "${alicloud_alikafka_instance.default.spec_type}"
}

output "vpc_id" {
  description = "The ID of attaching VPC to instance."
  value       = "${alicloud_alikafka_instance.default.vpc_id}"
}

output "vswitch_id" {
  description = "The ID of attaching vswitch to instance."
  value       = "${alicloud_alikafka_instance.default.vswitch_id}"
}

output "zone_id" {
  description = "The Zone to launch the kafka instance."
  value       = "${alicloud_alikafka_instance.default.zone_id}"
}

output "end_point" {
  description = "The EndPoint to access the kafka instance."
  value       = "${alicloud_alikafka_instance.default.end_point}"
}