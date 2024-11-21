---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_plugin_attachment"
sidebar_current: "docs-alicloud-resource-api-gateway-plugin-attachment"
description: |-
  Provides a Alicloud Api Gateway Plugin Attachment Resource.
---

# alicloud_api_gateway_plugin_attachment

Provides a plugin attachment resource.It is used for attaching a specific plugin to an api. 

For information about Api Gateway Plugin attachment and how to use it, see [Attach Plugin to specified API](https://www.alibabacloud.com/help/en/api-gateway/developer-reference/api-cloudapi-2016-07-14-attachplugin)

-> **NOTE:** Available since v1.219.0.

-> **NOTE:** Terraform will auto build plugin attachment while it uses `alicloud_api_gateway_plugin_attachment` to build.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_plugin_attachment&exampleId=25739863-f177-e95c-908d-cd9ce1d804d9fdb80ee5&activeTab=example&spm=docs.r.api_gateway_plugin_attachment.0.25739863f1&intl_lang=EN_US" target="_blank">
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

resource "alicloud_api_gateway_plugin" "example" {
  description = "tf_example"
  plugin_name = "tf_example"
  plugin_data = jsonencode({ "allowOrigins" : "api.foo.com", "allowMethods" : "GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH", "allowHeaders" : "Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid", "exposeHeaders" : "Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message", "maxAge" : 172800, "allowCredentials" : true })
  plugin_type = "cors"
}

resource "alicloud_api_gateway_plugin_attachment" "example" {
  api_id     = alicloud_api_gateway_api.example.api_id
  group_id   = alicloud_api_gateway_group.example.id
  plugin_id  = alicloud_api_gateway_plugin.example.id
  stage_name = "RELEASE"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required, ForceNew) The api_id that plugin attaches to.
* `group_id` - (Required, ForceNew) The group that the api belongs to.
* `plugin_id` - (Required, ForceNew) The plugin that attaches to the api.
* `stage_name` - (Required, ForceNew) Stage that the plugin attaches to.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the plugin attachment of api gateway., formatted as `<group_id>:<api_id>:<plugin_id>:<stage_name>`.
