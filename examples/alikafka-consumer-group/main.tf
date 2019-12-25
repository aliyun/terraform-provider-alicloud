resource "alicloud_alikafka_consumer_group" "default" {
  instance_id		= var.instance_id
  consumer_id       = var.consumer_id
}