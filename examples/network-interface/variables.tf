variable "vpc_cidr" {
  default = "192.168.0.0/24"
}

variable "vswitch_cidr" {
  default = "192.168.0.0/24"
}

variable "system_disk_category" {
  default = "cloud_efficiency"
}

variable "most_recent" {
  default = true
}

variable "image_owners" {
  default = "system"
}

variable "name_regex" {
  default = "^centos_6\\w{1,5}[64].*"
}

