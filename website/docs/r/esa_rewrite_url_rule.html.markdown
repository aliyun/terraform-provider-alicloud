---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_rewrite_url_rule"
description: |-
  Provides a Alicloud ESA Rewrite Url Rule resource.
---

# alicloud_esa_rewrite_url_rule

Provides a ESA Rewrite Url Rule resource.



For information about ESA Rewrite Url Rule and how to use it, see [What is Rewrite Url Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateRewriteUrlRule).

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

resource "alicloud_esa_rate_plan_instance" "resource_RewriteUrlRule_RatePlanInstance_example" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_RewriteUrlRule_Site_example" {
  site_name   = "gositecdn.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_RewriteUrlRule_RatePlanInstance_example.id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_rewrite_url_rule" "default" {
  rewrite_uri_type          = "static"
  rewrite_query_string_type = "static"
  site_id                   = alicloud_esa_site.resource_RewriteUrlRule_Site_example.id
  rule_name                 = "example"
  rule_enable               = "on"
  query_string              = "example=123"
  site_version              = "0"
  rule                      = "http.host eq \"video.example.com\""
  uri                       = "/image/example.jpg"
}
```

## Argument Reference

The following arguments are supported:
* `query_string` - (Optional) The desired query string to which you want to rewrite the query string in the original request.
* `rewrite_query_string_type` - (Optional) The query string rewrite method. Valid value:

  - static
* `rewrite_uri_type` - (Optional) The path rewrite method. Valid value:

  - static
* `rule` - (Optional) The rule content.
* `rule_enable` - (Optional) Indicates whether the rule is enabled. Valid values:

  - on
  - off
* `rule_name` - (Optional) Rule name. You can find the rule whose field is passed by the rule name. The rule takes effect only if functionName is passed.
* `site_id` - (Required, ForceNew, Int) The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.
      
* `site_version` - (Optional, ForceNew, Int) The version number of the website configurations.
* `uri` - (Optional) The desired URI to which you want to rewrite the path in the original request.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - ConfigId

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Rewrite Url Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Rewrite Url Rule.
* `update` - (Defaults to 5 mins) Used when update the Rewrite Url Rule.

## Import

ESA Rewrite Url Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_rewrite_url_rule.example <site_id>:<config_id>
```