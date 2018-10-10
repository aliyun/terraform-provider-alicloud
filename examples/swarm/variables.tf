variable "most_recent" {
  default = true
}

variable "image_owners" {
  default = ""
}

variable "name_regex" {
  default = "^centos_6\\w{1,5}[64].*"
}

variable "cpu_core_count" {
  default = 1
}

variable "memory_size" {
  default = 2
}

variable "vpc_name" {
  default = "alicloud_vpc"
}

variable "vpc_cidr" {
  default = "10.1.0.0/21"
}

variable "vswitch_cidr" {
  default = "10.1.1.0/24"
}

variable "password" {
  default = "Test12345"
}

variable "cidr_block" {
  default = "172.20.0.0/24"
}

variable "cluster_name" {
  default = "Alicloud-cluster"
}

variable "node_number" {
  default = 1
}

variable "disk_category" {
  default = "cloud_efficiency"
}

variable "disk_size" {
  default = 40
}
