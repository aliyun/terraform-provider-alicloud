# VPC variables
variable "vpc_id" {
  description = "The vpc id used to launch vswitch, security group and instance."
  default     = "vpc-hp32dbx98fhcjp9kr5m75"
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
  default     = "vsw-hp34pmwcw8i5was4sn07h"
}

variable "vswitch_name" {
  description = "The vswitch name used to launch a new vswitch when 'vswitch_id' is not specified."
  default     = "TF_VSwitch"
}

variable "vswitch_cidr" {
  description = "The cidr block used to launch a new vswitch when 'vswitch_id' is not specified."
  default     = "172.16.0.0/16"
}

variable "availability_zone" {
  description = "The available zone to launch ecs instance and other resources."
  default     = ""
}

variable "role" {
  default = "worder"
}

variable "security_group_name" {
  default = "TF_SECRITY_GROUP"
}

variable "security_group_id" {
  default = ""
}

variable "ram_name" {
  description = "The ram role name used to defined emr ecs role"
  default     = "testrole"
}