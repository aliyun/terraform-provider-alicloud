data "alicloud_bastionhost_host_accounts" "ids" {
  host_id     = "15"
  instance_id = "example_value"
  ids         = ["1", "2"]
}
output "bastionhost_host_account_id_1" {
  value = data.alicloud_bastionhost_host_accounts.ids.accounts.0.id
}

data "alicloud_bastionhost_host_accounts" "nameRegex" {
  host_id     = "15"
  instance_id = "example_value"
  name_regex  = "^my-HostAccount"
}
output "bastionhost_host_account_id_2" {
  value = data.alicloud_bastionhost_host_accounts.nameRegex.accounts.0.id
}

