variable "instance_id" {
  description = "InstanceId of your Kafka resource, the consumer group will create in this instance."
}
variable "consumer_id" {
  description = "Id of ALIKAFKA consumer group. The length can't exceed 64 characters."
}