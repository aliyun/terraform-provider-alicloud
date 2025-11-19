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
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - `on`: open.
  - `off`: close.
* `rule_name` - (Optional) Rule name.
* `sequence` - (Optional, Computed, Int, Available since v1.262.1) The rule execution order prioritizes lower numerical values. It is only applicable when setting or modifying the order of individual rule configurations.
* `site_id` - (Required, ForceNew) Site ID.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.
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