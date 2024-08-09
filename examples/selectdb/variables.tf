variable "engine_version" {
  default = "3.0"
}

variable "db_type" {
  default = "SelectDB"
}

# VPC variables
variable "vpc_id" {
  description = "The vpc id used to launch vswitch, security group and instance."
  default     = ""
}

variable "vpc_name" {
  description = "The vpc name used to launch a new vpc when 'vpc_id' is not specified."
  default     = "TF-VPC-example"
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
  default     = "TF_VSwitch-example"
}

variable "vswitch_cidr" {
  description = "The cidr block used to launch a new vswitch when 'vswitch_id' is not specified."
  default     = "172.16.0.0/16"
}

variable "availability_zone" {
  description = "The available zone to launch ecs instance and other resources."
  default     = ""
}

