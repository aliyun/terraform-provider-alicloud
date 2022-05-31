data "alicloud_alidns_instances" "example" {
  ids = ["dns-cn-oew1npk****"]
}
output "first_instance_id" {
  value = "${data.alicloud_alidns_instances.example.instances.0.id}"
}
