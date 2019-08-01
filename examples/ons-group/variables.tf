variable "name" {
  description = "Name of ONS Instance. Two instances on a single account in the same region cannot have the same name and the number of instances in the same region cannot exceed 8. The length must be 3 to 64 characters. Chinese characters, English letters and digits are allowed."
}

variable "instance_remark" {
  description = "This attribute is a concise description of instance. The length cannot exceed 128."
  default     = "tf-example-ons-instance-remark"
}

variable "group_id" {
  description = "Name of ONS Group. Two groups on a single instance cannot have the same name and the name must start with 'GID-' or 'GID_'. The length must be 7 to 64 characters."
}

variable "group_remark" {
  description = "This attribute is a concise description of group. The length cannot exceed 256."
  default     = "tf-example-ons-group-remark"
}