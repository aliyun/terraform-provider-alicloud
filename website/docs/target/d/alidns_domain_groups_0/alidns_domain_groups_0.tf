data "alicloud_alidns_domain_groups" "example" {
  ids = ["c5ef2bc43064445787adf182af2****"]
}
output "first_domain_group_id" {
  value = "${data.alicloud_alidns_domain_groups.example.groups.0.id}"
}
