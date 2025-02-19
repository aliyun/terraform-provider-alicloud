---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_http_response_header_modification_rule"
description: |-
  Provides a Alicloud ESA Http Response Header Modification Rule resource.
---

# alicloud_esa_http_response_header_modification_rule

Provides a ESA Http Response Header Modification Rule resource.



For information about ESA Http Response Header Modification Rule and how to use it, see [What is Http Response Header Modification Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateHttpResponseHeaderModificationRule).

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
* `response_header_modification` - (Required, List) The configurations of modifying response headers. You can add, delete, or modify a response header. See [`response_header_modification`](#response_header_modification) below.
* `rule` - (Optional) The rule content.
* `rule_enable` - (Optional) Indicates whether the rule is enabled. Valid values:

  - on
  - off
* `rule_name` - (Optional) The rule name.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Optional, ForceNew, Int) The version number of the website configurations.

### `response_header_modification`

The response_header_modification supports the following:
* `name` - (Required) The response header name.
* `operation` - (Required) Mode of operation.
* `value` - (Optional) The response header value.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Http Response Header Modification Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Http Response Header Modification Rule.
* `update` - (Defaults to 5 mins) Used when update the Http Response Header Modification Rule.

## Import

ESA Http Response Header Modification Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_http_response_header_modification_rule.example <site_id>:<config_id>
```