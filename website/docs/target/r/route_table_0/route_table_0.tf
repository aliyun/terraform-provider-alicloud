resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "vpc-example-name"
}

resource "alicloud_route_table" "foo" {
  vpc_id           = alicloud_vpc.foo.id
  route_table_name = "route-table-example-name"
  description      = "route-table-example-description"
}
