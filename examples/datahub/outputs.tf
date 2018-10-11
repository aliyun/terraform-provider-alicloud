output "ProjectName" {
  value = "${alicloud_datahub_project.example.name}"
}

output "ProjectCreateTime" {
  value = "${alicloud_datahub_project.example.create_time}"
}

output "TopicName" {
  value = "${alicloud_datahub_topic.example.name}"
}

output "TopicCreateTime" {
  value = "${alicloud_datahub_topic.example.create_time}"
}

output "ShardCount" {
  value = "${alicloud_datahub_topic.example.shard_count}"
}

output "SubscriptionId" {
  value = "${alicloud_datahub_subscription.example.sub_id}"
}

output "SubscriptionCreateTime" {
  value = "${alicloud_datahub_subscription.example.create_time}"
}
