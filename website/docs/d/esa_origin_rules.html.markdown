---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_origin_rules"
description: |-
  Provides a list of Esa Origin Rules to the user.
---

# alicloud_esa_origin_rules

This data source provides the Esa Origin Rules of the current Alibaba Cloud user.

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

resource "alicloud_esa_origin_rule" "default" {
  site_id                 = data.alicloud_esa_sites.default.sites.0.id
  site_version            = 0
  rule_enable             = "on"
  rule_name               = var.name
  rule                    = "true"
  sequence                = 1
  origin_host             = "origin.example.com"
  origin_scheme           = "http"
  origin_sni              = "origin.example.com"
  origin_https_port       = "443"
  origin_http_port        = "8080"
  origin_read_timeout     = "30"
  dns_record              = "test.example.com"
  origin_verify           = "on"
  origin_mtls             = "on"
  follow302_enable        = "on"
  follow302_max_tries     = "3"
  follow302_target_host   = "redirect.example.com"
  follow302_retain_header = "on"
  follow302_retain_args   = "on"
  range                   = "on"
  range_chunk_size        = "1MB"
}

data "alicloud_esa_origin_rules" "ids" {
  ids     = [alicloud_esa_origin_rule.default.id]
  site_id = alicloud_esa_origin_rule.default.site_id
}

output "esa_origin_rules_id_0" {
  value = data.alicloud_esa_origin_rules.ids.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List) A list of Origin Rule IDs. It formats as `<site_id>:<config_id>`.
* `name_regex` - (Optional) A regex string to filter results by Origin Rule name.
* `site_id` - (Required) The ID of the Site.
* `config_id` - (Optional) The ID of the Configuration.
* `config_type` - (Optional) The type of the Configuration. Valid values: `global`, `rule`.
* `rule_name` - (Optional) The name of the rule.
* `site_version` - (Optional) The version of the site.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Origin Rule names.
* `rules` - A list of Origin Rules. Each element contains the following attributes:
  * `id` - The ID of the Origin Rule.
  * `config_id` - The ID of the Configuration.
  * `config_type` - The type of the Configuration.
  * `site_version` - The version of the site.
  * `rule_enable` - Rule switch.
  * `rule_name` - Rule name.
  * `rule` - Rule content.
  * `sequence` - The rule execution order prioritizes lower numerical values.
  * `origin_host` - The Host header carried in the origin request.
  * `origin_scheme` - The protocol used for origin requests.
  * `origin_sni` - The SNI carried in the origin request.
  * `origin_http_port` - The origin server port used for origin requests over HTTP.
  * `origin_https_port` - The origin server port used for origin requests over HTTPS.
  * `origin_read_timeout` - The read timeout, in seconds, for the origin server.
  * `dns_record` - Overrides the DNS record for the origin request.
  * `origin_verify` - Specifies whether to verify the origin server certificate.
  * `origin_mtls` - Specifies whether mTLS is enabled.
  * `follow302_enable` - Specifies whether to follow 302 redirects from the origin.
  * `follow302_max_tries` - The maximum number of 302 redirects to follow.   
  * `follow302_target_host` - The host to use for the origin request after following a 302 redirect.
  * `follow302_retain_header` - Specifies whether to retain the original request header when following a redirect.
  * `follow302_retain_args` - Specifies whether to retain the original request parameters when following a redirect.
  * `range` - Specifies whether to use range-based requests to retrieve files from the origin.
  * `range_chunk_size` - The size of each chunk for range requests.
