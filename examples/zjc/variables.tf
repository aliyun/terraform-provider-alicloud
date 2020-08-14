variable "router_type" {
  default = "VRouter"
}

variable "instance_charge_type" {
  default = "PostPaid"
}

variable "role" {
  type = "list"
  default = ["AcceptingSide", "InitiatingSide"]
}

variable "specification" {
  type = "list"

  default = ["Negative", "Large.2", "Large.1"]
}

variable "name" {
  default = "test-name"
}

variable "region" {
  default = "cn-hangzhou"
}
