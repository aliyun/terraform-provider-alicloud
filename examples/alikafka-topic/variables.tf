variable "instance_id" {
  description = "InstanceId of your Kafka resource, the topic will create in this instance."
}

variable "topic" {
  description = "Name of ALIKAFKA topic. Two topics on a single instance cannot have the same name. The length cannot exceed 64 characters."
}

variable "local_topic" {
  description = "Whether the topic is localTopic or not."
}

variable "compact_topic" {
  description = "Whether the topic is compactTopic or not. Compact topic must be a localTopic."
}

variable "partition_num" {
  description = "The number of partitions of the topic. The number should between 1 and 48."
}

variable "remark" {
  description = "This attribute is a concise description of topic. The length cannot exceed 64 characters."
  default     = "tf-example-alikafka-topic-remark"
}