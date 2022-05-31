resource "alicloud_nas_access_group" "foo" {
  access_group_name = "CreateAccessGroup"
  access_group_type = "Vpc"
  description       = "test_AccessG"
  file_system_type  = "extreme"
}
