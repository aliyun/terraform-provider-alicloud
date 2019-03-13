resource "alicloud_nas_access_group" "main" {
		name = "tf-testAccNasConfigName"
		type = "Classic"
		description = "tf-testAccNasConfig"
}
resource "alicloud_nas_access_rule" "main" {
		access_group_name = "${alicloud_nas_access_group.main.id}"
		source_cidr_ip = "168.1.1.0/16"
		rw_access_type = "RDWR"
		user_access_type = "no_squash"
		priority = 2
}
