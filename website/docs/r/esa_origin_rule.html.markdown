---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_origin_rule"
description: |-
  Provides a Alicloud ESA Origin Rule resource.
---

# alicloud_esa_origin_rule

Provides a ESA Origin Rule resource.



For information about ESA Origin Rule and how to use it, see [What is Origin Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateOriginRule).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_origin_rule" "default" {
  origin_sni        = "origin.example.com"
  site_id           = data.alicloud_esa_sites.default.sites.0.id
  origin_host       = "origin.example.com"
  dns_record        = "tf.example.com"
  site_version      = "0"
  rule_name         = "tf"
  origin_https_port = "443"
  origin_scheme     = "http"
  range             = "on"
  origin_http_port  = "8080"
  rule              = "(http.host eq \"video.example.com\")"
  rule_enable       = "on"
}
```

## Argument Reference

The following arguments are supported:
* `dns_record` - (Optional) Overwrite the DNS resolution record of the origin request.
* `follow302_enable` - (Optional, Available since v1.263.0) Return Source 302 follow switch. Value range:
  - `on`: ON.
  - `off`: closed.
* `follow302_max_tries` - (Optional, Available since v1.263.0) 302 follows the upper limit of the number of times, with a value range of [1-5].
* `follow302_retain_args` - (Optional, Available since v1.263.0) Retain the original request parameter switch. Value range:
  - `on`: ON.
  - `off`: closed.
* `follow302_retain_header` - (Optional, Available since v1.263.0) Retain the original request header switch. Value range:
  - `on`: ON.
  - `off`: closed.
* `follow302_target_host` - (Optional, Available since v1.263.0) Modify the source host after 302.
* `origin_host` - (Optional) The HOST carried in the back-to-origin request.
* `origin_http_port` - (Optional) The port of the origin station accessed when the HTTP protocol is used to return to the origin.
* `origin_https_port` - (Optional) The port of the origin station accessed when the HTTPS protocol is used to return to the origin.
* `origin_mtls` - (Optional, Available since v1.263.0) The mtls switch. Value range:
  - `on`: ON.
  - `off`: closed.
* `origin_read_timeout` - (Optional, Available since v1.263.0) Read timeout interval of the source station (s).
* `origin_scheme` - (Optional) The protocol used by the back-to-origin request. Value range:
  - `http`: uses the http protocol to return to the source.
  - `https`: uses the https protocol to return to the source.
  - `follow`: follows the Client Protocol back to the source.
* `origin_sni` - (Optional) SNI carried in the back-to-origin request.
* `origin_verify` - (Optional, Available since v1.263.0) Source station certificate verification switch. Value range:
  - `on`: ON.
  - `off`: closed.
* `range` - (Optional) Use the range sharding method to download the file from the source. Value range:
  - `on`: Open.
  - `off`: off.
  - `force`: force.
* `range_chunk_size` - (Optional, Available since v1.263.0) range shard size.
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - `on`: open.
  - `off`: close.
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Int, Available since v1.263.0) The rule execution order prioritizes lower numerical values. It is only applicable when setting or modifying the order of individual rule configurations.
* `site_id` - (Required, ForceNew, Int) The site ID.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Back-to-source rule configuration ID

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Origin Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Origin Rule.
* `update` - (Defaults to 5 mins) Used when update the Origin Rule.

## Import

ESA Origin Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_origin_rule.example <site_id>:<config_id>
```