resource "alicloud_bastionhost_host_account" "example" {
  host_account_name = "example_value"
  host_id           = "15"
  instance_id       = "bastionhost-cn-tl32bh0no30"
  protocol_name     = "SSH"
  password          = "YourPassword12345"
}

