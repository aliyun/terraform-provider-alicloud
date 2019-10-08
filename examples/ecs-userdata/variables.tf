variable "vpc_cidr" {
  default = "172.16.0.0/12"
}

variable "vswitch_cidr" {
  default = "172.16.0.0/21"
}

variable "zone" {
  default = "cn-beijing-a"
}

variable "password" {
  default = "Test123456"
}

variable "image" {
  default = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
}

variable "ecs_type" {
  default = "ecs.n4.large"
}

