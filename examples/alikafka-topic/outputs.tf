output "instance_id" {
  description = "InstanceId of your Kafka resource, the topic will create in this instance."
  value       = "${alicloud_alikafka_topic.default.instance_id}"
}

output "topic" {
  description = "Name of ALIKAFKA topic."
  value       = "${alicloud_alikafka_topic.default.topic}"
}

output "local_topic" {
  description = "Whether the topic is localTopic or not."
  value       = "${alicloud_alikafka_topic.default.local_topic}"
}

output "compact_topic" {
  description = "Whether the topic is compactTopic or not."
  value       = "${alicloud_alikafka_topic.default.compact_topic}"
}

output "partition_num" {
  description = "Partition number of ALIKAFKA topic."
  value       = "${alicloud_alikafka_topic.default.partition_num}"
}

output "remark" {
  description = "This attribute is a concise description of topic."
  value       = "${alicloud_alikafka_topic.default.remark}"
}