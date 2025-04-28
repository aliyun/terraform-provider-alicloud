---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_rewrite_url_rule"
description: |-
  Provides a Alicloud ESA Rewrite Url Rule resource.
---

# alicloud_esa_rewrite_url_rule

Provides a ESA Rewrite Url Rule resource.



For information about ESA Rewrite Url Rule and how to use it, see [What is Rewrite Url Rule](https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/api-esa-2024-09-10-createrewriteurlrule).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_rewrite_url_rule&exampleId=390154be-9918-c9aa-dadc-e1ec2907dcac5a555fc4&activeTab=example&spm=docs.r.esa_rewrite_url_rule.0.390154be99&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `rewrite_query_string_type` - (Optional) Query string rewrite type. Value range:
  - `static`: Static mode.
  - `dynamic`: Dynamic mode.
* `rewrite_uri_type` - (Optional) URI rewrite type. Value range:
  - `static`: Static mode.
  - `dynamic`: Dynamic mode.
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
● Match all incoming requests: value set to true
● Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Indicates whether the rule is enabled. Valid values:

  - on
  - off
* `rule_name` - (Optional) The rule name. You do not need to set this parameter when adding a global configuration.
* `site_id` - (Required, ForceNew, Int) The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.
      
* `site_version` - (Optional, ForceNew, Int) Version number of the site configuration. For a site with configuration version management enabled, you can use this parameter to specify the site version in which the configuration takes effect. The default version is 0.
* `uri` - (Optional) The desired URI to which you want to rewrite the path in the original request.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - ConfigId

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Rewrite Url Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Rewrite Url Rule.
* `update` - (Defaults to 5 mins) Used when update the Rewrite Url Rule.

## Import

ESA Rewrite Url Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_rewrite_url_rule.example <site_id>:<config_id>
```