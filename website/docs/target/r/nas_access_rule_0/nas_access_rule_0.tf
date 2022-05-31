resource "alicloud_nas_access_group" "foo" {
  access_group_name = "tf-NasConfigName"
  access_group_type = "Vpc"
  description       = "tf-testAccNasConfig"
}

resource "alicloud_nas_access_rule" "foo" {
  access_group_name = alicloud_nas_access_group.foo.access_group_name
  source_cidr_ip    = "168.1.1.0/16"
  rw_access_type    = "RDWR"
  user_access_type  = "no_squash"
  priority          = 2
}


