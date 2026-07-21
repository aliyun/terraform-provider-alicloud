---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_waf_rulesets"
description: |-
  Provides a list of ESA Waf Rulesets to the user.
---

# alicloud_esa_waf_rulesets

This data source provides the ESA Waf Rulesets of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.274.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_waf_ruleset" "default" {
  site_id      = data.alicloud_esa_sites.default.sites.0.site_id
  phase        = "http_custom"
  site_version = "0"
  name         = var.name
}

data "alicloud_esa_waf_rulesets" "ids" {
  ids          = [alicloud_esa_waf_ruleset.default.id]
  site_id      = alicloud_esa_waf_ruleset.default.site_id
  phase        = alicloud_esa_waf_ruleset.default.phase
  site_version = alicloud_esa_waf_ruleset.default.site_version
}

output "esa_waf_rulesets_id_0" {
  value = data.alicloud_esa_waf_rulesets.ids.sets.0.id
}
```

## Argument Reference

The following attributes are exported:

* `ids` - (Optional, List) A list of Waf Ruleset IDs.
* `name_regex` - (Optional) A regex string to filter results by Waf Ruleset name.
* `site_id` - (Required) The ID of the Site.
* `phase` - (Required) The WAF operation phase.
* `site_version` - (Required) The version of the Site.
* `status` - (Optional) The status of the rule set. Valid values: `on`, `off`.
* `query_args` - (Optional, Set) The query parameters. See [`query_args`](#query_args) below.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

### `query_args`

The query_args supports the following:

* `any_like` - (Optional) The fuzzy search for rule set ID, rule set name, rule ID, and rule name.
* `name_like` - (Optional) The fuzzy search for rule set name.
* `order_by` - (Optional) Specify the column to sort by.
* `desc` - (Optional, Bool) Whether to sort in descending order. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Waf Ruleset names.
* `sets` -  A list of Waf Rulesets. Each element contains the following attributes:
  * `id` - The ID of the WAF Rule Set.
  * `ruleset_id` - The ID of the WAF rule set.
  * `phase` - The WAF operation phase.
  * `name` - The name of the rule set.
  * `target` - Protection target type in http_bot.
  * `status` - The status of the rule set.
  * `update_time` - The last modification time of the rule set.
  * `types` - The list of rule types.
  * `fields` - The list of match objects.
