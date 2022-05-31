resource "alicloud_privatelink_vpc_endpoint_service_resource" "example" {
  resource_id   = "lb-gw8nuym5xxxxx"
  resource_type = "slb"
  service_id    = "epsrv-gw8ii1xxxx"
}

