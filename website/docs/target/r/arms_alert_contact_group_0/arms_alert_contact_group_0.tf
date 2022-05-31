resource "alicloud_arms_alert_contact" "example" {
  alert_contact_name     = "example_value"
  ding_robot_webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=91f2f6****"
  email                  = "someone@example.com"
  phone_num              = "1381111****"
}
resource "alicloud_arms_alert_contact_group" "example" {
  alert_contact_group_name = "example_value"
  contact_ids              = [alicloud_arms_alert_contact.example.id]
}

