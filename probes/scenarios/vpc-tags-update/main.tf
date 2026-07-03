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

# step1:初始 VPC + 一组标签。step2/ 覆盖层改标签/名称/描述(cidr 不变,保持原地更新而非重建)。
resource "alicloud_vpc" "main" {
  vpc_name    = "${var.run_id}-vpc"
  cidr_block  = "10.0.0.0/16"
  description = "jarvis-probe update scenario v1"
  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
    env        = "dev"
    phase      = "initial"
  }
}
