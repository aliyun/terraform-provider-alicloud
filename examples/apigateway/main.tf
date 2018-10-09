data "alicloud_apigateway_groups" "data_apigatway"{
  output_file = "outgroups"
}

output "first_group_id" {
  value = "${data.alicloud_apigateway_groups.data_apigatway.groups.0.id}"
}

resource "alicloud_apigateway_api" "api" {
  name = "${var.apigateway_group_name}"
  description = "${var.apigateway_group_description}"
  group_id = "${data.alicloud_apigateway_groups.data_apigatway.groups.0.id}"
  auth_type = "APP"
  request_config = "{\"RequestProtocol\":\"HTTP\",\"RequestHttpMethod\":\"GET\",\"RequestPath\":\"/test\",\"BodyFormat\":\"\",\"PostBodyDescription\":\"\",\"RequestMode\":\"MAPPING\"}"
  service_config = "{\"ServiceProtocol\":\"HTTP\",\"ServiceHttpMethod\":\"GET\",\"ServiceAddress\":\"http://apigateway-backend.alicloudapi.com:8080\",\"ServiceTimeout\":\"10000\",\"ServicePath\":\"/web/cloudapi\",\"Mock\":\"FALSE\",\"MockResult\":\"\",\"ServiceVpcEnable\":\"FALSE\",\"VpcConfig\":{},\"FunctionComputeConfig\":{},\"ContentTypeCatagory\":\"DEFAULT\",\"ContentTypeValue\":\"application/x-www-form-urlencoded; charset=UTF-8\",\"AoneAppName\":\"cloudapi-openapi\"}"
  request_parameters = "[{\"ParameterLocation\":{\"name\":\"Head\",\"orderNumber\":2},\"ParameterType\":\"String\",\"Required\":\"OPTIONAL\",\"isHide\":false,\"ApiParameterName\":\"requestHead\",\"Location\":\"Head\"}]"
  service_parameters = "[{\"ServiceParameterName\":\"requestHead\",\"Location\":\"Query\",\"Type\":\"String\",\"ParameterCatalog\":\"REQUEST\"}]"
  service_parameters_map = "[{\"ServiceParameterName\":\"requestHead\",\"RequestParameterName\":\"requestHead\"}]"
}



