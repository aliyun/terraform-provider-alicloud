---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_network_optimizations"
description: |-
  Provides a list of Esa Network Optimizations to the user.
---

# alicloud_esa_network_optimizations

This data source provides the Esa Network Optimizations of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.279.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_network_optimization" "default" {
  site_id             = data.alicloud_esa_sites.default.sites.0.id
  site_version        = 0
  rule_enable         = "on"
  rule_name           = var.name
  rule                = "true"
  sequence            = 1
  smart_routing       = "on"
  websocket           = "on"
  http2_origin        = "on"
  grpc                = "on"
  upload_max_filesize = "100"
}

data "alicloud_esa_network_optimizations" "ids" {
  ids     = [alicloud_esa_network_optimization.default.id]
  site_id = alicloud_esa_network_optimization.default.site_id
}

output "esa_network_optimizations_id_0" {
  value = data.alicloud_esa_network_optimizations.ids.optimizations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List) A list of Network Optimization IDs. It formats as `<site_id>:<config_id>`.
* `name_regex` - (Optional) A regex string to filter results by Network Optimization name.
* `site_id` - (Required) The ID of the Site.
* `config_id` - (Optional) The ID of the Configuration.
* `config_type` - (Optional) The type of the Configuration. Valid values: `global`, `rule`.
* `rule_name` - (Optional) The name of the rule.
* `site_version` - (Optional) The version of the site.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Network Optimization names.
* `optimizations` - A list of Network Optimizations. Each element contains the following attributes:
  * `id` - The ID of the Network Optimization.
  * `config_id` - The ID of the Configuration.
  * `config_type` - The type of the Configuration.
  * `site_version` - The version of the site.
  * `rule_enable` - Rule switch.
  * `rule_name` - Rule name.
  * `rule` - Rule content.
  * `sequence` - The rule execution order prioritizes lower numerical values.
  * `smart_routing` - Whether to enable smart routing service.
  * `grpc` - Whether to enable GRPC.
  * `http2_origin` - Whether to enable HTTP2 origin.
  * `websocket` - Whether to enable Websocket.
  * `upload_max_filesize` - Maximum upload file size.
