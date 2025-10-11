variable "name" {
  default = "xdang_terraform2"
}

variable "creation" {
  default = ""
}

variable "db_version" {
  default = "8.0"
}

variable "db_type" {
  default = "MySQL"
}

variable "storage" {
  default = "10"
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

variable "db_node_num" {
  description = "Number of the PolarDB cluster nodes, default is 2(Each cluster must contain at least a primary node and a read-only node). Add/remove nodes by modifying this parameter, valid values: [2~16]."
  type        = number
  default     = 2
}

variable "db_cluster_nodes_configs" {
  description = "The advanced configuration for all nodes in the cluster except for the RW node, including db_node_class, hot_replica_mode, and imci_switch properties."
  type        = map(object({
    db_node_class           = string
    db_node_role            = optional(string,null)
    hot_replica_mode        = optional(string,null)
    imci_switch             = optional(string,null)
  }))
  default     = {}
}

