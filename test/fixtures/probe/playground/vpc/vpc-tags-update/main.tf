terraform {
  required_version = ">= 1.5.0"
  required_providers {
    alicloud = { source = "aliyun/alicloud", version = "1.284.0" }
  }
}

variable "run_id" { type = string }

resource "alicloud_vpc" "main" {
  vpc_name   = "probe-${var.run_id}"
  cidr_block = "172.16.0.0/16"
  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
  }
}
