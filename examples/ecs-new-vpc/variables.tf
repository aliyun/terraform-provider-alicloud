# common variables
variable "availability_zone" {
  description = "The available zone to launch ecs instance and other resources."
  default     = ""
}

variable "number_format" {
  description = "The number format used to output."
  default     = "%02d"
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

# Security Group variables
variable "sg_id" {
  description = "The security group id used to launch its rules."
  default     = ""
}

variable "sg_name" {
  description = "The security group name used to launch a new security group when 'sg_id' is not specified."
  default     = "TF_Security_Group"
}

variable "rule_directions" {
  description = "The security group rules direction used to set one or more rules."
  type        = list(string)
  default     = ["ingress"]
}

variable "ip_protocols" {
  description = "The security group rules ip protocol used to set one or more rules."
  type        = list(string)
  default     = []
}

variable "policies" {
  description = "The security group policy used to set one or more rules."
  type        = list(string)
  default     = ["accept"]
}

variable "port_ranges" {
  description = "The security group rules port range used to set one or more rules."
  type        = list(string)
  default     = ["-1/-1"]
}

variable "priorities" {
  description = "The security group rules priority used to set one or more rules."
  type        = list(string)
  default     = [1]
}

variable "cidr_ips" {
  description = "The security group rules cidr_ip used to set one or more rules."
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

# Key pair variables
variable "key_name" {
  description = "The key pair name used to attach one or more instances."
  default     = ""
}

# Disk variables
variable "disk_name" {
  description = "The data disk name used to mark one or more data disks."
  default     = "TF_ECS_Disk"
}

variable "disk_category" {
  description = "The data disk category used to launch one or more data disks."
  default     = "cloud_efficiency"
}

variable "disk_size" {
  description = "The data disk size used to launch one or more data disks."
  default     = "40"
}

variable "disk_tags" {
  description = "Used to mark specified ecs data disks."
  type        = map(string)

  default = {
    created_by   = "Terraform"
    created_from = "module-tf-alicloud-ecs-instance"
  }
}

variable "number_of_disks" {
  description = "The number of launching disks one time."
  default     = 2
}

# Ecs instance variables
variable "image_id" {
  description = "The image id used to launch one or more ecs instances."
  default     = ""
}

variable "instance_type" {
  description = "The instance type used to launch one or more ecs instances."
  default     = ""
}

variable "system_category" {
  description = "The system disk category used to launch one or more ecs instances."
  default     = "cloud_efficiency"
}

variable "system_size" {
  description = "The system disk size used to launch one or more ecs instances."
  default     = "40"
}

variable "instance_name" {
  description = "The instance name used to mark one or more instances."
  default     = "TF-ECS-Instance"
}

variable "host_name" {
  description = "The instance host name used to configure one or more instances.."
  default     = "TF-ECS-Host-Name"
}

variable "password" {
  description = "The password of instance."
  default     = ""
}

variable "internet_charge_type" {
  description = "The internet charge type of instance. Choices are 'PayByTraffic' and 'PayByBandwidth'."
  default     = "PayByTraffic"
}

variable "internet_max_bandwidth_out" {
  description = "The maximum internet out bandwidth of instance.."
  default     = 10
}

variable "instance_charge_type" {
  description = "The charge type of instance. Choices are 'PostPaid' and 'PrePaid'."
  default     = "PostPaid"
}

variable "period" {
  description = "The period of instance when instance charge type is 'PrePaid'."
  default     = 1
}

variable "instance_tags" {
  description = "Used to mark specified ecs instance."
  type        = map(string)

  default = {
    created_by   = "Terraform"
    created_from = "module-tf-alicloud-ecs-instance"
  }
}

variable "number_of_instances" {
  description = "The number of launching instances one time."
  default     = 2
}

