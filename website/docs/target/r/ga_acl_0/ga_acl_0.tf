resource "alicloud_ga_acl" "default" {
  acl_name           = "tf-testAccAcl"
  address_ip_version = "IPv4"
  acl_entries {
    entry             = "192.168.1.0/24"
    entry_description = "tf-test1"
  }
}
