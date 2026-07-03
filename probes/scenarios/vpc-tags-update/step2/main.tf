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

# step2 覆盖层(完整配置,覆盖 step1 的 main.tf):改名 + 改描述 + 改/加标签。
# cidr_block 保持 10.0.0.0/16 → 期望原地更新,不触发 ForceNew 重建。
resource "alicloud_vpc" "main" {
  vpc_name    = "${var.run_id}-vpc-renamed"
  cidr_block  = "10.0.0.0/16"
  description = "jarvis-probe update scenario v2"
  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
    env        = "prod"
    phase      = "updated"
    extra      = "added-in-step2"
  }
}
