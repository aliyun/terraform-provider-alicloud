---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_error_pages_redirect_rule"
description: |-
  Provides a Alicloud ESA Error Pages Redirect Rule resource.
---

# alicloud_esa_error_pages_redirect_rule

Provides a ESA Error Pages Redirect Rule resource.

Error Code Redirection Rule.

For information about ESA Error Pages Redirect Rule and how to use it, see [What is Error Pages Redirect Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateErrorPagesRedirectRule).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_example_ErrorPagesRedirectRule" {
  type         = "NS"
  auto_renew   = false
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_example_ErrorPagesRedirectRule" {
  site_name   = "gositecdn.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_example_ErrorPagesRedirectRule.id
  coverage    = "overseas"
  access_type = "NS"
}


resource "alicloud_esa_error_pages_redirect_rule" "default" {
  site_id      = alicloud_esa_site.resource_Site_example_ErrorPagesRedirectRule.id
  rule_enable  = "off"
  rule         = "(http.host eq \"video.example.com\")"
  sequence     = "1"
  site_version = "0"
  rule_name    = "rule_example"
  error_pages_redirect {
    target_url  = "https://example.com/foo/bar"
    status_code = "500"
  }
  error_pages_redirect {
    target_url  = "https://example.com/foo"
    status_code = "400"
  }
  error_pages_redirect {
    target_url  = "https://example.com/example"
    status_code = "503"
  }
}
```

## Argument Reference

The following arguments are supported:
* `error_pages_redirect` - (Required, List) Error code redirection. See [`error_pages_redirect`](#error_pages_redirect) below.
* `rule` - (Optional) Rule content, which uses conditional expressions to match user requests. This parameter is not required when adding a global configuration. There are two usage scenarios:
  - Match all incoming requests: Set the value to `true`.
  - Match specific requests: Set the value to a custom expression, such as `(http.host eq "video.example.com")`.
* `rule_enable` - (Optional) Rule switch. This parameter is not required when adding a global configuration. Valid values:
  - `on`: Enable.
  - `off`: Disable.
* `rule_name` - (Optional) Rule name. This parameter is not required when adding a global configuration.
* `sequence` - (Optional, ForceNew, Int) Rule execution order. Rules with smaller values take higher priority.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the [ListSites](https://help.aliyun.com/document_detail/2850189.html) operation.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites with configuration version management enabled, you can use this parameter to specify the site version for which the configuration takes effect. The default value is version 0.

### `error_pages_redirect`

The error_pages_redirect supports the following:
* `status_code` - (Required) The response status code used by the node when returning the redirection address to the client. Valid values:
  - 400
  - 403
  - 404
  - 405
  - 414
  - 416
  - 500
  - 501
  - 502
  - 503
  - 504
* `target_url` - (Required) The target URL after redirection.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Configuration ID.
* `config_type` - Configuration type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Error Pages Redirect Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Error Pages Redirect Rule.
* `update` - (Defaults to 5 mins) Used when update the Error Pages Redirect Rule.

## Import

ESA Error Pages Redirect Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_error_pages_redirect_rule.example <site_id>:<config_id>
```