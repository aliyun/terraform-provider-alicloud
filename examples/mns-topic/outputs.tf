output "topic_name"{
	description =Two topics on a single account in the same region cannot have the same name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters."
	value = "${alicloud_mns_topic.topic.name}"
}