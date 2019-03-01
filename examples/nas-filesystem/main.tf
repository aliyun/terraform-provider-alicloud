resource "alicloud_nas_filesystem" "main" {
		protocol_type = "NFS"
		storage_type = "Performance"
		description = "Create_FileSystem"
}
