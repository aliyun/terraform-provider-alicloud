# Stack Group Configuration
variable "stack_group_name" {
  description = "The name of the ROS stack group"
  type        = string
  default     = "example-multi-region-stack-group"
}

# Deployment Targets
variable "target_regions" {
  description = "List of target region IDs for stack instance deployment (1-20 regions)"
  type        = list(string)
  default     = ["cn-beijing", "cn-shanghai"]
}

variable "target_accounts" {
  description = "List of target Alibaba Cloud account IDs for self-managed permissions (1-50 accounts)"
  type        = list(string)
  default     = ["123456789012****"]
}

variable "rd_folder_ids" {
  description = "List of Resource Directory folder IDs for service-managed permissions (optional)"
  type        = list(string)
  default     = []
}

# Template Parameters
variable "vpc_cidr_block" {
  description = "CIDR block for the VPC to be created in each region"
  type        = string
  default     = "172.16.0.0/12"
}

variable "vswitch_cidr_block" {
  description = "CIDR block for the VSwitch to be created in each region"
  type        = string
  default     = "172.16.0.0/16"
}

# Operation Preferences
variable "max_concurrent_count" {
  description = "Maximum number of concurrent operations across regions and accounts"
  type        = number
  default     = 5
}

variable "failure_tolerance_count" {
  description = "Number of failures tolerated before stopping the operation"
  type        = number
  default     = 0
}

variable "operation_timeout" {
  description = "Timeout in minutes for the stack instances operation (1-1440 minutes)"
  type        = number
  default     = 60
}

# Metadata
variable "environment" {
  description = "Environment name (e.g., dev, staging, production)"
  type        = string
  default     = "development"
}

variable "project_name" {
  description = "Project name for resource tagging"
  type        = string
  default     = "ros-stack-instances-example"
}
