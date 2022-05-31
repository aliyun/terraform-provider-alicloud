data "alicloud_alidns_address_pools" "ids" {
  instance_id = "example_value"
  ids         = ["example_value-1", "example_value-2"]
}
output "alidns_address_pool_id_1" {
  value = data.alicloud_alidns_address_pools.ids.pools.0.id
}

data "alicloud_alidns_address_pools" "nameRegex" {
  instance_id = "example_value"
  name_regex  = "^my-AddressPool"
}
output "alidns_address_pool_id_2" {
  value = data.alicloud_alidns_address_pools.nameRegex.pools.0.id
}

