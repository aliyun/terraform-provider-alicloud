// default vpc variables
variable "vpc_cidr" {
  default = "10.1.0.0/21"
}

variable "availability_zones" {
  default = "cn-beijing-c"
}

variable "cidr_blocks" {
  default = "10.1.1.0/24"
}

// default vpn gateway variables
variable "bandwidth" {
  default = 10
}

variable "instance_type" {
  default = "ecs.n4.small"
}

variable "internet_charge_type" {
  default = "PayByTraffic"
}

variable "instance_charge_type" {
  default = "PostPaid"
}

// default vpn varibles
variable "weight" {
  default = 100
}

variable "publish_vpc" {
  default = false
}

variable "route_dest" {
  default = "10.0.0.0/24"
}
