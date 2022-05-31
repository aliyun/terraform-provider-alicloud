resource "alicloud_express_connect_physical_connection" "domestic" {
  access_point_id          = "ap-cn-hangzhou-yh-B"
  line_operator            = "CT"
  peer_location            = "example_value"
  physical_connection_name = "example_value"
  type                     = "VPC"
  description              = "my domestic connection"
  port_type                = "1000Base-LX"
  bandwidth                = 100
}

resource "alicloud_express_connect_physical_connection" "international" {
  access_point_id          = "ap-sg-singpore-A"
  line_operator            = "Other"
  peer_location            = "example_value"
  physical_connection_name = "example_value"
  type                     = "VPC"
  description              = "my domestic connection"
  port_type                = "1000Base-LX"
  bandwidth                = 100
}
