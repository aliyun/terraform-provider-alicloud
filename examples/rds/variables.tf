variable "engine" {
  default = "MySQL"
}

variable "engine_version" {
  default = "5.6"
}

variable "instance_class" {
  default = "rds.mysql.t1.small"
}

variable "storage" {
  default = "10"
}

variable "net_type" {
  default = "Intranet"
}

variable "user_name" {
  default = "tf_tester"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "bookstore"
}

variable "database_character" {
  default = "utf8"
}

variable "availability_zone" {
  description = "The available zone to launch ecs instance and other resources."
  default     = ""
}

# VPC variables
variable "vpc_id" {
  description = "The vpc id used to launch vswitch, security group and instance."
  default     = ""
}

variable "vpc_name" {
  description = "The vpc name used to launch a new vpc when 'vpc_id' is not specified."
  default     = "TF-VPC"
}

variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "172.16.0.0/12"
}

# VSwitch variables
variable "vswitch_id" {
  description = "The vswitch id used to launch one or more instances."
  default     = ""
}

variable "vswitch_name" {
  description = "The vswitch name used to launch a new vswitch when 'vswitch_id' is not specified."
  default     = "TF_VSwitch"
}

variable "vswitch_cidr" {
  description = "The cidr block used to launch a new vswitch when 'vswitch_id' is not specified."
  default     = "172.16.0.0/16"
}

variable "role" {
  default = "worder"
}

