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

resource "alicloud_esa_site" "default" {
  site_name   = "gositecdn.cn"
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
* `grpc` - (Optional) Indicates whether to enable GRPC, disabled by default. Possible values:
  - on: Enable
  - off: Disable
* `http2_origin` - (Optional) Indicates whether to enable HTTP2 origin, disabled by default. Possible values:
  - on: Enable
  - off: Disable
* `rule` - (Optional) Rule content.
* `rule_enable` - (Optional) Rule switch. Possible values:
  - on: Enable
  - off: Disable
* `rule_name` - (Optional) Rule name, which can be used to find the rule with the specified name.
* `site_id` - (Required, ForceNew, Int) Site ID, which can be obtained by calling the [ListSites](~~ListSites~~) interface.
* `site_version` - (Optional, ForceNew, Int) Version number of the site configuration. For sites with version management enabled, this parameter specifies the version to which the configuration applies, defaulting to version 0.
* `smart_routing` - (Optional) Indicates whether to enable smart routing service, disabled by default. Possible values:
  - on: Enable
  - off: Disable
* `upload_max_filesize` - (Optional) Maximum upload file size, in MB, value range: 100～500.
* `websocket` - (Optional) Indicates whether to enable Websocket, enabled by default. Possible values:
  - on: Enable
  - off: Disable

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - ConfigId of the configuration, which can be obtained by calling the ListNetworkOptimizations.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Optimization.
* `delete` - (Defaults to 5 mins) Used when delete the Network Optimization.
* `update` - (Defaults to 5 mins) Used when update the Network Optimization.

## Import

ESA Network Optimization can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_network_optimization.example <site_id>:<config_id>
```