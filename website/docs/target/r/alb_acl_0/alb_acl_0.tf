resource "alicloud_alb_acl" "example" {
  acl_name = "example_value"
  acl_entries {
    description = "example_value"
    entry       = "10.0.0.0/24"
  }
}

