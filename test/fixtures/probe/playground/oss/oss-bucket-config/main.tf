terraform {
  required_version = ">= 1.5.0"
  required_providers {
    alicloud = { source = "aliyun/alicloud", version = "1.284.0" }
  }
}

variable "run_id" { type = string }

resource "alicloud_oss_bucket" "main" {
  bucket = var.run_id
  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
  }
}
