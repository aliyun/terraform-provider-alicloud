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

# 空桶免费。版本控制 / 生命周期 / 标签全部用 1.284.0 主资源内联块(照抄 r/oss_bucket
# 的 bucket-versioning-lifecycle + bucket-tag-lifecycle 示例)。force_destroy=true 便于探测后清理。
resource "alicloud_oss_bucket" "main" {
  bucket        = var.run_id
  storage_class = "Standard"
  force_destroy = true

  versioning {
    status = "Enabled"
  }

  lifecycle_rule {
    id      = "expire-noncurrent"
    prefix  = "logs/"
    enabled = true

    expiration {
      expired_object_delete_marker = true
    }

    noncurrent_version_expiration {
      days = 90
    }
  }

  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
  }
}
