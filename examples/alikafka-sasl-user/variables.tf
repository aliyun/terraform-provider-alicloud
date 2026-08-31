variable "instance_id" {
  description = "InstanceId of your Kafka resource, the sasl user will create in this instance."
}

variable "username" {
  description = "Username of ALIKAFKA sasl user. The length should between 1 to 64 characters."
}

variable "password" {
  description = "Password of ALIKAFKA sasl user. The length should between 1 to 64 characters."
}