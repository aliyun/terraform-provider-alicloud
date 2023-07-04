---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_plugin"
sidebar_current: "docs-alicloud-resource-api-gateway-plugin"
description: |-
  Provides a Alicloud Api Gateway Plugin resource.
---

# alicloud_api_gateway_plugin

Provides a Api Gateway Plugin resource.

For information about Api Gateway Plugin and how to use it, see [What is Plugin](https://www.alibabacloud.com/help/en/api-gateway/latest/create-an-plugin).

-> **NOTE:** Available since v1.187.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_api_gateway_plugin" "default" {
  description = "tf_example"
  plugin_name = "tf_example"
  plugin_data = "{\"allowOrigins\": \"api.foo.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}"
  plugin_type = "cors"
  tags = {
    Created = "TF",
    For     = "example",
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the plug-in, which cannot exceed 200 characters.
* `plugin_data` - (Required) The definition statement of the plug-in. Plug-in definition statements in the JSON and YAML formats are supported.
* `plugin_name` - (Required) The name of the plug-in that you want to create. It can contain uppercase English letters, lowercase English letters, Chinese characters, numbers, and underscores (_). It must be 4 to 50 characters in length and cannot start with an underscore (_).
* `plugin_type` - (Required, ForceNew) The type of the plug-in. Valid values: `backendSignature`, `caching`, `cors`, `ipControl`, `jwtAuth`, `trafficControl`.
  - ipControl: indicates IP address-based access control.
  - trafficControl: indicates throttling.
  - backendSignature: indicates backend signature.
  - jwtAuth: indicates JWT (OpenId Connect).
  - cors: indicates cross-origin resource access (CORS).
  - caching: indicates caching.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Plugin.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Plugin.
* `update` - (Defaults to 1 mins) Used when update the Plugin.
* `delete` - (Defaults to 1 mins) Used when delete the Plugin.


## Import

Api Gateway Plugin can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_plugin.example <id>
```