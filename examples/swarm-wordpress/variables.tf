# common variables
variable "alicloud_access_key" {
  description = "The Alicloud Access Key ID to launch resources. Support to environment 'ALICLOUD_ACCESS_KEY'."
}

variable "alicloud_secret_key" {
  description = "The Alicloud Access Secret Key to launch resources.  Support to environment 'ALICLOUD_SECRET_KEY'."
}

variable "region" {
  description = "The region to launch resources."
  default     = "cn-hongkong"
}

variable "most_recent" {
  default = true
}

variable "image_owners" {
  default = ""
}

variable "image_name_regex" {
  description = "The name regex of image used to filter image."
  default     = "^centos_6\\w{1,5}[64].*"
}

variable "resource_group_name" {
  description = "A default resource name and it can be used when other resource name is empty."
  default     = "tf-module-concourse"
}

variable "cpu_core_count" {
  default = 1
}

variable "memory_size" {
  default = 2
}

variable "vpc_name" {
  default = "alicloud_vpc"
}

variable "vpc_cidr" {
  default = "10.1.0.0/21"
}

variable "vswitch_cidr" {
  default = "10.1.1.0/24"
}

variable "password" {
  default = "Test12345"
}

variable "cidr_block" {
  default = "172.20.0.0/24"
}

variable "cluster_name" {
  default = "wp-cluster2"
}

variable "node_number" {
  description = "The number of swarm cluster nodes."
  default     = 1
}

variable "disk_category" {
  default = "cloud_efficiency"
}

variable "disk_size" {
  default = 40
}

// Application variables
variable "app_name" {
  description = "The app resource name. Default to variable `resource_group_name`"
  default     = "wordpress"
}

variable "app_version" {
  description = "The app resource version."
  default     = "1.0"
}

variable "latest_image" {
  description = "Whether use the latest image while each update."
  default     = true
}

variable "blue_green" {
  description = "Whether use blue-green release while each update."
  default     = true
}

variable "confirm_blue_green" {
  description = "Confirm a application release which in blue_green."
  default     = true
}

