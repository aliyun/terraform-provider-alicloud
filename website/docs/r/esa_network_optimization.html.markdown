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
* `grpc` - (Optional) Whether to enable GRPC, which is disabled by default. Value range:
  - on: Open
  - off: off
* `http2_origin` - (Optional) Whether to enable HTTP2 back to the source, the default is off. Value range:
  - on: Open
  - off: off
* `rule` - (Optional) The rule content.
* `rule_enable` - (Optional) Rule switch. Value:
  - on: Open
  - off: off
* `rule_name` - (Optional) Rule name, which can be used to find the rule with the transmitted field as its rule name. This only takes effect when the functionName is provided.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Optional, ForceNew, Int) The version number of the site. For a site with version management enabled, you can use this parameter to specify the effective site version. The default version is 0.
* `smart_routing` - (Optional) Whether to enable the intelligent routing service is disabled by default. Value range:
  - on: Open
  - off: off
* `upload_max_filesize` - (Optional) The maximum upload file size, in MB. The value range is 100 to 500.
* `websocket` - (Optional) Whether to enable Websocket, enabled by default. Value range:
  - on: Open
  - off: off

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

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