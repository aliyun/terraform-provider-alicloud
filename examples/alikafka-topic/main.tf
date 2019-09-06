resource "alicloud_alikafka_topic" "default" {
  instance_id       = "${var.instance_id}"
  topic             = "${var.topic}"
  local_topic       = "${var.local_topic}"
  compact_topic     = "${var.compact_topic}"
  partition_num     = "${var.partition_num}"
  remark            = "${var.remark}"
}