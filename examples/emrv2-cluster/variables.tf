# VPC variables
variable "vpc_name" {
  description = "The vpc name used to launch a new vpc when 'vpc_id' is not specified."
  default     = "TF-VPC"
}

variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "172.16.0.0/12"
}

# VSwitch variables
variable "vswitch_name" {
  description = "The vswitch name used to launch a new vswitch when 'vswitch_id' is not specified."
  default     = "TF_VSwitch"
}

variable "vswitch_cidr" {
  description = "The cidr block used to launch a new vswitch when 'vswitch_id' is not specified."
  default     = "172.16.0.0/21"
}

variable "key_pair_name" {
  description = "The name of the key pair of the ecs instance"
  default     = "terraform-kp"
}

variable "security_group_name" {
  default = "TF_SECURITY_GROUP"
}

variable "ram_name" {
  description = "The ram role name used to defined emr ecs role"
  default     = "tftest"
}

variable "cluster_name" {
  description = "The name of emr v2 cluster"
  default     = "terraform-emrv2-cluster"
}