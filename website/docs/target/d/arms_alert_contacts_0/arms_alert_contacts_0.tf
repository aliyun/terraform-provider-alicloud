data "alicloud_arms_alert_contacts" "ids" {}
output "arms_alert_contact_id_1" {
  value = data.alicloud_arms_alert_contacts.ids.contacts.0.id
}

data "alicloud_arms_alert_contacts" "nameRegex" {
  name_regex = "^my-AlertContact"
}
output "arms_alert_contact_id_2" {
  value = data.alicloud_arms_alert_contacts.nameRegex.contacts.0.id
}

