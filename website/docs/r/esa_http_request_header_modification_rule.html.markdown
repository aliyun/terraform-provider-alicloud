---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_http_request_header_modification_rule"
description: |-
  Provides a Alicloud ESA Http Request Header Modification Rule resource.
---

# alicloud_esa_http_request_header_modification_rule

Provides a ESA Http Request Header Modification Rule resource.



For information about ESA Http Request Header Modification Rule and how to use it, see [What is Http Request Header Modification Rule](https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/api-esa-2024-09-10-createhttprequestheadermodificationrule).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_http_request_header_modification_rule&exampleId=f266c577-541b-05f6-08cc-fe780ddccfddc44c05d4&activeTab=example&spm=docs.r.esa_http_request_header_modification_rule.0.f266c57754&intl_lang=EN_US" target="_blank">
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

resource "alicloud_esa_rate_plan_instance" "instance" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "site" {
  site_name   = "gositecdn.cn"
  instance_id = alicloud_esa_rate_plan_instance.instance.id
  coverage    = "overseas"
  access_type = "NS"
}


resource "alicloud_esa_http_request_header_modification_rule" "default" {
  rule_name = "example_modify"
  request_header_modification {
    value     = "modify1"
    operation = "modify"
    name      = "example_modify1"
  }

  site_id      = alicloud_esa_site.site.id
  rule_enable  = "off"
  rule         = "(http.request.uri eq \"/content?page=1234\")"
  site_version = "0"
}
```

## Argument Reference

The following arguments are supported:
* `request_header_modification` - (Required, List) The configurations of modifying request headers. You can add, delete, or modify a request header. See [`request_header_modification`](#request_header_modification) below.
* `rule` - (Optional) The rule content.
* `rule_enable` - (Optional) Rule switch. Value range:
on: Open.
off: off.
* `rule_name` - (Optional) Rule Name.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Optional, ForceNew, Int) The version number of the website configurations.

### `request_header_modification`

The request_header_modification supports the following:
* `name` - (Required) Request Header Name.
* `operation` - (Required) Mode of operation. Value range:
add: add.
del: delete
modify: change.
* `value` - (Optional) Request header value

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Http Request Header Modification Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Http Request Header Modification Rule.
* `update` - (Defaults to 5 mins) Used when update the Http Request Header Modification Rule.

## Import

ESA Http Request Header Modification Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_http_request_header_modification_rule.example <site_id>:<config_id>
```