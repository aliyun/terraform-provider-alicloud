variable "instance_type" {
  default = "cassandra.c.large"
}

# 2~20
variable "node_count" {
  default = 2
}

variable "disk_type" {
  default = "cloud_ssd"
}

# node disk size, unit: GB,  disk size per core node; all disk size = disk_size * node_count(2 * 160 =320GB)
variable "disk_size" {
  default = 160
}

variable "pay_type" {
  default = "PayAsYouGo"
}

# valid when pay_type = "PayAsYouGo"
variable "auto_renew_period" {
  default = 1
}
# valid when pay_type = "PayAsYouGo"
variable "auto_renew" {
  default = false
}

variable "major_version" {
  default = "3.11"
}

variable "dc_name_1" {
  default = "dc-1"
}

variable "dc_name_2" {
  default = "dc-2"
}

variable "zone_id_1" {
  description = "The zone-2 id used to launch one or more clusters."
  default     = ""
}

variable "zone_id_2" {
  description = "The zone-2 id used to launch one or more clusters."
  default     = ""
}

variable "vswitch_id_1" {
  description = "The vswitch id used to launch one or more clusters."
  default     = ""
}

variable "vswitch_id_2" {
  description = "The vswitch id used to launch one or more clusters."
  default     = ""
}

variable "maintain_start_time" {
  default = "08:00Z"
}

variable "maintain_end_time" {
  default = "10:00Z"
}

variable "password" {
  default = "Admin123"
}

variable "enable_public" {
  default =  false
}

variable "ip_white" {
  default =  "127.0.0.1"
}
