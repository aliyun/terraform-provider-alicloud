resource "alicloud_msc_sub_subscription" "example" {
  item_name      = "Notifications of Product Expiration"
  sms_status     = "1"
  email_status   = "1"
  pmsg_status    = "1"
  tts_status     = "1"
  webhook_status = "0"
}
