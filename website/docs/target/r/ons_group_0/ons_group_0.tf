variable "name" {
  default = "onsInstanceName"
}

variable "group_name" {
  default = "GID-onsGroupDatasourceName"
}

resource "alicloud_ons_instance" "default" {
  name   = var.name
  remark = "default_ons_instance_remark"
}

resource "alicloud_ons_group" "default" {
  group_name  = var.group_name
  instance_id = alicloud_ons_instance.default.id
  remark      = "dafault_ons_group_remark"
}
