output "first_group_id" {
  value = "${data.alicloud_api_gateway_groups.data_apigatway_groups.groups.0.id}"
}

output "first_api_id" {
  value = "${data.alicloud_api_gateway_apis.data_apigatway_apis.apis.0.id}"
}
