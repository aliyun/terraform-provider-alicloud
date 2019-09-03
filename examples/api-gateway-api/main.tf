resource "alicloud_api_gateway_group" "apiGatewayGroup" {
  name        = "${var.apigateway_group_name}"
  description = "${var.apigateway_group_description}"
}

resource "alicloud_api_gateway_api" "apiGatewayApi" {
  name        = "terraformapi"
  group_id    = "${alicloud_api_gateway_group.apiGatewayGroup.id}"
  description = "description"
  auth_type   = "APP"

  request_config {
    protocol = "HTTP"
    method   = "GET"
    path     = "/test/path1"
    mode     = "MAPPING"
  }

  service_type = "HTTP"

  http_service_config {
    address   = "http://apigateway-backend.alicloudapi.com:8080"
    method    = "GET"
    path      = "/web/cloudapi"
    timeout   = 12
    aone_name = "cloudapi-openapi"
  }

  request_parameters {
    name         = "aa"
    type         = "STRING"
    required     = "REQUIRED"
    in           = "QUERY"
    in_service   = "QUERY"
    name_service = "testparams"
  }

  stage_names = [
    "RELEASE",
    "PRE",
    "TEST",
  ]
}

resource "alicloud_api_gateway_app" "apiGatewayApp" {
  name        = "${var.apigateway_app_name_test}"
  description = "${var.apigateway_app_description_test}"
}

resource "alicloud_api_gateway_app_attachment" "foo" {
  api_id     = "${alicloud_api_gateway_api.apiGatewayApi.api_id}"
  group_id   = "${alicloud_api_gateway_group.apiGatewayGroup.id}"
  stage_name = "PRE"
  app_id     = "${alicloud_api_gateway_app.apiGatewayApp.id}"
}
