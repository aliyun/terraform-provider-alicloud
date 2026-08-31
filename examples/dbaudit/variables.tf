variable "description" {
  default = "TerraformTest"
}

variable "plan_code" {
  default = "alpha.basic"
}

variable "period" {
  default = 1
}

variable "available_zone_resource_creation" {
  default = "VSwitch"
}

variable "vpc_cidr_block" {
  default = "172.16.0.0/12"
}

variable "vswitch_cidr_block" {
  default = "172.16.0.0/21"
}


