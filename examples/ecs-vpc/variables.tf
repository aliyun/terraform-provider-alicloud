variable "number" {
  default = "3"
}

variable "count_format" {
  default = "%02d"
}

variable "image_id" {
  default = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
}

variable "role" {
  default = "example-ecs-vpc"
}

variable "short_name" {
  default = "hi"
}

variable "ecs_password" {
  default = "Test12345"
}

variable "ecs_type" {
  default = "ecs.n4.small"
}

variable "ssh_username" {
  default = "root"
}

//if instance_charge_type is "PrePaid", then must be set period, the value is 1 to 30, unit is month
variable "instance_charge_type" {
  default = "PostPaid"
}

variable "system_disk_category" {
  default = "cloud_efficiency"
}

variable "internet_charge_type" {
  default = "PayByTraffic"
}

variable "internet_max_bandwidth_out" {
  default = 5
}

variable "disk_category" {
  default = "cloud_ssd"
}

variable "disk_size" {
  default = "40"
}

variable "vpc_name" {
  description = "The vpc name used to launch a new vpc."
  default     = "TF-VPC"
}

variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc."
  default     = "172.16.0.0/12"
}

variable "vswitch_id" {
  description = "The vswitch id of existing vswitch."
  default     = ""
}

variable "vswitch_name" {
  description = "The vswitch name used to launch a new vswitch when vswitch_id is not set."
  default     = "TF_VSwitch"
}

variable "vswitch_cidr" {
  description = "The cidr block used to launch a new vswitch when vswitch_id is not set."
  default     = "172.16.0.0/16"
}

# Security Group variables
variable "sg_id" {
  description = "The security group id of existing security group."
  default     = ""
}

variable "sg_name" {
  description = "The security group name used to launch a new security group when sg is not set."
  default     = "TF_Security_Group"
}
