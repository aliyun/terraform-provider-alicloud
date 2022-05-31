resource "alicloud_dms_enterprise_user" "example" {
  uid        = "uid"
  user_name  = "tf-test"
  role_names = ["DBA"]
  mobile     = "1591066xxxx"
}
