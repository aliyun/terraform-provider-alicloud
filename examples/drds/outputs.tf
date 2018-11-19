output "instance" {
  value = "${alicloud_drds_instance.instance.id}"
}

output "description" {
  value = "${alicloud_drds_instance.instance.description}"
}

output "type" {
  value = "${alicloud_drds_instance.instance.type}"
}

output "zone_id" {
  value = "${alicloud_drds_instance.instance.zone_id}"
}

output "specification" {
  value = "${alicloud_drds_instance.instance.specification}"
}

output "pay_type" {
  value = "${alicloud_drds_instance.instance.pay_type}"
}

output "instance_series" {
  value = "${alicloud_drds_instance.instance.instance_series}"
}
