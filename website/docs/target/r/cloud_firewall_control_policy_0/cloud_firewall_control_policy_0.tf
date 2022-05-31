resource "alicloud_cloud_firewall_control_policy" "example" {
  application_name = "ANY"
  acl_action       = "accept"
  description      = "example"
  destination_type = "net"
  destination      = "100.1.1.0/24"
  direction        = "out"
  proto            = "ANY"
  source           = "1.2.3.0/24"
  source_type      = "net"
}

