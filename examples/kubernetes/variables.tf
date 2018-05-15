# common variables
variable "alicloud_access_key" {
  description = "The Alicloud Access Key ID to launch resources."
}
variable "alicloud_secret_key" {
  description = "The Alicloud Access Secret Key to launch resources."
}
variable "region" {
  description = "The region to launch resources."
  default = "cn-hangzhou"
}
variable "availability_zone" {
  description = "The available zone to launch ecs instance and other resources."
  default = ""
}
variable "number_format" {
  description = "The number format used to output."
  default = "%02d"
}
variable "example_name" {
  default = "tf-example-kubernetes"
}

# Instance typs variables
variable "cpu_core_count" {
  description = "CPU core count is used to fetch instance types."
  default = 1
}

variable "memory_size" {
  description = "Memory size used to fetch instance types."
  default = 2
}

# VPC variables
variable "vpc_name" {
  description = "The vpc name used to create a new vpc when 'vpc_id' is not specified. Default to variable `example_name`"
  default = ""
}
variable "vpc_id" {
  description = "A existing vpc id used to create several vswitches and other resources."
  default = ""
}
variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default = "10.1.0.0/21"
}

# VSwitch variables
variable "vswitch_name_prefix" {
  description = "The vswitch name prefix used to create several new vswitches. Default to variable `example_name`"
  default = ""
}

variable "vswitch_ids" {
  description = "List of existing vswitch id."
  type = "list"
  default = []
}

variable "vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'vswitch_ids' is not specified."
  type = "list"
  default = ["10.1.2.0/24"]
}

variable "new_nat_gateway" {
  description = "Whether to create a new nat gateway. In this template, a new nat gateway will create a nat gateway, eip and server snat entries."
  default = "true"
}

# Cluster nodes variables

variable "master_instance_type" {
  description = "The ecs instance type used to launch master nodes. Default from instance typs datasource."
  default = ""
}

variable "worker_instance_type" {
  description = "The ecs instance type used to launch worker nodes. Default from instance typs datasource."
  default = ""
}

variable "master_disk_category" {
  description = "The system disk category used to launch one or more master nodes."
  default = "cloud_efficiency"
}

variable "worker_disk_category" {
  description = "The system disk category used to launch one or more worker nodes."
  default = "cloud_efficiency"
}

variable "master_disk_size" {
  description = "The system disk size used to launch one or more master nodes."
  default = "40"
}

variable "worker_disk_size" {
  description = "The system disk size used to launch one or more worker nodes."
  default = "40"
}

variable "ecs_password" {
  description = "The password of instance."
  default = "Abc12345"
}

variable "k8s_number" {
  description = "The number of kubernetes cluster."
  default = 1
}

variable "k8s_worker_number" {
  description = "The number of worker nodes in each kubernetes cluster."
  default = 3
}

variable "k8s_name_prefix" {
  description = "The name prefix used to create several kubernetes clusters. Default to variable `example_name`"
  default = ""
}

variable "k8s_pod_cidr" {
  description = "The kubernetes pod cidr block. It cannot be equals to vpc's or vswitch's and cannot be in them."
  default = "172.20.0.0/16"
}

variable "k8s_service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  default = "172.21.0.0/20"
}