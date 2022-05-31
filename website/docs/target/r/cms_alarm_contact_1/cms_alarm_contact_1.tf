# If you use this template, you can ignore the diff of the alarm contact information by `lifestyle`. We recommend the above usage and activate the link in time.
resource "alicloud_cms_alarm_contact" "example" {
  alarm_contact_name = "zhangsan"
  describe           = "For Test"
  channels_mail      = "terraform.test.com"
  lifecycle {
    ignore_changes = [channels_mail]
  }
}
