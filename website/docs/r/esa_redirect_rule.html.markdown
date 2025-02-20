---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_redirect_rule"
description: |-
  Provides a Alicloud ESA Redirect Rule resource.
---

# alicloud_esa_redirect_rule

Provides a ESA Redirect Rule resource.



For information about ESA Redirect Rule and how to use it, see [What is Redirect Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateRedirectRule).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_redirect_rule&exampleId=c9456582-adcc-430d-667d-39ff87ec6940ece1773d&activeTab=example&spm=docs.r.esa_redirect_rule.0.c9456582ad&intl_lang=EN_US" target="_blank">
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

resource "alicloud_esa_rate_plan_instance" "resource_RedirectRule_example" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_RedirectRule_example" {
  site_name   = "gositecdn.cn"
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
* `rule` - (Optional) The rule content.
* `rule_enable` - (Optional) Indicates whether the rule is enabled. Valid values:

  - `on`
  - `off`
* `rule_name` - (Optional) 规则名，可以查出规则名为所传字段的那条规则，只有传了functionName才生效
* `site_id` - (Required, ForceNew, Int) The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.
* `site_version` - (Optional, ForceNew, Int) The version of the website configurations.
* `status_code` - (Required) The response code that you want to use to indicate URL redirection. Valid values:

  - 301
  - 302
  - 303
  - 307
  - 308
* `target_url` - (Required) The destination URL to which requests are redirected.
* `type` - (Required) The redirect type. Valid value:

  - static

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Redirect Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Redirect Rule.
* `update` - (Defaults to 5 mins) Used when update the Redirect Rule.

## Import

ESA Redirect Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_redirect_rule.example <site_id>:<config_id>
```