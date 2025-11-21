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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_waf_ruleset&exampleId=7a144afa-4769-0103-01a9-8ad20499c71a6c720991&activeTab=example&spm=docs.r.esa_waf_ruleset.0.7a144afa47&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `name` - (Optional, ForceNew, Computed) The ruleset name.
* `phase` - (Required, ForceNew) The WAF phase
* `site_id` - (Required, ForceNew) The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.
* `site_version` - (Optional, Int) The site version.

-> **NOTE:** This parameter only applies during resource creation, update or deletion. If modified in isolation without other property changes, Terraform will not trigger any action.

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