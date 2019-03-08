variable "count" {
  default = "1"
}

variable "count_format" {
  default = "%02d"
}

variable "image_id" {
  default = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
}

variable "role" {
  default = "work"
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

variable "internet_charge_type" {
  default = "PayByTraffic"
}

variable "internet_max_bandwidth_out" {
  default = 5
}

variable "disk_category" {
  default = "cloud_efficiency"
}

variable "disk_size" {
  default = "40"
}

variable "nic_type" {
  default = "intranet"
}
