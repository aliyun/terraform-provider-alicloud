output "instance_id" {
  description = "The ID used to identify alikafka.instance resource."
  value       = "${alicloud_alikafka_consumer_group.default.instance_id}"
}

output "consumer_id" {
  description = "Name of ALIKAFKA Consumer Group."
  value       = "${alicloud_alikafka_consumer_group.default.consumer_id}"
}