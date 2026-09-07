variable "name" {
  default = "terratest_vpc_flow_log"
}

variable "log_store_name" {
  default = "vpc-flow-log-for-eni"
}

variable "project_name" {
  default = "vpc-flow-log-for-eni"
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/24"
  vpc_name   = var.name
}

resource "alicloud_vpc_flow_log" "default" {
  depends_on     = [alicloud_vpc.default]
  resource_id    = alicloud_vpc.default.id
  resource_type  = "VPC"
  traffic_type   = "All"
  log_store_name = var.log_store_name
  project_name   = var.project_name
  flow_log_name  = var.name
  status         = "Active"
}
