---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_redirect_rule"
description: |-
  Provides a Alicloud ESA Redirect Rule resource.
---

# alicloud_esa_redirect_rule

Provides a ESA Redirect Rule resource.



For information about ESA Redirect Rule and how to use it, see [What is Redirect Rule](https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/api-esa-2024-09-10-createredirectrule).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_esa_rate_plan_instance" "resource_RedirectRule_example" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_esa_site" "resource_Site_RedirectRule_example" {
  site_name   = "gositecdn-${random_integer.default.result}.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_RedirectRule_example.id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_redirect_rule" "default" {
  status_code          = "301"
  rule_name            = "example"
  site_id              = alicloud_esa_site.resource_Site_RedirectRule_example.id
  type                 = "static"
  reserve_query_string = "on"
  target_url           = "http://www.exapmle.com/index.html"
  rule_enable          = "on"
  site_version         = "0"
  rule                 = "(http.host eq \"video.example.com\")"
}
```

## Argument Reference

The following arguments are supported:
* `reserve_query_string` - (Required) Indicates whether the feature of retaining the query string is enabled. Valid values:

  - on
  - off
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
● Match all incoming requests: value set to true
● Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:

  - `on`
  - `off`
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Int, Available since v1.262.1) Order of rule execution. The smaller the value, the higher the priority for execution.
* `site_id` - (Required, ForceNew, Int) The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.
* `status_code` - (Required) The response code that you want to use to indicate URL redirection. Valid values:

  - 301
  - 302
  - 303
  - 307
  - 308
* `target_url` - (Required) The destination URL to which requests are redirected.
* `type` - (Required) The redirection type. Value range:
  - static: static mode.
  - dynamic: dynamic mode.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Redirect Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Redirect Rule.
* `update` - (Defaults to 5 mins) Used when update the Redirect Rule.

## Import

ESA Redirect Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_redirect_rule.example <site_id>:<config_id>
```