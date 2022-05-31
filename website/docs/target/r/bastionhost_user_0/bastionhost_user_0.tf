resource "alicloud_bastionhost_user" "Local" {
  instance_id         = "example_value"
  mobile_country_code = "CN"
  mobile              = "13312345678"
  password            = "YourPassword-123"
  source              = "Local"
  user_name           = "my-local-user"
}

resource "alicloud_bastionhost_user" "Ram" {
  instance_id         = "example_value"
  mobile_country_code = "CN"
  mobile              = "13312345678"
  password            = "YourPassword-123"
  source              = "Ram"
  source_user_id      = "1234567890"
  user_name           = "my-ram-user"
}
