output "name" {
  description = "Name of ONS Instance."
  value       = "${alicloud_ons_instance.instance.name}"
}

output "instance_id" {
  description = "The ID used to identify ons.instance resource."
  value       = "${alicloud_ons_instance.instance.id}"
}

output "topic" {
  description = "Name of ONS Topic."
  value       = "${alicloud_ons_topic.default.topic}"
}

output "topic_remark" {
  description = "This attribute is a concise description of topic."
  value       = "${alicloud_ons_topic.default.remark}"
}

