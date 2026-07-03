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

# RAM 与 region 无关(全局服务)。字段用 1.284.0 新名:role_name / assume_role_policy_document /
# policy_name / policy_document(旧 name/document 已废弃)。JSON 均照抄 r/ram_role、r/ram_policy 示例。
resource "alicloud_ram_role" "main" {
  role_name                   = "${var.run_id}-role"
  description                 = "jarvis-probe ram role"
  force                       = true
  assume_role_policy_document = <<EOF
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ecs.aliyuncs.com"
        ]
      }
    }
  ],
  "Version": "1"
}
EOF
  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
  }
}

resource "alicloud_ram_policy" "main" {
  policy_name     = "${var.run_id}-policy"
  description     = "jarvis-probe ram policy"
  force           = true
  policy_document = <<EOF
{
  "Statement": [
    {
      "Action": [
        "oss:ListObjects",
        "oss:GetObject"
      ],
      "Effect": "Allow",
      "Resource": [
        "acs:oss:*:*:mybucket",
        "acs:oss:*:*:mybucket/*"
      ]
    }
  ],
  "Version": "1"
}
EOF
  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
  }
}

# 附加:引用非废弃属性(policy_name / type / role_name)
resource "alicloud_ram_role_policy_attachment" "main" {
  policy_name = alicloud_ram_policy.main.policy_name
  policy_type = alicloud_ram_policy.main.type
  role_name   = alicloud_ram_role.main.role_name
}
