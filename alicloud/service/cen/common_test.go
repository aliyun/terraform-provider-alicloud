package cen_test

// EcsInstanceCommonTestCase is duplicated from the root-package test fixture
// (alicloud/service_alicloud_common_test.go). The cen acceptance tests live in
// package cen_test and cannot import a _test.go const from the root package,
// so an identical copy is kept here. Keep this in sync with the root fixture.
const EcsInstanceCommonTestCase = `
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name              = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}
`
