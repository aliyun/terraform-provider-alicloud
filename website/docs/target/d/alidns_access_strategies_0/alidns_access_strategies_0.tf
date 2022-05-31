data "alicloud_alidns_access_strategies" "ids" {
  instance_id   = "example_value"
  strategy_mode = "example_value"
  ids           = ["example_value-1", "example_value-2"]
  name_regex    = "the_resource_name"
}
output "alidns_access_strategy_id_1" {
  value = data.alicloud_alidns_access_strategies.ids.strategies.0.id
}

