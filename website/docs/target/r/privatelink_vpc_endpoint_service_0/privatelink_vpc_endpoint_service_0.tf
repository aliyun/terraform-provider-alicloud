resource "alicloud_privatelink_vpc_endpoint_service" "example" {
  service_description    = "tftest"
  connect_bandwidth      = 103
  auto_accept_connection = false
}

