variable "cluster_name" {
  description = "The name of the cluster that you want to create."
}

variable "cluster_type" {
  description = "The type of the cluster that you want to create. Valid values only: 2: ECS cluster."
}

variable "network_mode" {
  description = "The network type of the cluster that you want to create. Valid values: 1: classic network. 2: VPC."
}

variable "logical_region_id" {
  description = "The ID of the namespace where you want to create the application."
}

variable "vpc_id" {
  description = "The ID of the Virtual Private Cloud (VPC) for the cluster that you want to create. This parameter needs to be specified if the ClusterType is set as VPC."
}