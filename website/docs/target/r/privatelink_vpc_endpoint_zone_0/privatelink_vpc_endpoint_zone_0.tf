resource "alicloud_privatelink_vpc_endpoint_zone" "example" {
  endpoint_id = "ep-gw8boxxxxx"
  vswitch_id  = "vsw-rtycxxxxx"
  zone_id     = "eu-central-1a"
}

