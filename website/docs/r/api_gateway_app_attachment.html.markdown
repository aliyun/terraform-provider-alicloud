---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_app_attachment"
sidebar_current: "docs-alicloud-resource-api-gateway-app-attachment"
description: |-
  Provides a Alicloud Api Gateway App Attachment Resource.
---

# alicloud_api_gateway_app_attachment

Provides an app attachment resource.It is used for authorizing a specific api to an app accessing. 

For information about Api Gateway App attachment and how to use it, see [Add specified API access authorities](https://www.alibabacloud.com/help/en/api-gateway/latest/api-cloudapi-2016-07-14-setappsauthorities)

-> **NOTE:** Available since v1.23.0.

-> **NOTE:** Terraform will auto build app attachment while it uses `alicloud_api_gateway_app_attachment` to build.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_app_attachment&exampleId=b3521212-61d9-9101-54f5-70e09d2e2585464e5839&activeTab=example&spm=docs.r.api_gateway_app_attachment.0.b352121261&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-beijing"
}

variable "name" {
  default = "terraform_example"
}
resource "alicloud_api_gateway_group" "example" {
  name        = var.name
  description = var.name
}

resource "alicloud_api_gateway_api" "example" {
  group_id          = alicloud_api_gateway_group.example.id
  name              = var.name
  description       = var.name
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

resource "alicloud_api_gateway_app" "example" {
  name        = var.name
  description = var.name
}

resource "alicloud_api_gateway_app_attachment" "example" {
  api_id     = alicloud_api_gateway_api.example.api_id
  group_id   = alicloud_api_gateway_group.example.id
  app_id     = alicloud_api_gateway_app.example.id
  stage_name = "PRE"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required, ForceNew) The api_id that app apply to access.
* `group_id` - (Required, ForceNew) The group that the api belongs to.
* `app_id` - (Required, ForceNew) The app that apply to the authorization.
* `stage_name` - (Required, ForceNew) Stage that the app apply to access.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the app attachment of api gateway., formatted as `<group_id>:<api_id>:<app_id>:<stage_name>`.
