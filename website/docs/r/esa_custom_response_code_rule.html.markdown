---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_custom_response_code_rule"
description: |-
  Provides a Alicloud ESA Custom Response Code Rule resource.
---

# alicloud_esa_custom_response_code_rule

Provides a ESA Custom Response Code Rule resource.



For information about ESA Custom Response Code Rule and how to use it, see [What is Custom Response Code Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateCustomResponseCodeRule).

-> **NOTE:** Available since v1.276.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_CustomResponseCodeRule_example" {
  type         = "NS"
  auto_renew   = false
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "basic"
}

resource "alicloud_esa_site" "resource_Site_CustomResponseCodeRule_example" {
  site_name   = "hyhexample.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_CustomResponseCodeRule_example.id
  coverage    = "overseas"
  access_type = "NS"
}


resource "alicloud_esa_custom_response_code_rule" "default" {
  page_id      = "0"
  site_id      = alicloud_esa_site.resource_Site_CustomResponseCodeRule_example.id
  return_code  = "400"
  rule_enable  = "on"
  rule         = "(http.host eq \"video.example.com\")"
  sequence     = "1"
  site_version = "0"
  rule_name    = var.name
}
```

## Argument Reference

The following arguments are supported:
* `page_id` - (Required) Response page.
* `return_code` - (Required) The response code.
* `rule` - (Optional) The content of the rule. A conditional expression is used to match a user request. You do not need to set this parameter when you add global configurations. Use cases:

  - Match all incoming requests: Set the value to true.
  - Set the value to a custom expression, for example, (http.host eq "video.example.com"): Match the specified request.
* `rule_enable` - (Optional) Specifies whether to enable the rule. Valid values: You do not need to set this parameter when you add global configurations. Valid values:

  - on
  - off
* `rule_name` - (Optional) The rule name.
* `sequence` - (Optional, Computed, Int) The order in which the rule is executed. A smaller value gives priority to the rule.
* `site_id` - (Required, ForceNew) The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.
* `site_version` - (Optional, ForceNew, Int) The version number of the website configurations. You can use this parameter to specify a version of your website to apply the feature settings. By default, version 0 is used.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<site_id>:<config_id>`.
* `config_id` - The ID of the configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Custom Response Code Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Response Code Rule.
* `update` - (Defaults to 5 mins) Used when update the Custom Response Code Rule.

## Import

ESA Custom Response Code Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_custom_response_code_rule.example <site_id>:<config_id>
```
