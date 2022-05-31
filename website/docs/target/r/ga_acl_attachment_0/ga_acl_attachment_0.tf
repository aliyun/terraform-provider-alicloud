variable "name" {
  default = "tf-testacc-ga"
}

data "alicloud_ga_accelerators" "default" {
  status = "active"
}

data "alicloud_ga_bandwidth_packages" "default" {
  status = "active"
}

resource "alicloud_ga_accelerator" "default" {
  count           = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? 0 : 1
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  count           = length(data.alicloud_ga_bandwidth_packages.default.packages) > 0 ? 0 : 1
  bandwidth       = 20
  type            = "Basic"
  bandwidth_type  = "Basic"
  duration        = 1
  ratio           = 30
  auto_pay        = true
  auto_use_coupon = true
}

locals {
  accelerator_id       = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? data.alicloud_ga_accelerators.default.accelerators.0.id : alicloud_ga_accelerator.default.0.id
  bandwidth_package_id = length(data.alicloud_ga_bandwidth_packages.default.packages) > 0 ? data.alicloud_ga_bandwidth_packages.default.packages.0.id : alicloud_ga_bandwidth_package.default.0.id
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = local.accelerator_id
  bandwidth_package_id = local.bandwidth_package_id
}

resource "alicloud_ga_listener" "default" {
  depends_on     = [alicloud_ga_bandwidth_package_attachment.default]
  accelerator_id = local.accelerator_id
  port_ranges {
    from_port = 60
    to_port   = 70
  }
}

resource "alicloud_ga_acl" "default" {
  acl_name           = var.name
  address_ip_version = "IPv4"
  acl_entries {
    entry             = "192.168.1.0/24"
    entry_description = "tf-test1"
  }
}

resource "alicloud_ga_acl_attachment" "default" {
  acl_id      = alicloud_ga_acl.default.id
  listener_id = alicloud_ga_listener.default.id
  acl_type    = "white"
}
