---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_plugin"
description: |-
  Provides a Alicloud Api Gateway Plugin resource.
---

# alicloud_api_gateway_plugin

Provides a Api Gateway Plugin resource. 

For information about Api Gateway Plugin and how to use it, see [What is Plugin](https://www.alibabacloud.com/help/en/api-gateway/developer-reference/api-cloudapi-2016-07-14-createplugin).

-> **NOTE:** Available since v1.187.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_plugin&exampleId=0459c282-24f0-8aa1-1585-b60a576f4cbc7a5bd10e&activeTab=example&spm=docs.r.api_gateway_plugin.0.0459c28224&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_api_gateway_plugin" "default" {
  description = var.name
  plugin_name = var.name
  plugin_data = jsonencode({
    "routes" : [
      {
        "name" : "Vip",
        "condition" : "$CaAppId = 123456",
        "backend" : {
          "type" : "HTTP-VPC",
          "vpcAccessName" : "slbAccessForVip"
        }
      },
      {
        "name" : "MockForOldClient",
        "condition" : "$ClientVersion < '2.0.5'",
        "backend" : {
          "type" : "MOCK",
          "statusCode" : 400,
          "mockBody" : "This version is not supported!!!"
        }
      },
      {
        "name" : "BlueGreenPercent05",
        "condition" : "1 = 1",
        "backend" : {
          "type" : "HTTP",
          "address" : "https://beta-version.api.foo.com"
        },
        "constant-parameters" : [
          {
            "name" : "x-route-blue-green",
            "location" : "header",
            "value" : "route-blue-green"
          }
        ]
      }
    ]
  })
  plugin_type = "routing"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the plug-in, which cannot exceed 200 characters.
* `plugin_data` - (Required) The definition statement of the plug-in. Plug-in definition statements in the JSON and YAML formats are supported.
* `plugin_name` - (Required) The name of the plug-in that you want to create. It can contain uppercase English letters, lowercase English letters, Chinese characters, numbers, and underscores (_). It must be 4 to 50 characters in length and cannot start with an underscore (_).
* `plugin_type` - (Required, ForceNew) The type of the plug-in. Valid values:
  - "trafficControl"
  - "ipControl"
  - "backendSignature"
  - "jwtAuth"
  - "basicAuth"
  - "cors"
  - "caching"
  - "routing"
  - "accessControl"
  - "errorMapping"
  - "circuitBreaker"
  - "remoteAuth"
  - "logMask"
  - "transformer".
* `tags` - (Optional, Map) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Create time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Plugin.
* `delete` - (Defaults to 5 mins) Used when delete the Plugin.
* `update` - (Defaults to 5 mins) Used when update the Plugin.

## Import

Api Gateway Plugin can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_plugin.example <id>
```