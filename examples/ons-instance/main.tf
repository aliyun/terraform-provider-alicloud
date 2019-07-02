resource "alicloud_ons_instance" "instance" {
  name              = "${var.name}"
  remark            = "terraform-test-remark"
}
