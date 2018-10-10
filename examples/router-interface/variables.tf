# Initiating side variables
variable "region" {
  description = "The region to launch resources."
  default     = "cn-beijing"
}

variable "opposite_region" {
  description = "The opposite region to launch resources."
  default     = "cn-hangzhou"
}

variable "init_vpc_id" {
  description = "The vpc id used to launch several vswitches."
  default     = ""
}

variable "accept_vpc_id" {
  description = "The vpc id used to launch several vswitches."
  default     = ""
}

variable "vpc_name" {
  description = "The vpc name used to launch a new vpc when 'vpc_id' is not specified."
  default     = "example-router-interface"
}

variable "vpc_description" {
  description = "The vpc description used to launch a new vpc when 'vpc_id' is not specified."
  default     = "A new VPC created by Terrafrom example router-interface"
}

variable "init_vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "10.0.0.0/8"
}

variable "accept_vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "172.16.0.0/12"
}

variable "interface_spec" {
  description = "The router interface specification."
  default     = "Large.2"
}
