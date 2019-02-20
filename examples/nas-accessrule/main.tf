resource "alicloud_nas_accessgroup" "main" {
		accessgroup_name = "Test_AccessGroup"
		accessgroup_type = "Classic"
		description = "test_wang"
}
resource "alicloud_nas_accessrule" "main" {
		accessgroup_name = "${alicloud_nas_accessgroup.main.id}"
		sourcecidr_ip = "168.1.1.0/16"
		rwaccess_type = "RDWR"
		useraccess_type = "no_squash"
}