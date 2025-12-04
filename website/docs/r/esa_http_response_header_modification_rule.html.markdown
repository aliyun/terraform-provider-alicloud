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

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_esa_site" "resource_Site_HttpResponseHeaderModificationRule_example" {
  site_name   = "gositecdn-${random_integer.default.result}.cn"
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
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - `on`: open.
  - `off`: close.
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Computed, Int, Available since v1.262.1) The rule execution order prioritizes lower numerical values. It is only applicable when setting or modifying the order of individual rule configurations.
* `site_id` - (Required, ForceNew) The site ID.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.

### `response_header_modification`

The response_header_modification supports the following:
* `name` - (Required) The response header name.
* `operation` - (Required) Operation method. Possible values:
  - `add`: Add
  - `del`: Delete
  - `modify`: Modify
* `type` - (Optional, Computed, Available since v1.262.1) The value type. Value range:
  - `static`: Static mode.
  - `dynamic`: Dynamic mode.
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