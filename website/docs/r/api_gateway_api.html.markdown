---
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_api"
sidebar_current: "docs-alicloud-resource-api-gateway-api"
description: |-
  Provides a Alicloud Api Gateway Api Resource.
---

# alicloud_api_gateway_api

Provides an api resource.When you create an API, you must enter the basic information about the API, and define the API request information, the API backend service and response information.

For information about Api Gateway Api and how to use it, see [Create an API](https://www.alibabacloud.com/help/doc-detail/29478.htm)

~> **NOTE:** Terraform will auto build api while it uses `alicloud_api_gateway_api` to build api.

## Example Usage

Basic Usage

```
name = "terraformapi"
group_id = "${alicloud_api_gateway_group.apiGatewayGroup.id}"
description = "description"
auth_type = "APP"
request_config = [
  {
    protocol = "HTTP"
    method = "GET"
    path = "/test/path"
    mode = "MAPPING"
  },
]
service_type = "HTTP"
http_service_config = [
  {
    address = "http://apigateway-backend.alicloudapi.com:8080"
    method = "GET"
    path = "/web/cloudapi"
    timeout = 20
    aone_name = "cloudapi-openapi"
  },
]
request_parameters = [
  {
    name = "testparam"
    type = "STRING"
    required = "OPTIONAL"
    in = "QUERY"
    in_service = "QUERY"
    name_service = "testparams"
  },
]
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the api gateway api. Defaults to null.
* `group_id` - (Required, ForcesNew) The api gateway that the api belongs to. Defaults to null.
* `description` - (Required) The description of the api. Defaults to null.
* `auth_type` - (Required) The authorization Type including APP and ANONYMOUS. Defaults to null.
* `request_config` - (Required, Type: list) Request_config defines how users can send requests to your API.
* `service_type` - (Required) The type of backend service. Type including HTTP,VPC and MOCK. Defaults to null.
* `http_service_config` - (Required, Type: list) http_service_config defines the config when service_type selected 'HTTP'.
* `http_vpc_service_config` - (Required, Type: list) http_service_config defines the config when service_type selected 'HTTP'.
* `mock_service_config` - (Required, Type: list) http_service_config defines the config when service_type selected 'HTTP'.
* `request_parameters` - (Required, Type: list) request_parameters defines .
* `constant_parameters` - (Required, Type: list) http_service_config defines the config when service_type selected 'HTTP'.
* `system_parameters` - (Required, Type: list) http_service_config defines the config when service_type selected 'HTTP'.

### Block request_config

The request_config mapping supports the following:

* `protocol` - (Required) The protocol of api which supports values of 'HTTP','HTTPS' or 'HTTP,HTTPS'
* `method` - (Required) The method of the api, including 'GET','POST','PUT' and etc..
* `path` - (Required) The request path of the api.
* `mode` - (Required) The mode of the parameters between request parameters and service parameters, which support the values of 'MAPPING' and 'PASSTHROUGH'
* `body_format` - (Optional) The body format of the api, which support the values of 'STREAM' and 'FORM'

### Block http_service_config

The http_service_config mapping supports the following:

* `address` - (Required) The address of backend service.
* `path` - (Required) The path of backend service.
* `method` - (Required) The http method of backend service.
* `timeout` - (Optional) Backend service time-out time; unit: millisecond.

### Block http_vpc_service_config

The http_vpc_service_config mapping supports the following:

* `name` - (Required) The name of vpc instance.
* `path` - (Required) The path of backend service.
* `method` - (Required) The http method of backend service.
* `timeout` - (Optional) Backend service time-out time; unit: millisecond.

### Block mock_service_config

The mock_service_config mapping supports the following:

* `result` - (Required) The result of the mock service.

### Block request_parameters

The request_parameters mapping supports the following:

* `name` - (Required) Request's parameter name.
* `type` - (Required) Parameter type which supports values of 'STRING','INT','BOOLEAN','LONG',"FLOAT" and "DOUBLE"
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

* `name` - (Required) System parameter name which supports values including in [system parameter list](https://www.alibabacloud.com/help/doc-detail/43677.html)
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
