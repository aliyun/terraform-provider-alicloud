resource "alicloud_cloud_firewall_control_policy" "example1" {
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

resource "alicloud_cloud_firewall_control_policy_order" "example2" {
  acl_uuid  = alicloud_cloud_firewall_control_policy.example1.acl_uuid
  direction = alicloud_cloud_firewall_control_policy.example1.direction
  order     = 1
}

