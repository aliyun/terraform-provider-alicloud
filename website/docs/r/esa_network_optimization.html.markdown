---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_network_optimization"
description: |-
  Provides a Alicloud ESA Network Optimization resource.
---

# alicloud_esa_network_optimization

Provides a ESA Network Optimization resource.



For information about ESA Network Optimization and how to use it, see [What is Network Optimization](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateNetworkOptimization).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_network_optimization&exampleId=eaf35a11-00a4-5fe2-c66e-24d5e3e0180f1ccb1b6d&activeTab=example&spm=docs.r.esa_network_optimization.0.eaf35a1100&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_esa_site" "default" {
  site_name   = "gositecdn-${random_integer.default.result}.cn"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_network_optimization" "default" {
  site_version        = "0"
  site_id             = alicloud_esa_site.default.id
  rule_enable         = "on"
  websocket           = "off"
  rule                = "(http.host eq \"tf.example.com\")"
  grpc                = "off"
  http2_origin        = "off"
  smart_routing       = "off"
  upload_max_filesize = "100"
  rule_name           = "network_optimization"
}
```

## Argument Reference

The following arguments are supported:
* `grpc` - (Optional) Whether to enable GRPC, default is disabled. Value range:
  - `on`: Enabled
  - `off`: Disabled
* `http2_origin` - (Optional) Whether to enable HTTP2 origin, default is disabled. Value range:
  - `on`: Enabled
  - `off`: Disabled
* `rule` - (Optional) Rule content.
* `rule_enable` - (Optional) Rule switch. Values:
  - `on`: Enabled
  - `off`: Disabled
* `rule_name` - (Optional) Rule name.
* `site_id` - (Required, ForceNew, Int) Site ID.
* `site_version` - (Optional, ForceNew, Int) Site version number.
* `smart_routing` - (Optional) Whether to enable smart routing service, default is disabled. Value range:
  - `on`: Enabled
  - `off`: Disabled
* `upload_max_filesize` - (Optional) Maximum upload file size, in MB, value range: 100～500.
* `websocket` - (Optional) Whether to enable Websocket, default is enabled. Value range:
  - `on`: Enabled
  - `off`: Disabled

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - ConfigId of the configuration, which can be obtained by calling the ListNetworkOptimizations.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Optimization.
* `delete` - (Defaults to 5 mins) Used when delete the Network Optimization.
* `update` - (Defaults to 5 mins) Used when update the Network Optimization.

## Import

ESA Network Optimization can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_network_optimization.example <site_id>:<config_id>
```