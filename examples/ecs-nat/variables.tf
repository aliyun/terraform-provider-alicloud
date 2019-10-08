variable "vpc_cidr" {
  default = "10.1.0.0/21"
}

variable "vswitch_cidr" {
  default = "10.1.1.0/24"
}

variable "zone" {
  default = "cn-hangzhou-i"
}

variable "image" {
  default = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
}

variable "instance_nat_type" {
  default = "ecs.n4.small"
}

variable "instance_worker_type" {
  default = "ecs.n4.large"
}

variable "instance_pwd" {
  default = "Test123456"
}

