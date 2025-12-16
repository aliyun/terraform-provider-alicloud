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
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_http_request_header_modification_rule&exampleId=4e7df8b4-3c60-887a-f0a2-72e0c816d15c1bf8882e&activeTab=example&spm=docs.r.esa_http_request_header_modification_rule.0.4e7df8b43c&intl_lang=EN_US" target="_blank">
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

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_esa_site" "site" {
  site_name   = "gositecdn-${random_integer.default.result}.cn"
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_http_request_header_modification_rule&spm=docs.r.esa_http_request_header_modification_rule.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `request_header_modification` - (Required, List) The configurations of modifying request headers. You can add, delete, or modify a request header. See [`request_header_modification`](#request_header_modification) below.
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - on: open.
  - off: close.
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Int, Available since v1.263.0) The rule execution order prioritizes lower numerical values. It is only applicable when setting or modifying the order of individual rule configurations.
* `site_id` - (Required, ForceNew) The site ID.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.

### `request_header_modification`

The request_header_modification supports the following:
* `name` - (Required) Request Header Name.
* `operation` - (Required) Mode of operation. Value range:
  - `add`: add.
  - `del`: delete
  - `modify`: change.
* `type` - (Optional, Available since v1.263.0) Value type. Value range:
  - `static`:static mode.
  - `dynamic`:dynamic mode.
* `value` - (Optional) Request header value

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Http Request Header Modification Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Http Request Header Modification Rule.
* `update` - (Defaults to 5 mins) Used when update the Http Request Header Modification Rule.

## Import

ESA Http Request Header Modification Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_http_request_header_modification_rule.example <site_id>:<config_id>
```