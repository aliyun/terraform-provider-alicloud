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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_origin_rule&exampleId=c5336247-6e35-1d85-2621-f1d3dac196e4f636f38c&activeTab=example&spm=docs.r.esa_origin_rule.0.c53362476e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - on: open.
  - off: close.
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
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