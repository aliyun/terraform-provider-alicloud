resource "alicloud_nas_filesystem" "main" {
		protocol_type = "NFS"
		storage_type = "Performance"
		description = "Create_FileSystem"
}

resource "alicloud_nas_accessgroup" "main" {
		accessgroup_name = "Test_AccessGroup"
		accessgroup_type = "Classic"
		description = "test_wang"
}
resource "alicloud_nas_mounttarget" "main" {
		filesystem_id = "${alicloud_nas_filesystem.main.id}"
		accessgroup_name = "${alicloud_nas_accessgroup.main.id}"
		networktype = "Classic"
}