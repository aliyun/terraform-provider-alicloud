---
subcategory: "Anti-DDoS Pro (DdosCoo)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_web_cc_rule"
description: |-
  Provides a Alicloud DdosCoo Web Cc Rule resource.
---

# alicloud_ddoscoo_web_cc_rule

Provides a DdosCoo Web Cc Rule resource.

CC frequency control rules.

For information about DdosCoo Web Cc Rule and how to use it, see [What is Web Cc Rule](https://next.api.alibabacloud.com/document/ddoscoo/2020-01-01/ConfigWebCCRuleV2).

-> **NOTE:** Available since v1.271.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform"
}

variable "domain" {
  default = "terraform-example.alibaba.com"
}

data "alicloud_ddoscoo_instances" "default" {
}

resource "alicloud_ddoscoo_domain_resource" "default" {
  domain       = var.domain
  instance_ids = [data.alicloud_ddoscoo_instances.default.ids.0]
  proxy_types {
    proxy_ports = [443]
    proxy_type  = "https"
  }
  real_servers = ["177.167.32.11"]
  rs_type      = 0
}

resource "alicloud_ddoscoo_web_cc_rule" "default" {
  rule_detail {
    action = "block"
    rate_limit {
      interval  = "11"
      threshold = "2"
      ttl       = "840"
      target    = "header"
      sub_key   = "33"
    }
    condition {
      match_method = "belong"
      field        = "ip"
      content      = "1.1.1.1"
    }
    condition {
      match_method = "contain"
      field        = "uri"
      content      = "/a"
    }
    condition {
      match_method = "contain"
      field        = "header"
      header_name  = "123"
      content      = "1234"
    }
    statistics {
      mode        = "distinct"
      field       = "header"
      header_name = "12"
    }
    status_code {
      enabled         = true
      code            = "100"
      use_ratio       = false
      count_threshold = "2"
      ratio_threshold = "5"
    }
  }
  name   = var.name
  domain = alicloud_ddoscoo_domain_resource.default.id
}
```

## Argument Reference

The following arguments are supported:
* `domain` - (Required, ForceNew) The domain name of the website service.  

-> **NOTE:**  The domain name must already have website service forwarding rules configured. You can call [DescribeDomains](https://help.aliyun.com/document_detail/91724.html) to query all domain names.  

* `name` - (Required, ForceNew) Rule name.
* `rule_detail` - (Required, Set) Rule details.   See [`rule_detail`](#rule_detail) below.

### `rule_detail`

The rule_detail supports the following:
* `action` - (Required) The action to take when a match occurs. Valid values:
  - `accept`: Allow
  - `block`: Block
  - `challenge`: Challenge
  - `watch`: Monitor
* `condition` - (Required, List) List of matching conditions.   See [`condition`](#rule_detail-condition) below.
* `rate_limit` - (Required, Set) Rate limiting statistics. See [`rate_limit`](#rule_detail-rate_limit) below.
* `statistics` - (Optional, Set) Deduplicated statistics. This parameter is optional. If omitted, deduplication is not applied. See [`statistics`](#rule_detail-statistics) below.
* `status_code` - (Optional, Set) The HTTP status code. See [`status_code`](#rule_detail-status_code) below.

### `rule_detail-condition`

The rule_detail-condition supports the following:
* `content` - (Required) Matching content.
* `field` - (Required) Matching field.  
* `header_name` - (Optional) Custom HTTP header field name.

-> **NOTE:**  Valid only when `Field` is set to `header`.

* `match_method` - (Required) Matching method.  

### `rule_detail-rate_limit`

The rule_detail-rate_limit supports the following:
* `interval` - (Required, Int) Statistical interval. Unit: seconds.
* `sub_key` - (Optional) Header field name (required only when the statistic source is `header`).
* `target` - (Required) Statistic source. Valid values:
  - `ip`: Statistics are collected by IP address.
  - `header`: Statistics are collected by header.
* `threshold` - (Required, Int) The trigger threshold.
* `ttl` - (Required, Int) Block duration. Unit: seconds.

### `rule_detail-statistics`

The rule_detail-statistics supports the following:
* `field` - (Required) The statistic source. Valid values:
  - `ip`: Count by IP address
  - `header`: Count by HTTP header
  - `uri`: Count by URI
* `header_name` - (Optional) Set this parameter only when the statistic source is `header`.
* `mode` - (Required) The deduplication mode. Valid values:
  - `count`: No deduplication
  - `distinct`: Deduplicated count

### `rule_detail-status_code`

The rule_detail-status_code supports the following:
* `code` - (Required, Int) Status code. The value range is `100` to `599`:
* `count_threshold` - (Optional, Int) When the ratio is not used, the enforcement action is triggered only when the corresponding status code reaches `CountThreshold`. The value range is `2` to `50000`.
* `enabled` - (Required) Whether the rule is enabled. Valid values:
  - `true`: Enabled.
  - `false`: Disabled.
* `ratio_threshold` - (Optional, Int) When the ratio is used, the enforcement action is triggered only when the corresponding status code reaches `RatioThreshold`. The value range is `1` to `100`.
* `use_ratio` - (Required) Whether to use a ratio:
  - `true`: Use ratio.
  - `false`: Do not use ratio.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<domain>:<name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Web Cc Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Web Cc Rule.
* `update` - (Defaults to 5 mins) Used when update the Web Cc Rule.

## Import

DdosCoo Web Cc Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddoscoo_web_cc_rule.example <domain>:<name>
```
