resource "alicloud_vpc" "example" {
  name       = "tf-testaccvpc"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "example" {
  availability_zone = "cn-hangzhou-g"
  cidr_block        = "192.168.0.0/16"
  vpc_id            = alicloud_vpc.example.id
}

resource "alicloud_resource_manager_resource_share" "example" {
  resource_share_name = "example_value"
}

resource "alicloud_resource_manager_shared_resource" "example" {
  resource_id       = alicloud_vswitch.example.id
  resource_share_id = alicloud_resource_manager_resource_share.example.resource_share_id
  resource_type     = "VSwitch"
}

