resource "alicloud_cms_monitor_group" "example" {
  monitor_group_name = "tf-testaccmonitorgroup"
}

resource "alicloud_cms_monitor_group" "default2" {
  contact_groups      = ["your_contact_groups"]
  resource_group_id   = "your_resource_group_id"
  resource_group_name = "resource_group_name"
  tags = {
    Created = "TF"
    For     = "Acceptance-test"
  }
}
