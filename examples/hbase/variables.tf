# override  data.alicloud_zones.default.zones[0].id
variable "availability_zone" {
  description = "The available zone to launch ecs instance and other resources."
  default     = ""
}

variable "engine" {
  default = "hbase"
}

variable "engine_version" {
  default = "2.0"
}

variable "master_instance_type" {
  default = "hbase.sn1.large"
}

variable "core_instance_type" {
  default = "hbase.sn1.large"
}

# 2~20
variable "core_instance_quantity" {
  default = 2
}

variable "core_disk_type" {
  default = "cloud_ssd"
}

# node disk size, unit: GB,  disk size per core node; all disk size = coreNodeSize * core_disk_size(2 * 400 =800GB)
variable "core_disk_size" {
  default = 400
}

variable "pay_type" {
  default = "PostPaid"
}

# valid when pay_type = "PrePaid"
variable "duration" {
  default = 1
}
# valid when pay_type = "PrePaid"
variable "auto_renew" {
  default = false
}

# VSwitch variables, if vswitch_id is empty, then the net_type = classic
variable "vswitch_id" {
  description = "The vswitch id used to launch one or more instances."
  default     = ""
}

# 0 mean is_cold_storage = false.
variable "cold_storage_size" {
  default = 0
}

variable "maintain_start_time" {
  default = "02:00Z"
}

variable "maintain_end_time" {
  default = "04:00Z"
}
variable "deletion_protection" {
  default = true
}
