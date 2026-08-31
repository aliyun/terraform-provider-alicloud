variable "zone_name" {
  default = "www.test.com"
}

variable "resource_record" {
  default = "www"
}

variable "type" {
  default = "A"
}

variable "value" {
  default = "1.1.1.1"
}

// Only MX supports priority
variable "priority" {
  default = "10"
}

variable "long_name" {
  default = "alicloud.com"
}

variable "vpc_cidr" {
  default = "10.1.0.0/21"
}

