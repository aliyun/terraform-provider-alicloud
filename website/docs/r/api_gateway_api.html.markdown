---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_api"
sidebar_current: "docs-alicloud-resource-api-gateway-api"
description: |-
  Provides a Alicloud Api Gateway Api Resource.
---

# alicloud_api_gateway_api

Provides an api resource.When you create an API, you must enter the basic information about the API, and define the API request information, the API backend service and response information.

For information about Api Gateway Api and how to use it, see [Create an API](https://www.alibabacloud.com/help/en/api-gateway/latest/api-cloudapi-2016-07-14-createapi)

-> **NOTE:** Available since v1.22.0.

-> **NOTE:** Terraform will auto build api while it uses `alicloud_api_gateway_api` to build api.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_api&exampleId=efb02e53-1ba1-7435-c04d-29fbb459a7e2e612d86a&activeTab=example&spm=docs.r.api_gateway_api.0.efb02e531b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_api_gateway_group" "example" {
  name        = "tf-example"
  description = "tf-example"
  base_path   = "/"
}

resource "alicloud_api_gateway_api" "example" {
  group_id          = alicloud_api_gateway_group.example.id
  name              = "tf-example"
  description       = "tf-example"
  auth_type         = "APP"
  force_nonce_check = false

  request_config {
    protocol = "HTTP"
    method   = "GET"
    path     = "/example/path"
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
    name         = "example"
    type         = "STRING"
    required     = "OPTIONAL"
    in           = "QUERY"
    in_service   = "QUERY"
    name_service = "exampleservice"
  }

  stage_names = [
    "RELEASE",
    "TEST",
  ]
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_api_gateway_api&spm=docs.r.api_gateway_api.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the api gateway api. Defaults to null.
* `group_id` - (Required, ForceNew) The api gateway that the api belongs to. Defaults to null.
* `description` - (Required) The description of the api. Defaults to null.
* `auth_type` - (Required) The authorization Type including APP and ANONYMOUS. Defaults to null.
* `request_config` - (Required, List) Request_config defines how users can send requests to your API. See [`request_config`](#request_config) below.
* `service_type` - (Required) The type of backend service. Type including HTTP, VPC, FunctionCompute and MOCK. Defaults to null.
* `http_service_config` - (Optional, List) http_service_config defines the config when service_type selected 'HTTP'. See [`http_service_config`](#http_service_config) below.
* `http_vpc_service_config` - (Optional, List) http_vpc_service_config defines the config when service_type selected 'HTTP-VPC'. See [`http_vpc_service_config`](#http_vpc_service_config) below.
* `fc_service_config` - (Optional, List) fc_service_config defines the config when service_type selected 'FunctionCompute'. See [`fc_service_config`](#fc_service_config) below.
* `mock_service_config` - (Optional, List) http_service_config defines the config when service_type selected 'MOCK'. See [`mock_service_config`](#mock_service_config) below.
* `request_parameters` - (Optional, List) request_parameters defines the request parameters of the api. See [`request_parameters`](#request_parameters) below.
* `constant_parameters` - (Optional, List) constant_parameters defines the constant parameters of the api. See [`constant_parameters`](#constant_parameters) below.
* `system_parameters` - (Optional, List) system_parameters defines the system parameters of the api. See [`system_parameters`](#system_parameters) below.
* `stage_names` - (Optional, Type: list) Stages that the api need to be deployed. Valid value: `RELEASE`,`PRE`,`TEST`.
* `force_nonce_check` - (Optional, Type: bool, Available in v1.140+) Whether to prevent API replay attack. Default value: `false`.

### `request_config`

The request_config mapping supports the following:

* `protocol` - (Required) The protocol of api which supports values of 'HTTP','HTTPS' or 'HTTP,HTTPS'.
* `method` - (Required) The method of the api, including 'GET','POST','PUT' etc.
* `path` - (Required) The request path of the api.
* `mode` - (Required) The mode of the parameters between request parameters and service parameters, which support the values of 'MAPPING' and 'PASSTHROUGH'.
* `body_format` - (Optional) The body format of the api, which support the values of 'STREAM' and 'FORM'.

### `http_service_config`

The http_service_config mapping supports the following:

* `address` - (Required) The address of backend service.
* `path` - (Required) The path of backend service.
* `method` - (Required) The http method of backend service.
* `timeout` - (Required) Backend service time-out time; unit: millisecond.
* `aone_name` - (Optional) The name of aone.
* `content_type_category` - (Optional, Available since v1.228.0) The content type category of backend service which supports values of 'DEFAULT','CUSTOM' and 'CLIENT'.
* `content_type_value` - (Optional, Available since v1.228.0) The content type value of backend service.

### `http_vpc_service_config`

The http_vpc_service_config mapping supports the following:

* `name` - (Required) The name of vpc instance.
* `path` - (Required) The path of backend service.
* `method` - (Required) The http method of backend service.
* `timeout` - (Required) Backend service time-out time. Unit: millisecond.
* `aone_name` - (Optional) The name of aone.
* `vpc_scheme` - (Optional, Available since v1.228.0) The vpc scheme of backend service which supports values of `HTTP` and `HTTPS`.
* `content_type_category` - (Optional, Available since v1.228.0) The content type category of backend service which supports values of 'DEFAULT','CUSTOM' and 'CLIENT'.
* `content_type_value` - (Optional, Available since v1.228.0) The content type value of backend service.

### `fc_service_config`

The fc_service_config mapping supports the following:

* `function_version` - (Optional) The function compute version of function compute service. Supports values of `2.0`, `3.0`. Default value: `2.0`.
* `function_type` - (Optional, Available in v1.219.0) The type of function compute service. Supports values of `FCEvent`,`HttpTrigger`. Default value: `FCEvent`.
* `region` - (Required) The region that the function compute service belongs to.
* `function_name` - (Optional) The function name of function compute service. Required if `function_type` is `FCEvent`.
* `service_name` - (Optional) The service name of function compute service. Required if `function_type` is `FCEvent` and `function_version` is `2.0`.
* `qualifier` - (Optional, Available in v1.219.0) The qualifier of function name of compute service.
* `function_base_url` - (Optional, Available in v1.219.0) The base url of function compute service. Required if `function_type` is `HttpTrigger`.
* `path` - (Optional, Available in v1.219.0) The path of function compute service. Required if `function_type` is `HttpTrigger`.
* `method` - (Optional, Available in v1.219.0) The http method of function compute service. Required if `function_type` is `HttpTrigger`.
* `only_business_path` - (Optional, Available in v1.219.0) Whether to filter path in `function_base_url`. Optional if `function_type` is `HttpTrigger`.
* `arn_role` - (Required) RAM role arn attached to the Function Compute service. This governs both who / what can invoke your Function, as well as what resources our Function has access to. See [User Permissions](https://www.alibabacloud.com/help/doc-detail/52885.htm) for more details.
* `timeout` - (Required) Backend service time-out time; unit: millisecond.

### `mock_service_config`

The mock_service_config mapping supports the following:

* `result` - (Required) The result of the mock service.
* `aone_name` - (Optional) The name of aone.

### `request_parameters`

The request_parameters mapping supports the following:

* `name` - (Required) Request's parameter name.
* `type` - (Required) Parameter type which supports values of 'STRING','INT','BOOLEAN','LONG',"FLOAT" and "DOUBLE".
* `required` - (Required) Parameter required or not; values: REQUIRED and OPTIONAL.
* `in` - (Required) Request's parameter location; values: BODY, HEAD, QUERY, and PATH.
* `in_service` - (Required) Backend service's parameter location; values: BODY, HEAD, QUERY, and PATH.
* `name_service` - (Required) Backend service's parameter name.
* `description` - (Optional) The description of parameter.
* `default_value` - (Optional) The default value of the parameter.

### `constant_parameters`

The constant_parameters mapping supports the following:

* `name` - (Required) Constant parameter name.
* `in` - (Required) Constant parameter location; values: 'HEAD' and 'QUERY'.
* `value` - (Required) Constant parameter value.
* `description` - (Optional) The description of Constant parameter.

### `system_parameters`

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

```shell
$ terraform import alicloud_api_gateway_api.example "ab2351f2ce904edaa8d92a0510832b91:e4f728fca5a94148b023b99a3e5d0b62"
```
