resource "alicloud_api_gateway_group" "apiGatewayGroup"{
  name = "${var.apigateway_group_name}"
  description = "${var.apigateway_group_description}"
}


