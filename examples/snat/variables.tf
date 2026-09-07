variable "vpc_cidr" {
  default = "10.1.0.0/21"
}

variable "vswitch_cidr" {
  default = "10.1.1.0/24"
}

variable "rule_policy" {
  default = "accept"
}

variable "instance_type" {
  default = "ecs.n4.small"
}

variable "image_id" {
  default = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
}

variable "disk_category" {
  default = "cloud_efficiency"
}

