variable "name" {
  default = "example_name"
}

resource "alicloud_dfs_access_group" "default" {
  network_type      = "VPC"
  access_group_name = var.name
  description       = var.name
}

resource "alicloud_dfs_access_rule" "default" {
  network_segment = "192.0.2.0/24"
  access_group_id = alicloud_dfs_access_group.default.id
  description     = var.name
  rw_access_type  = "RDWR"
  priority        = "10"
}

