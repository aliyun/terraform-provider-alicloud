---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_compression_rule"
description: |-
  Provides a Alicloud ESA Compression Rule resource.
---

# alicloud_esa_compression_rule

Provides a ESA Compression Rule resource.



For information about ESA Compression Rule and how to use it, see [What is Compression Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateCompressionRule).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_compression_rule&exampleId=f6a0eaa3-19a6-a29f-028e-efe54ad7c3f445f6a674&activeTab=example&spm=docs.r.esa_compression_rule.0.f6a0eaa319&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "example" {
  site_name   = "compression.example.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "domestic"
  access_type = "NS"
}

resource "alicloud_esa_compression_rule" "default" {
  gzip         = "off"
  brotli       = "off"
  rule         = "http.host eq \"video.example.com\""
  site_version = "0"
  rule_name    = "rule_example"
  site_id      = alicloud_esa_site.example.id
  zstd         = "off"
  rule_enable  = "off"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_compression_rule&spm=docs.r.esa_compression_rule.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `brotli` - (Optional) Brotli compression. Value range:
  - `on`: on.
  - `off`: off.
* `gzip` - (Optional) Gzip compression. Value range:
  - `on`: on.
  - `off`: off.
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - `on`: open.
  - `off`: close
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Computed, Int, Available since v1.262.0) The rule execution order prioritizes lower numerical values. It is only applicable when setting or modifying the order of individual rule configurations.
* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.
* `zstd` - (Optional) Zstd compression. Value range:
  - `on`: on.
  - `off`: off.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Compression Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Compression Rule.
* `update` - (Defaults to 5 mins) Used when update the Compression Rule.

## Import

ESA Compression Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_compression_rule.example <site_id>:<config_id>
```