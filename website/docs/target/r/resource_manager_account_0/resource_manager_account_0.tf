# Add a Resource Manager Account.
resource "alicloud_resource_manager_folder" "f1" {
  folder_name = "test1"
}

resource "alicloud_resource_manager_account" "example" {
  display_name = "RDAccount"
  folder_id    = alicloud_resource_manager_folder.f1.id
}
