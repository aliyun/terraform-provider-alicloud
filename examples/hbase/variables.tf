// export ALICLOUD_REGION="***"
variable "zone_id" {
  default = "cn-shenzhen-b"
}

variable "engine_version" {
  default = "2.0"
}

variable "master_instance_type" {
  default = "hbase.n1.medium"
}

variable "core_instance_type" {
  default = "hbase.n1.large"
}

# 2~20
variable "core_instance_quantity" {
  default = 2
}

variable "core_disk_type" {
  default = "cloud_ssd"
}
// one disk size, unit: GB, default 4 disk per core node; all disk size = coreNodeSize * 4 * core_disk_size(2 * 4 * 100 =800GB)
variable "core_disk_size" {
  default = 100
}
variable "pay_type" {
  default = "Postpaid"
}

// valid when pay_type = "Prepaid"
variable "duration" {
  default = 0
}
// valid when pay_type = "Prepaid"
variable "auto_renew" {
  default = "false"
}

variable "security_ip_list" {
  description = "The security ip list."
  default     = ["127.0.0.1", "127.0.0.2"]
}

variable "is_cold_storage" {
  default     = "false"
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
  default     = "172.16.132.0/24"
}

# VSwitch variables
# VSwitch variables, if vswitch_id is  empty  then the net_type = classic
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
  default     = "172.16.132.0/24"
}

// override  data.alicloud_zones.default.zones[0].id
variable "availability_zone" {
  description = "The available zone to launch ecs instance and other resources."
  default     = ""
}