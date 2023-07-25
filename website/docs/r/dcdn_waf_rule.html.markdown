---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_waf_rule"
sidebar_current: "docs-alicloud-resource-dcdn-waf-rule"
description: |-
  Provides a Alicloud Dcdn Waf Rule resource.
---

# alicloud_dcdn_waf_rule

Provides a Dcdn Waf Rule resource.

For information about Dcdn Waf Rule and how to use it, see [What is Waf Rule](https://www.alibabacloud.com/help/en/dynamic-route-for-cdn/latest/configure-protection-rules).

-> **NOTE:** Available since v1.201.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_dcdn_waf_policy" "example" {
  defense_scene = "waf_group"
  policy_name   = var.name
  policy_type   = "custom"
  status        = "on"
}

resource "alicloud_dcdn_waf_rule" "example" {
  policy_id = alicloud_dcdn_waf_policy.example.id
  rule_name = var.name
  conditions {
    key      = "URI"
    op_value = "ne"
    values   = "/login.php"
  }
  conditions {
    key      = "Header"
    sub_key  = "a"
    op_value = "eq"
    values   = "b"
  }
  status    = "on"
  cc_status = "on"
  action    = "monitor"
  effect    = "rule"
  rate_limit {
    target    = "IP"
    interval  = "5"
    threshold = "5"
    ttl       = "1800"
    status {
      code  = "200"
      ratio = "60"
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `action` - (Optional) Specifies the action of the rule. Valid values: `block`, `monitor`, `js`.
* `cc_status` - (Optional) Specifies whether to enable rate limiting. Valid values: `on` and `off`. **NOTE:** This parameter is required when policy is of type `custom_acl`.
* `cn_region_list` - (Optional) The blocked regions in the Chinese mainland, separated by commas (,).
* `conditions` - (Optional) Conditions that trigger the rule. See [`conditions`](#conditions) below. **NOTE:** This parameter is required when policy is of type `custom_acl` or `whitelist`.
* `effect` - (Optional) The effective scope of the rate limiting blacklist. If you set ccStatus to on, you must configure this parameter. Valid values: `rule` (takes effect for the current rule) and `service` (takes effect globally).
* `other_region_list` - (Optional) Blocked regions outside the Chinese mainland, separated by commas (,).
* `policy_id` - (Required, ForceNew) The protection policy ID.
* `rate_limit` - (Optional) The rules of rate limiting. If you set `cc_status` to on, you must configure this parameter. See [`rate_limit`](#rate_limit) below.
* `regular_rules` - (Optional) The regular expression.e, when waf_group appears in tags, this value can be filled in, and only one list of six digits in string format can appear with regultypes.
* `regular_types` - (Optional) Regular rule type, when waf_group appears in tags, this value can be filled in, optional values:["sqli", "xss", "code_exec", "crlf", "lfileii", "rfileii", "webshell", "vvip", "other"]
* `remote_addr` - (Optional) Filter by IP address.
* `rule_name` - (Required) The name of the protection rule. The name can be up to 64 characters in length and can contain letters, digits, and underscores (_). **NOTE:** This parameter cannot be modified when policy is of type `region_block`.
* `scenes` - (Optional) The types of the protection policies.
* `status` - (Optional) The status of the waf rule. Valid values: `on` and `off`. Default value: on.
* `waf_group_ids` - (Optional) The id of the waf rule group. The default value is "1012". Multiple rules are separated by commas.

### `conditions`

The conditions supports the following:
* `key` - (Required) The match field.
* `op_value` - (Required) The logical symbol.
* `sub_key` - (Optional) The match subfield.
* `values` - (Optional) The match content. Separate multiple values with commas (,).

### `rate_limit`

The rate_limit supports the following:
* `interval` - (Optional) Statistical duration, 5-1800.
* `status` - (Optional) Response code statistics. See [`status`](#rate_limit-status) below.
* `sub_key` - (Optional) The subfield of the target field. If you set `target` to `Header`, `Query String Parameter`, or `Cookie Name`, you must configure `sub_key`.
* `target` - (Optional) The statistical field for frequency control. Currently, `IP`, `Header`, `Query String Parameter`, `Cookie Name`, `Session` is supported.
* `threshold` - (Optional) The trigger threshold of rate limiting. Valid values: 2 to 500000. Unit: requests.
* `ttl` - (Optional) The validity period of the blacklist. Valid values: 60 to 86400. Unit: seconds.

### `rate_limit-status`

The status supports the following:
* `code` - (Optional) The HTTP status code returned.
* `count` - (Optional) The number of times that the HTTP status code that was returned. Valid values: 2 to 50000. You can configure only one of the `ratio` and `count` fields.
* `ratio` - (Optional) The percentage of HTTP status codes. Valid values: 1 to 100. You can configure only one of the `ratio` and `count` fields.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `defense_scene` - The type of protection policy. The following scenarios are supported:-waf_group:Web basic protection-custom_acl: Custom protection policy-whitelist: whitelist
* `gmt_modified` - Revised the time. The date format is based on ISO8601 notation and uses UTC +0 time in the format of yyyy-MM-ddTHH:mm:ssZ.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Waf Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Waf Rule.
* `update` - (Defaults to 5 mins) Used when update the Waf Rule.

## Import

Dcdn Waf Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_waf_rule.example <id>
```