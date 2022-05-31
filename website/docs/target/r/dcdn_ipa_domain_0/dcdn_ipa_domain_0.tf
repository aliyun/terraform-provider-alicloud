data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
resource "alicloud_dcdn_ipa_domain" "example" {
  domain_name       = "example.com"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  sources {
    content  = "1.1.1.1"
    port     = 80
    priority = "20"
    type     = "ipaddr"
    weight   = 10
  }
  scope  = "overseas"
  status = "online"
}
