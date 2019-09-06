variable "name" {
  description = "Name of ONS Instance. Two instances on a single account in the same region cannot have the same name and the number of instances in the same region cannot exceed 8. The length must be 3 to 64 characters. Chinese characters, English letters and digits are allowed."
}

variable "instance_remark" {
  description = "This attribute is a concise description of instance. The length cannot exceed 128."
  default     = "tf-example-ons-instance-remark"
}

variable "topic" {
  description = "Name of ONS Topic. Two topics on a single instance cannot have the same name and the name cannot start with 'GID' or 'CID'. The length cannot exceed 64 characters."
}

variable "message_type" {
  description = "The type of the message. Valid values: 0: normal message; 1: partitionally ordered message; 2: globally ordered message; 4: transactional message; 5: scheduled/delayed message."
  default     = "0"
}

variable "topic_remark" {
  description = "This attribute is a concise description of topic. The length cannot exceed 128."
  default     = "tf-example-ons-topic-remark"
}

