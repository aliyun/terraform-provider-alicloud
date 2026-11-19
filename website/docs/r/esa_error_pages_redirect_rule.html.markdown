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

-> **NOTE:** Available since v2.4.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}

resource "alicloud_esa_error_pages_redirect_rule" "default" {
  site_id      = data.alicloud_esa_sites.default.sites.0.id
  rule_enable  = "off"
  rule_name    = "rule_example"
  rule         = "(http.host eq \"video.example.com\")"
  sequence     = "1"
  site_version = "0"
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

* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the [ListSites](https://help.aliyun.com/document_detail/2850189.html) operation.
* `rule_enable` - (Optional) The rule switch. You do not need to set this parameter when adding a global configuration. Valid values:
  - `on`: Enable.
  - `off`: Disable.
* `rule_name` - (Optional) Rule name. This parameter is not required when adding a global configuration.
* `rule` - (Optional) Rule content, which uses conditional expressions to match user requests. This parameter is not required when adding a global configuration. There are two usage scenarios:
  - Match all incoming requests: Set the value to `true`.
  - Match specific requests: Set the value to a custom expression, such as `(http.host eq "video.example.com")`.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites with configuration version management enabled, you can use this parameter to specify the site version for which the configuration takes effect. Default value: `0`.
* `sequence` - (Optional, ForceNew, Int) The execution sequence of the rule.
* `error_pages_redirect` - (Required, Set) Error code redirection. See [`error_pages_redirect`](#error_pages_redirect) below.

### `error_pages_redirect`

The error_pages_redirect supports the following:

* `status_code` - (Required) The response status code used when the node responds to the client with a redirection address. Valid values: `400`, `403`, `404`, `405`, `414`, `416`, `500`, `501`, `502`, `503`, `504`.
* `target_url` - (Required) The target URL after redirection.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Error Pages Redirect Rule. It formats as `<site_id>:<config_id>`.
* `config_id` - The configuration ID.
* `config_type` - The configuration type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Error Pages Redirect Rule.
* `update` - (Defaults to 5 mins) Used when update the Error Pages Redirect Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Error Pages Redirect Rule.

## Import

ESA Error Pages Redirect Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_error_pages_redirect_rule.example <site_id>:<config_id>
```
