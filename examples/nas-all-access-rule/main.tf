resource "alicloud_nas_access_group" "main" {
  name        = "tf-testAccNasConfigName"
  type        = "Vpc"
  description = "tf-testAccNasConfig"
  file_system_type = "extreme"
}

resource "alicloud_nas_all_access_rule" "main" {
  access_group_name = alicloud_nas_access_group.main.name
  file_system_type  = alicloud_nas_access_group.main.file_system_type
  source_cidr_ip    = "168.1.1.0/16"
  rw_access_type    = "RDWR"
  user_access_type  = "no_squash"
  priority          = 3
}

