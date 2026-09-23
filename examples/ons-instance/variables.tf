variable "name" {
  description = "Two instances on a single account in the same region cannot have the same name and the number of instances in the same region cannot exceed 8. The length must be 3 to 64 characrers. Chinese characters, English letters and digits are allowed."
}

variable "remark" {
  description = "This attribute is a concise description of instance. The length cannot exceed 128."
  default     = "tf-example-ons-instance-remark"
}

