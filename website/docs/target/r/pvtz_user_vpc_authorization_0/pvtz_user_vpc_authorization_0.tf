resource "alicloud_pvtz_user_vpc_authorization" "example" {
  authorized_user_id = "example_value"
  auth_channel       = "RESOURCE_DIRECTORY"
  auth_type          = "NORMAL"
}
