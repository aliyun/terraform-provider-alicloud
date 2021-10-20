---
subcategory: "API Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_api"
sidebar_current: "docs-alicloud-resource-api-gateway-api"
description: |-
  Provides a Alicloud Api Gateway Api Resource.
---

# alicloud_api_gateway_api

Provides an api resource.When you create an API, you must enter the basic information about the API, and define the API request information, the API backend service and response information.

For information about Api Gateway Api and how to use it, see [Create an API](https://www.alibabacloud.com/help/doc-detail/29478.htm)

-> **NOTE:** Terraform will auto build api while it uses `alicloud_api_gateway_api` to build api.

## Example Usage

Basic Usage

```
resource "alicloud_api_gateway_group" "apiGroup" {
  name        = "ApiGatewayGroup"
  description = "description of the api group"
}

resource "alicloud_api_gateway_api" "apiGatewayApi" {
  name        = alicloud_api_gateway_group.apiGroup.name
  group_id    = alicloud_api_gateway_group.apiGroup.id
  description = "your description"
  auth_type   = "APP"
  force_nonce_check = false

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
    name         = "aaa"
    type         = "STRING"
    required     = "OPTIONAL"
    in           = "QUERY"
    in_service   = "QUERY"
    name_service = "testparams"
  }

  stage_names = [
    "RELEASE",
    "TEST",
  ]
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the api gateway api. Defaults to null.
* `group_id` - (Required, ForcesNew) The api gateway that the api belongs to. Defaults to null.
* `description` - (Required) The description of the api. Defaults to null.
* `auth_type` - (Required) The authorization Type including APP and ANONYMOUS. Defaults to null.
* `request_config` - (Required, Type: list) Request_config defines how users can send requests to your API.
* `service_type` - (Required) The type of backend service. Type including HTTP,VPC and MOCK. Defaults to null.
* `http_service_config` - (Optional, Type: list) http_service_config defines the config when service_type selected 'HTTP'.
* `http_vpc_service_config` - (Optional, Type: list) http_vpc_service_config defines the config when service_type selected 'HTTP-VPC'.
* `fc_service_config` - (Optional, Type: list) fc_service_config defines the config when service_type selected 'FunctionCompute'.
* `mock_service_config` - (Optional, Type: list) http_service_config defines the config when service_type selected 'MOCK'.
* `request_parameters` - (Required, Type: list) request_parameters defines the request parameters of the api.
* `constant_parameters` - (Required, Type: list) constant_parameters defines the constant parameters of the api.
* `system_parameters` - (Required, Type: list) system_parameters defines the system parameters of the api.
* `stage_names` - (Optional, Type: list) Stages that the api need to be deployed. Valid value: `RELEASE`,`PRE`,`TEST`.
* `force_nonce_check` - (Optional, Type: bool, Available in v1.140+) Whether to prevent API replay attack. Default value: `false`.

### Block request_config

The request_config mapping supports the following:

* `protocol` - (Required) The protocol of api which supports values of 'HTTP','HTTPS' or 'HTTP,HTTPS'.
* `method` - (Required) The method of the api, including 'GET','POST','PUT' etc.
* `path` - (Required) The request path of the api.
* `mode` - (Required) The mode of the parameters between request parameters and service parameters, which support the values of 'MAPPING' and 'PASSTHROUGH'.
* `body_format` - (Optional) The body format of the api, which support the values of 'STREAM' and 'FORM'.

### Block http_service_config

The http_service_config mapping supports the following:

* `address` - (Required) The address of backend service.
* `path` - (Required) The path of backend service.
* `method` - (Required) The http method of backend service.
* `timeout` - (Required) Backend service time-out time; unit: millisecond.

### Block http_vpc_service_config

The http_vpc_service_config mapping supports the following:

* `name` - (Required) The name of vpc instance.
* `path` - (Required) The path of backend service.
* `method` - (Required) The http method of backend service.
* `timeout` - (Required) Backend service time-out time. Unit: millisecond.

### Block fc_vpc_service_config

The fc_service_config mapping supports the following:

* `region` - (Required) The region that the function compute service belongs to.
* `function_name` - (Required) The function name of function compute service.
* `service_name` - (Required) The service name of function compute service.
* `arn_role` - (Optional) RAM role arn attached to the Function Compute service. This governs both who / what can invoke your Function, as well as what resources our Function has access to. See [User Permissions](https://www.alibabacloud.com/help/doc-detail/52885.htm) for more details.
* `timeout` - (Required) Backend service time-out time; unit: millisecond.

### Block mock_service_config

The mock_service_config mapping supports the following:

* `result` - (Required) The result of the mock service.

### Block request_parameters

The request_parameters mapping supports the following:

* `name` - (Required) Request's parameter name.
* `type` - (Required) Parameter type which supports values of 'STRING','INT','BOOLEAN','LONG',"FLOAT" and "DOUBLE".
* `required` - (Required) Parameter required or not; values: REQUIRED and OPTIONAL.
* `in` - (Required) Request's parameter location; values: BODY, HEAD, QUERY, and PATH.
* `in_service` - (Required) Backend service's parameter location; values: BODY, HEAD, QUERY, and PATH.
* `name_service` - (Required) Backend service's parameter name.
* `description` - (Optional) The description of parameter.
* `default_value` - (Optional) The default value of the parameter.

### Block constant_parameters

The constant_parameters mapping supports the following:

* `name` - (Required) Constant parameter name.
* `in` - (Required) Constant parameter location; values: 'HEAD' and 'QUERY'.
* `value` - (Required) Constant parameter value.
* `description` - (Optional) The description of Constant parameter.

### Block system_parameters

The system_parameters mapping supports the following:

* `name` - (Required) System parameter name which supports values including in [system parameter list](https://www.alibabacloud.com/help/doc-detail/43677.html).
* `in` - (Required) System parameter location; values: 'HEAD' and 'QUERY'.
* `name_service` - (Required) Backend service's parameter name.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the api resource of api gateway.
* `api_id` - The ID of the api of api gateway.

## Import

Api gateway api can be imported using the id.Format to `<API Group Id>:<API Id>` e.g.

```
$ terraform import alicloud_api_gateway_api.example "ab2351f2ce904edaa8d92a0510832b91:e4f728fca5a94148b023b99a3e5d0b62"
```
