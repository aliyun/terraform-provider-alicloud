---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_apis"
sidebar_current: "docs-alicloud-datasource-api-gateway-apis"
description: |-
  Provides a list of Api Gateway APIs to the user.
---

# alicloud_api_gateway_apis

This data source provides the Api Gateway APIs of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.22.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_api_gateway_group" "default" {
  name        = var.name
  description = var.name
}

resource "alicloud_api_gateway_api" "default" {
  group_id     = alicloud_api_gateway_group.default.id
  name         = var.name
  description  = var.name
  auth_type    = "APP"
  service_type = "HTTP"
  request_config {
    protocol = "HTTP"
    method   = "GET"
    path     = "/test/path"
    mode     = "MAPPING"
  }
  http_service_config {
    address   = "http://apigateway-backend.alicloudapi.com:8080"
    method    = "GET"
    path      = "/web/cloudapi"
    timeout   = 20
    aone_name = "cloudapi-openapi"
  }
  request_parameters {
    name         = var.name
    type         = "STRING"
    required     = "OPTIONAL"
    in           = "QUERY"
    in_service   = "QUERY"
    name_service = var.name
  }
}

data "alicloud_api_gateway_apis" "ids" {
  ids = [alicloud_api_gateway_api.default.id]
}

output "api_gateway_apis_id_0" {
  value = data.alicloud_api_gateway_apis.ids.apis.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List, Available since v1.52.2) A list of API IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by API name.
* `group_id` - (Optional, ForceNew) The ID of the API group.
* `api_id` - (Optional, ForceNew) The ID of the API.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of API names.
* `apis` - A list of APIs. Each element contains the following attributes:
  * `id` - The resource ID in terraform of API. It formats as `<group_id>:<api_id>`.
  * `group_id` - The ID of the API group.
  * `api_id` - (Available since v1.224.0) The ID of the API.
  * `name` - The name of the API.
  * `description` - The description of the API.
  * `group_name` - The name of the API group.
  * `region_id` - The region ID of the API.
