resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = "example_value"
  describe                 = "example_value"
  enable_subscribed        = true
}
resource "alicloud_cms_dynamic_tag_group" "default" {
  contact_group_list = [alicloud_cms_alarm_contact_group.default.id]
  tag_key            = "your_tag_key"
  match_express {
    tag_value                = "your_tag_value"
    tag_value_match_function = "all"
  }
}

