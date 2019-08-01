resource "alicloud_ons_instance" "instance" {
  name              = "${var.name}"
  remark            = "terraform-test-instance-remark"
}

resource "alicloud_ons_group" "default" {
  instance_id		= "${alicloud_ons_instance.instance.id}"
  group_id          = "${var.group_id}"
  remark            = "terraform-test-group-remark"
}
