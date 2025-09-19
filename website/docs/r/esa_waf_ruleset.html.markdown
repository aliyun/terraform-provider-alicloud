---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_waf_ruleset"
description: |-
  Provides a Alicloud ESA Waf Ruleset resource.
---

# alicloud_esa_waf_ruleset

Provides a ESA Waf Ruleset resource.

waf rule set.

For information about ESA Waf Ruleset and how to use it, see [What is Waf Ruleset](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateWafRuleset).

-> **NOTE:** Available since v1.260.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
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
```

## Argument Reference

The following arguments are supported:
* `name` - (Optional, ForceNew) The ruleset name.
* `phase` - (Required, ForceNew) The WAF phase
* `site_id` - (Required, ForceNew, Int) The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.
* `site_version` - (Optional, Int) The site version.
* `status` - (Optional, Computed) Rule Set Status

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<ruleset_id>:<site_id>`.
* `ruleset_id` - waf rule set id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Waf Ruleset.
* `delete` - (Defaults to 5 mins) Used when delete the Waf Ruleset.
* `update` - (Defaults to 5 mins) Used when update the Waf Ruleset.

## Import

ESA Waf Ruleset can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_waf_ruleset.example <ruleset_id>:<site_id>
```