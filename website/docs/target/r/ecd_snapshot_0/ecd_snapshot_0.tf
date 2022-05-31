variable "name" {
  default = "example_value"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block             = "172.16.0.0/12"
  desktop_access_type    = "Internet"
  office_site_name       = var.name
  enable_internet_access = false
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = "example_value"
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = "example_value"
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}


resource "alicloud_ecd_desktop" "default" {
  office_site_id  = alicloud_ecd_simple_office_site.default.id
  policy_group_id = alicloud_ecd_policy_group.default.id
  bundle_id       = data.alicloud_ecd_bundles.default.bundles.0.id
  desktop_name    = var.name
}

resource "alicloud_ecd_snapshot" "default" {
  description      = var.name
  desktop_id       = alicloud_ecd_desktop.default.id
  snapshot_name    = var.name
  source_disk_type = "SYSTEM"
}
