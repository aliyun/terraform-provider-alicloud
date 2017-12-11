# common variables
variable "alicloud_access_key" {
  description = "The Alicloud Access Key ID to launch resources."
  default = ""
}
variable "alicloud_secret_key" {
  description = "The Alicloud Access Secret Key to launch resources."
  default = ""
}
variable "region" {
  description = "The region to launch resources."
  default = "cn-hangzhou"
}
variable "availability_zones" {
  description = "List available zones to launch several VSwitches."
  type = "list"
  default = [""]

}
variable "number_format" {
  description = "The number format used to output."
  default = "%02d"
}

# VPC variables
variable "vpc_id" {
  description = "The vpc id used to launch several vswitches."
  default = ""
}
variable "vpc_name" {
  description = "The vpc name used to launch a new vpc when 'vpc_id' is not specified."
  default = "TF-VPC"
}
variable "vpc_description" {
  description = "The vpc description used to launch a new vpc when 'vpc_id' is not specified."
  default = "A new VPC created by Terrafrom module tf-alicloud-vpc-cluster"
}
variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default = "172.16.0.0/12"
}

# VSwitch variables
variable "vswitch_cidrs" {
  description = "List of cidr blocks used to launch several new vswitches."
  type = "list"
  default = ["172.16.1.0/24"]
}
variable "vswitch_name" {
  description = "The vswitch name prefix used to launch several new vswitch."
  default = "TF_VSwitch"
}

variable "vswitch_description" {
  description = "The vswitch description used to launch several new vswitch."
  default = "New VSwitch created by Terrafrom module tf-alicloud-vpc-cluster."
}

// According to the vswitch cidr blocks to launch several vswitches
variable "route_table_id" {
  description = "The route table ID of virtual router in the specified VPC."
  default = ""
}
variable "destination_cidrs" {
  description = "List of destination CIDR block of virtual router in the specified VPC."
  type = "list"
  default = []
}
variable "nexthop_ids" {
  description = "List of next hop instance IDs of virtual router in the specified VPC."
  type = "list"
  default = []
}

// Interface parameters
variable "opposite_region" {
  description = "The opposite region to launch resources."
  default = "cn-beijing"
}
variable "interface_role" {
  description = "The router interface role. Choices are 'InitiatingSide' and 'AcceptingSide'."
  default = "InitiatingSide"
}
variable "interface_spec" {
  description = "The router interface specification."
  default = "Small.2"
}