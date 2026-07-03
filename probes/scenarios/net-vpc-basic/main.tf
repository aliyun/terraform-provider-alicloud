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

# 可用区:照抄 d/zones + r/vswitch 官方示例,让 provider 选一个能建 VSwitch 的 AZ
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "main" {
  vpc_name   = "${var.run_id}-vpc"
  cidr_block = "172.16.0.0/16"
  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
  }
}

resource "alicloud_vswitch" "main" {
  vswitch_name = "${var.run_id}-vsw"
  vpc_id       = alicloud_vpc.main.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
  }
}

resource "alicloud_security_group" "main" {
  security_group_name = "${var.run_id}-sg"
  vpc_id              = alicloud_vpc.main.id
  tags = {
    managed_by = "jarvis-probe"
    run_id     = var.run_id
  }
}

# VPC 安全组:nic_type 只能是 intranet(照抄 r/security_group_rule 示例)
resource "alicloud_security_group_rule" "allow_ssh" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alicloud_security_group.main.id
  cidr_ip           = "0.0.0.0/0"
}
