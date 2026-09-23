variable "availability_zone" {
  default = ""
}

variable "security_group_name" {
  default = "tf-sg"
}

variable "scaling_min_size" {
  default = 1
}

variable "scaling_max_size" {
  default = 1
}

variable "enable" {
  default = true
}

variable "removal_policies" {
  type    = list(string)
  default = ["OldestInstance", "NewestInstance"]
}

variable "ecs_instance_type" {
  default = "ecs.n4.large"
}

# VPC variables
variable "vpc_id" {
  description = "The vpc id is used to launch scaling group."
  default     = ""
}

variable "vpc_name" {
  description = "The vpc name used to launch a new vpc when 'vpc_id' is not specified."
  default     = "TF-VPC-For-Scaling-Group"
}

variable "vpc_description" {
  description = "The vpc description used to launch a new vpc when 'vpc_id' is not specified."
  default     = "A new VPC created by Terrafrom example ess-scaling"
}

variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "172.16.0.0/12"
}

# VSwitch variables
variable "vswitch_id" {
  description = "The vswitch id is used to launch a scaling group."
  default     = ""
}

variable "vswitch_cidr" {
  description = "The cidr block used to launch a new vswitch when 'vswitch_id' is not specified."
  default     = "172.16.0.0/16"
}

variable "vswitch_name" {
  description = "The vswitch name prefix used to launch a new vswitch when 'vswitch_id' is not specified."
  default     = "TF-VSwitch-For-Scaling-Group"
}

variable "vswitch_description" {
  description = "The vswitch description used to launch a new vswitch when 'vswitch_id' is not specified."
  default     = "New VSwitch created by Terrafrom module tf-alicloud-vpc-cluster."
}

