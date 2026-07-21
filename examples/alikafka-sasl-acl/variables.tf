variable "instance_id" {
  description = "InstanceId of your Kafka resource, the sasl acl will create in this instance."
}

variable "username" {
  description = "Username of ALIKAFKA sasl acl. The length should between 1 to 64 characters."
}

variable "acl_resource_type" {
  description = "Resource type of ALIKAFKA sasl acl. The resource type can only be \"Topic\" and \"Group\"."
}

variable "acl_resource_name" {
  description = "Resource name of ALIKAFKA sasl acl. The resource name should be a topic or consumer group name."
}

variable "acl_resource_pattern_type" {
  description = "Resource pattern type of ALIKAFKA sasl acl. The resource pattern support two types \"LITERAL\" and \"PREFIXED\". \"LITERAL\": A literal name defines the full name of a resource. The special wildcard character \"*\" can be used to represent a resource with any name. \"PREFIXED\": A prefixed name defines a prefix for a resource."
}

variable "acl_operation_type" {
  description = "Acl operation type of ALIKAFKA sasl acl. The operation type can only be \"Write\" and \"Read\"."
}
