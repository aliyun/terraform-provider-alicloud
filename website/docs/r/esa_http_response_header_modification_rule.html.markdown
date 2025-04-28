---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_http_response_header_modification_rule"
description: |-
  Provides a Alicloud ESA Http Response Header Modification Rule resource.
---

# alicloud_esa_http_response_header_modification_rule

Provides a ESA Http Response Header Modification Rule resource.



For information about ESA Http Response Header Modification Rule and how to use it, see [What is Http Response Header Modification Rule](https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/api-esa-2024-09-10-createhttpresponseheadermodificationrule).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_http_response_header_modification_rule&exampleId=77b37dbc-9f42-bc83-09b3-62c7969cebd6e13219af&activeTab=example&spm=docs.r.esa_http_response_header_modification_rule.0.77b37dbc9f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_esa_rate_plan_instance" "resource_HttpResponseHeaderModificationRule_example" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_HttpResponseHeaderModificationRule_example" {
  site_name   = "gositecdn.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_HttpResponseHeaderModificationRule_example.id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_http_response_header_modification_rule" "default" {
  rule_enable = "on"
  response_header_modification {
    value     = "add"
    operation = "add"
    name      = "exampleadd"
  }
  response_header_modification {
    operation = "del"
    name      = "exampledel"
  }
  response_header_modification {
    operation = "modify"
    name      = "examplemodify"
    value     = "modify"
  }

  rule         = "(http.host eq \"video.example.com\")"
  site_version = "0"
  rule_name    = "exampleResponseHeader"
  site_id      = alicloud_esa_site.resource_Site_HttpResponseHeaderModificationRule_example.id
}
```

## Argument Reference

The following arguments are supported:
* `response_header_modification` - (Required, List) Modify response headers, supporting add, delete, and modify operations. See [`response_header_modification`](#response_header_modification) below.
* `rule` - (Optional) Rule content.
* `rule_enable` - (Optional) Rule switch. Possible values:
  - `on`: Enable.
  - `off`: Disable.
* `rule_name` - (Optional) Rule name.
* `site_id` - (Required, ForceNew, Int) The site ID.
* `site_version` - (Optional, ForceNew, Int) The version number of the website configurations.

### `response_header_modification`

The response_header_modification supports the following:
* `name` - (Required) The response header name.
* `operation` - (Required) Operation method. Possible values:
  - `add`: Add
  - `del`: Delete
  - `modify`: Modify
* `value` - (Optional) The response header value.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Http Response Header Modification Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Http Response Header Modification Rule.
* `update` - (Defaults to 5 mins) Used when update the Http Response Header Modification Rule.

## Import

ESA Http Response Header Modification Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_http_response_header_modification_rule.example <site_id>:<config_id>
```