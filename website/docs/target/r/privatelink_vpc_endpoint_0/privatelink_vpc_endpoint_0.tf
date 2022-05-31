resource "alicloud_privatelink_vpc_endpoint" "example" {
  service_id        = "YourServiceId"
  security_group_id = ["sg-ercx1234"]
  vpc_id            = "YourVpcId"
}
