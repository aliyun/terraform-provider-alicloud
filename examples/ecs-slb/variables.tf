variable "number" {
  default = "1"
}

variable "count_format" {
  default = "%02d"
}

variable "image_id" {
  default = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
}

variable "role" {
  default = "worder"
}

variable "datacenter" {
  default = "beijing"
}

variable "short_name" {
  default = "hi"
}

variable "ecs_type" {
  default = "ecs.n4.small"
}

variable "ecs_password" {
  default = "Test12345"
}

variable "availability_zones" {
  default = "cn-beijing-b"
}

variable "ssh_username" {
  default = "root"
}

variable "internet_charge_type" {
  default = "PayByTraffic"
}

variable "slb_internet_charge_type" {
  default = "paybytraffic"
}

variable "internet_max_bandwidth_out" {
  default = 5
}

variable "slb_name" {
  default = "slb_worder"
}

variable "internet" {
  default = true
}

