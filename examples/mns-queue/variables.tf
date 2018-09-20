variable "name" {
  description = "Two queues on a single account in the same region cannot have the same name. A queue name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters."
  default = "tf-example-mnsqueue"
}

variable "delay_seconds" {
  description = "This attribute defines the length of time, in seconds, after which every message sent to the queue is dequeued. Valid value range: 0-604800 seconds, i.e., 0 to 7 days."
  default = 0
}

variable "maximum_message_size" {
  description = "This indicates the maximum length, in bytes, of any message body sent to the queue. Valid value range: 1024-65536, i.e., 1K to 64K."
  default = 65536
}

variable "message_retention_period" {
  description = "Messages are deleted from the queue after a specified length of time, whether they have been activated or not. This attribute defines the viability period, in seconds, for every message in the queue. Valid value range: 60-259200 seconds, i.e., 1 minutes to 3 days."
  default = 259200
}

variable "visibility_timeout" {
  description = "Dequeued messages change from active (visible) status to inactive (invisible) status. This attribute defines the length of time, in seconds, that messages remain invisible. Messages return to active status after the set period. Valid value range: 1-43200 seconds, i.e., 1 seconds to 12 hours."
  default = 30
}
variable "polling_wait_seconds" {
  description = "Long polling is measured in seconds. When this attribute is set to 0, long polling is disabled. When it is not set to 0, long polling is enabled and message dequeue requests will be processed only when valid messages are received or when long polling times out. The value range is 0-30 seconds."
  default = 0
}