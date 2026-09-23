output "instance_id" {
  value = alicloud_instance.default.id
}

output "ip_addresses" {
  value = join(",", alicloud_eip_address.default.*.ip_address)
}

