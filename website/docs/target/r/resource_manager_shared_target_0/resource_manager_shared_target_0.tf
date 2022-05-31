data "alicloud_resource_manager_accounts" "example" {}

resource "alicloud_resource_manager_resource_share" "example" {
  resource_share_name = "example_value"
}

resource "alicloud_resource_manager_shared_target" "example" {
  resource_share_id = alicloud_resource_manager_resource_share.example.resource_share_id
  target_id         = data.alicloud_resource_manager_accounts.example.ids.0
}

