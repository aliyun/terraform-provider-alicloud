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
* `origin_host` - (Optional) The HOST carried in the back-to-origin request.
* `origin_http_port` - (Optional) The port of the origin station accessed when the HTTP protocol is used to return to the origin.
* `origin_https_port` - (Optional) The port of the origin station accessed when the HTTPS protocol is used to return to the origin.
* `origin_scheme` - (Optional) The protocol used by the back-to-origin request. Value range:
  - `http`: uses the http protocol to return to the source.
  - `https`: uses the https protocol to return to the source.
  - `follow`: follows the Client Protocol back to the source.
* `origin_sni` - (Optional) SNI carried in the back-to-origin request.
* `range` - (Optional) Use the range sharding method to download the file from the source. Value range:
  - `on`: Open.
  - `off`: off.
  - `force`: force.
* `rule` - (Optional) Rule Content.
* `rule_enable` - (Optional) Rule switch. Value range:
  - `on`: Open.
  - `off`: off.
* `rule_name` - (Optional) Rule Name.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Optional, ForceNew, Int) Version number of the site.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Origin Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Origin Rule.
* `update` - (Defaults to 5 mins) Used when update the Origin Rule.

## Import

ESA Origin Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_origin_rule.example <site_id>:<config_id>
```