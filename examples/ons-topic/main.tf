resource "alicloud_ons_instance" "instance" {
  name              = "${var.name}"
  remark            = "terraform-test-instance-remark"
}

resource "alicloud_ons_topic" "default" {
  instance_id		= "${alicloud_ons_instance.instance.id}"
  topic             = "${var.topic}"
  message_type		= "${var.message_type}"
  remark            = "terraform-test-topic-remark"
}
