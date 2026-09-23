---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_group_plugin_attachment"
sidebar_current: "docs-alicloud-resource-api-gateway-group-plugin-attachment"
description: |-
  Provides a Alicloud Api Gateway Group Plugin Attachment Resource.
---

# alicloud_api_gateway_group_plugin_attachment

Provides a plugin attachment resource.It is used for attaching a specific plugin to an api group.

For information about Api Gateway Plugin attachment and how to use it, see [Attach Plugin to specified API GROUP](https://www.alibabacloud.com/help/en/api-gateway/traditional-api-gateway/developer-reference/api-cloudapi-2016-07-14-attachgroupplugin)

-> **NOTE:** Available since v1.278.0.

-> **NOTE:** Terraform will auto build plugin attachment while it uses `alicloud_api_gateway_group_plugin_attachment` to build.

## Example Usage

Basic Usage


<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_group_plugin_attachment&exampleId=ef0412d5-71cd-9485-98aa-fe053e23faa9bdf11f41&activeTab=example&spm=docs.r.api_gateway_group_plugin_attachment.0.ef0412d571&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_api_gateway_group" "example" {
  name        = "tf-example-api-gateway-group"
  description = "tf-example-api-gateway-group"
}

resource "alicloud_api_gateway_plugin" "example" {
  description = "tf_example"
  plugin_name = "tf-example-api-gateway-plugin"
  plugin_data = jsonencode({ "allowOrigins" : "api.foo.com", "allowMethods" : "GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH", "allowHeaders" : "Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid", "exposeHeaders" : "Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message", "maxAge" : 172800, "allowCredentials" : true })
  plugin_type = "cors"
}

resource "alicloud_api_gateway_group_plugin_attachment" "example" {
  group_id   = alicloud_api_gateway_group.example.id
  plugin_id  = alicloud_api_gateway_plugin.example.id
  stage_name = "RELEASE"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_api_gateway_group_plugin_attachment&spm=docs.r.api_gateway_group_plugin_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `group_id` - (Required, ForceNew) The group that plugin attaches to.
* `plugin_id` - (Required, ForceNew) The plugin that attaches to the group.
* `stage_name` - (Required, ForceNew) Stage that the plugin attaches to.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID of the group plugin attachment. The value formats as `<group_id>:<plugin_id>:<stage_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the api gateway group plugin attachment.
* `delete` - (Defaults to 5 mins) Used when deleting the api gateway group plugin attachment.

## Import
Api Gateway group plugin attachment a can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_group_plugin_attachment.example <group_id>:<plugin_id>:<stage_name>
```
