variable "name" {
  description = "Name of your Kafka instance. The length cannot exceed 64 characters."
}

variable "topic_quota" {
  description = "The max num of topic can be create of the instance. When modify this value, it only adjust to a greater value."
}

variable "disk_type" {
  description = "The disk type of the instance. 0: efficient cloud disk , 1: SSD."
}

variable "disk_size" {
  description = "The disk size of the instance. When modify this value, it only adjust to a greater value."
}

variable "deploy_type" {
  description = "The deploy type of the instance. Currently only support two deploy type, 4: eip/vpc instance, 5: vpc instance."
}

variable "io_max" {
  description = "The peak value of io of the instance. When modify this value, it only adjust to a greater value."
}

variable "eip_max" {
  description = "The peak bandwidth of the instance. When modify this value, it only adjust to a greater value."
}

variable "vswitch_id" {
  description = "The ID of attaching vswitch to instance."
}