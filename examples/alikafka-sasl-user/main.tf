resource "alicloud_alikafka_sasl_user" "default" {
  instance_id   = "${var.instance_id}"
  username      = "${var.username}"
  password      = "${var.password}"
}