---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_error_pages_redirect_rules"
description: |-
  Provides a list of Esa Error Pages Redirect Rules to the user.
---

# alicloud_esa_error_pages_redirect_rules

This data source provides the Esa Error Pages Redirect Rules of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_error_pages_redirect_rule" "default" {
  site_id      = data.alicloud_esa_sites.default.sites.0.id
  rule_enable  = "off"
  rule         = "(http.host eq \"video.example.com\")"
  sequence     = 1
  site_version = 0
  rule_name    = var.name
  error_pages_redirect {
    target_url  = "https://example.com/foo/bar"
    status_code = "500"
  }
  error_pages_redirect {
    target_url  = "https://example.com/foo"
    status_code = "400"
  }
}

data "alicloud_esa_error_pages_redirect_rules" "default" {
  site_id   = data.alicloud_esa_sites.default.sites.0.id
  rule_name = var.name
  ids       = [alicloud_esa_error_pages_redirect_rule.default.id]
}

output "rule_id" {
  value = data.alicloud_esa_error_pages_redirect_rules.default.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, ForceNew) The website ID, which can be obtained by calling the ListSites operation.
* `config_id` - (Optional) The configuration ID.
* `config_type` - (Optional) The configuration type. Valid values: `global`, `rule`.
* `rule_name` - (Optional) The rule name.
* `site_version` - (Optional) The version number of the site configuration.
* `ids` - (Optional) A list of rule IDs.
* `name_regex` - (Optional) A regex string to filter results by the rule name.
* `output_file` - (Optional) File path where data results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the `ids` and `names` arguments:

* `names` - A list of rule names.
* `rules` - A list of ESA Error Pages Redirect Rules. Each element contains the following attributes:
  * `id` - The ID of the rule, formatted as `<site_id>:<config_id>`.
  * `config_id` - The configuration ID.
  * `config_type` - The configuration type. Valid values: `global`, `rule`.
  * `site_version` - The version number of the site configuration.
  * `rule_enable` - The rule switch. Valid values: `on`, `off`.
  * `rule_name` - The rule name.
  * `rule` - The rule content, using conditional expressions to match user requests.
  * `sequence` - The order of rule execution. The smaller the value, the higher the priority for execution.
  * `error_pages_redirect` - The configurations of error pages redirect.
    * `target_url` - The destination URL after the redirect.
    * `status_code` - The response code that you want to use to indicate URL redirection.
