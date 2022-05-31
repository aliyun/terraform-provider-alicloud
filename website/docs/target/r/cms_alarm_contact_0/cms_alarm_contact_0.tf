# If you use this template, you need to activate the link before you can return to the alarm contact information, otherwise diff will appear in terraform. So please confirm the activation link as soon as possible.
resource "alicloud_cms_alarm_contact" "example" {
  alarm_contact_name = "zhangsan"
  describe           = "For Test"
  channels_mail      = "terraform.test.com"
}
