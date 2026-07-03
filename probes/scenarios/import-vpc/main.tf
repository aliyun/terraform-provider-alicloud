terraform {
  required_version = ">= 1.5.0"
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "1.284.0"
    }
  }
}

variable "run_id" {
  type = string
}

resource "alicloud_vpc" "main" {
  vpc_name   = "${var.run_id}-vpc"
  cidr_block = "192.168.0.0/16"
  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
  }
}

# import 往返用:runner 从此 output 取真实 id 做 terraform import
output "vpc_id" {
  value = alicloud_vpc.main.id
}
