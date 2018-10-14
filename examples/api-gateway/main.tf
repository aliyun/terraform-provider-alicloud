resource "alicloud_api_gateway_group" "apiGatewayGroup" {
  name        = "${var.apigateway_group_name}"
  description = "${var.apigateway_group_description}"
}

data "alicloud_api_gateway_groups" "data_apigatway_groups" {
  name_regex  = "${alicloud_api_gateway_group.apiGatewayGroup.name}"
  output_file = "output_ApiGatawayGroups"
}

output "first_group_id" {
  value = "${data.alicloud_api_gateway_groups.data_apigatway_groups.groups.0.id}"
}
