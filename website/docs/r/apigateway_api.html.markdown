---
layout: "alicloud"
page_title: "Alicloud: alicloud_apigateway_api"
sidebar_current: "docs-alicloud-resource-apigateway-api"
description: |-
  Provides a Alicloud Apigateway api resource.
---

# alicloud_apigateway_api

Provides an api resource.

~> **NOTE:** Terraform will auto build api while it uses `alicloud_apigateway_api` to build api.

## Example Usage

Basic Usage

```
resource "alicloud_apigateway_api_group" "apiGroup" {
  name = "apiGroup"
  description = "Description of api group"
}

resource "alicloud_apigateway_api" "api" {
  name = "api"
  description = "description of api"
  group_id = "${alicloud_apigateway_api_group.apiGroup.id}"
  auth_type = "APP"
  request_config = "{\"RequestProtocol\":\"HTTP\",\"RequestHttpMethod\":\"GET\",\"RequestPath\":\"/test\",\"BodyFormat\":\"\",\"PostBodyDescription\":\"\",\"RequestMode\":\"MAPPING\"}"
  service_config = "{\"ServiceProtocol\":\"HTTP\",\"ServiceHttpMethod\":\"GET\",\"ServiceAddress\":\"http://apigateway-backend.alicloudapi.com:8080\",\"ServiceTimeout\":\"10000\",\"ServicePath\":\"/web/cloudapi\",\"Mock\":\"FALSE\",\"MockResult\":\"\",\"ServiceVpcEnable\":\"FALSE\",\"VpcConfig\":{},\"FunctionComputeConfig\":{},\"ContentTypeCatagory\":\"DEFAULT\",\"ContentTypeValue\":\"application/x-www-form-urlencoded; charset=UTF-8\",\"AoneAppName\":\"cloudapi-openapi\"}"
  request_parameters = "[{\"ParameterLocation\":{\"name\":\"Head\",\"orderNumber\":2},\"ParameterType\":\"String\",\"Required\":\"OPTIONAL\",\"isHide\":false,\"ApiParameterName\":\"requestHead\",\"Location\":\"Head\"}]"
  service_parameters = "[{\"ServiceParameterName\":\"requestHead\",\"Location\":\"Query\",\"Type\":\"String\",\"ParameterCatalog\":\"REQUEST\"}]"
  service_parameters_map = "[{\"ServiceParameterName\":\"requestHead\",\"RequestParameterName\":\"requestHead\"}]"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the api, which is unique in the group. Names must be a string of 4–50 characters, which can contain Chinese characters, English letters, numbers, and underscores, and starts with an English letter or Chinese character.
* `description` - (Required, Forces new resource) Description on the API, at most 180 characters.
* `group_id` - (Required, Forces new resource) ID of the specified group.
* `auth_type` - (Required, Forces new resource) API security authorization type; values: 
	- APP: Only authorized apps can call the API.
	- ANONYMOUS: The API can be called anonymously.
* `request_config` - (Required, Forces new resource) Configuration items of the API request sent by a consumer to the apigateway. The value is a json object based on [Request Config](https://www.alibabacloud.com/help/doc-detail/43985.htm?spm=a2c63.p38356.a3.1.16f460dbdLnSrH) 
* `service_config` - (Required, Forces new resource) Configuration items of the API request sent by the gateway to a backend service.. The value is a json object based on [Service Config](https://www.alibabacloud.com/help/doc-detail/43987.htm?spm=a2c63.p38356.a3.2.16f460dbdLnSrH) 
* `request_parameters` - (Required, Forces new resource) Description on parameters of the API request sent by a consumer to the gateway.. The value is a json object based on [Request Parameters](https://www.alibabacloud.com/help/doc-detail/43986.htm?spm=a2c63.p38356.a3.3.16f460dbdLnSrH) 
* `service_parameters` - (Required, Forces new resource) Description on parameters of the API request sent by the gateway to a backend service. The value is a json object based on [Service Parameters](https://www.alibabacloud.com/help/doc-detail/43988.htm?spm=a2c63.p38356.a3.4.16f460dbZ5cTNp) 
* `service_parameters_map` - (Required, Forces new resource) The mappings between parameters of a request sent by a consumer to the gateway and parameters of a request sent by the gateway to a backend service. The value is a json object based on [Service Parameter Map](https://www.alibabacloud.com/help/doc-detail/43989.htm?spm=a2c63.p38356.a3.5.16f460dbvrGD97) 

## Attributes Reference

The following attributes are exported:

* `api_id` - The ID of the api.