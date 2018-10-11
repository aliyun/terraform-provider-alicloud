provider "alicloud" {
  region = "cn-hangzhou"
}
 resource "alicloud_drds_instance" "instance" {
  provider = "alicloud"
  description = "${var.description}"
  type = "${var.type}"
  zone_id = "${var.zone_id}"
  specification = "${var.specification}"
  pay_type = "${var.pay_type}"
  vswitch_id = "${var.vswitch_id}"
  instance_series = "${var.instance_series}"
}