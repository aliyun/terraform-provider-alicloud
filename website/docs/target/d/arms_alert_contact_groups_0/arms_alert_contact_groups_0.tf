
data "alicloud_arms_alert_contact_groups" "nameRegex" {
  name_regex = "^my-AlertContactGroup"
}
output "arms_alert_contact_group_id" {
  value = data.alicloud_arms_alert_contact_groups.nameRegex.groups.0.id
}

