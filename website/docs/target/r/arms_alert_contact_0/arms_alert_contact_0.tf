resource "alicloud_arms_alert_contact" "example" {
  alert_contact_name     = "example_value"
  ding_robot_webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=91f2f6****"
  email                  = "someone@example.com"
  phone_num              = "1381111****"
}

