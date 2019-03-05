resource "alicloud_nas_file_system" "main" {
		protocol_type = "NFS"
		storage_type = "Performance"
		description = "Create_FileSystem"
}
