resource "alicloud_vpc" "default" {
  vpc_name   = "tf-testacc-vpcname"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_cms_monitor_group" "default" {
  monitor_group_name = "tf-testaccmonitorgroup"
}

resource "alicloud_cms_monitor_group_instances" "example" {
  group_id = alicloud_cms_monitor_group.default.id
  instances {
    instance_id   = alicloud_vpc.default.id
    instance_name = "tf-testacc-vpcname"
    region_id     = "cn-hangzhou"
    category      = "vpc"
  }
}
