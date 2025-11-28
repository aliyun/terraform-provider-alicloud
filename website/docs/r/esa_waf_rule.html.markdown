---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_waf_rule"
description: |-
  Provides a Alicloud ESA Waf Rule resource.
---

# alicloud_esa_waf_rule

Provides a ESA Waf Rule resource.

The detailed configuration of a Web Application Firewall (WAF) rule.

For information about ESA Waf Rule and how to use it, see [What is Waf Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/BatchCreateWafRules).

-> **NOTE:** Available since v1.261.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_waf_rule&exampleId=032985c1-00b2-e5cf-d900-02f61f3b1d5843fd5677&activeTab=example&spm=docs.r.esa_waf_rule.0.032985c100&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_waf_ruleset" "default" {
  site_id      = data.alicloud_esa_sites.default.sites.0.site_id
  phase        = "http_custom"
  site_version = "0"
}

resource "alicloud_esa_waf_rule" "default" {
  ruleset_id = alicloud_esa_waf_ruleset.default.ruleset_id
  phase      = "http_custom"
  config {
    status     = "on"
    action     = "deny"
    expression = "(http.host in {\"123.example.top\"})"
    actions {
      response {
        id   = "0"
        code = "403"
      }

    }

    name = "111"
  }

  site_version = "0"
  site_id      = data.alicloud_esa_sites.default.sites.0.site_id
}
```

## Argument Reference

The following arguments are supported:
* `config` - (Optional, List) The specific configuration of the WAF rule. See [`config`](#config) below.
* `phase` - (Required, ForceNew) The phase in which the WAF processes this rule.
* `ruleset_id` - (Optional, ForceNew, Int) The ID of the WAF ruleset, which can be obtained by calling the [ListWafRulesets](https://www.alibabacloud.com/help/en/doc-detail/2850233.html) operation.
* `shared` - (Optional, List) Shared configuration attributes used across multiple rules. See [`shared`](#shared) below.

-> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `site_id` - (Required, ForceNew) The unique identifier of the website, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.
* `site_version` - (Optional, Int) The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.

-> **NOTE:** This parameter only applies during resource creation, update or deletion. If modified in isolation without other property changes, Terraform will not trigger any action.


### `config`

The config supports the following:
* `action` - (Optional) The action to perform when a request matches this rule.
* `actions` - (Optional, List) Extended action configurations, including custom responses and bypass settings. See [`actions`](#config-actions) below.
* `app_package` - (Optional, List) Security mechanism to prevent apps from being repackaged. See [`app_package`](#config-app_package) below.
* `app_sdk` - (Optional, List) Mobile app SDK-related configurations. See [`app_sdk`](#config-app_sdk) below.
* `expression` - (Optional) The match expression used to evaluate incoming requests.
* `managed_list` - (Optional) The name of the managed list applied to this rule.
* `managed_rulesets` - (Optional, Set) The managed rulesets referenced by this rule and their configurations. See [`managed_rulesets`](#config-managed_rulesets) below.
* `name` - (Optional) The display name of the WAF rule.
* `notes` - (Optional) Additional notes about this rule.
* `rate_limit` - (Optional, List) Configuration of the rate limiting rule. See [`rate_limit`](#config-rate_limit) below.
* `security_level` - (Optional, List) The overall security protection level of WAF.
Valid values:
  - off
  - essentially_off
  - low
  - medium
  - high
  - under_attack See [`security_level`](#config-security_level) below.
* `sigchl` - (Optional, Set) Configuration items for token verification mechanisms.
* `status` - (Optional) The status of the WAF rule: whether it is enabled or disabled.
* `timer` - (Optional, List) Configuration for the time schedule when the rule takes effect. See [`timer`](#config-timer) below.
* `type` - (Optional, ForceNew) The type category of the WAF rule.
* `value` - (Optional) The IP address allowed or denied in IP access control.

### `config-actions`

The config-actions supports the following:
* `bypass` - (Optional, List) The skip configuration specified by the whitelist rule. See [`bypass`](#config-actions-bypass) below.
* `response` - (Optional, List) The custom error page returned when the rule is triggered. See [`response`](#config-actions-response) below.

### `config-app_package`

The config-app_package supports the following:
* `package_signs` - (Optional, Set) Security mechanism to prevent apps from being repackaged. See [`package_signs`](#config-app_package-package_signs) below.

### `config-app_sdk`

The config-app_sdk supports the following:
* `custom_sign` - (Optional, List) Custom fields used for mobile app signature validation. See [`custom_sign`](#config-app_sdk-custom_sign) below.
* `custom_sign_status` - (Optional) Indicates whether the custom signature field validation is enabled.
* `feature_abnormal` - (Optional, Set) Detected abnormal behaviors of the application.

### `config-managed_rulesets`

The config-managed_rulesets supports the following:
* `action` - (Optional) The default action applied to all rules in this ruleset.
* `attack_type` - (Optional, Int) The primary attack type targeted by this ruleset.
* `managed_rules` - (Optional, Set) The individual managed rules included in this ruleset. See [`managed_rules`](#config-managed_rulesets-managed_rules) below.
* `protection_level` - (Optional, Int) The protection strength level assigned to this ruleset.

### `config-rate_limit`

The config-rate_limit supports the following:
* `characteristics` - (Optional, List) The statistical dimensions to which the rate limiting rule applies. See [`characteristics`](#config-rate_limit-characteristics) below.
* `interval` - (Optional, Int) The statistical interval.
* `on_hit` - (Optional) Indicates whether the rule applies to requests that hit the cache.
* `threshold` - (Optional, List) Threshold settings for the rate limiting rule. See [`threshold`](#config-rate_limit-threshold) below.
* `ttl` - (Optional, Int) The timeout period for creating the stack used in rate limiting.

### `config-security_level`

The config-security_level supports the following:
* `value` - (Optional) The security protection level of WAF.
Valid values:
  - off
  - essentially_off
  - low
  - medium
  - high
  - under_attack

### `config-timer`

The config-timer supports the following:
* `periods` - (Optional, Set) One-time effective time intervals. See [`periods`](#config-timer-periods) below.
* `scopes` - (Optional) Timing type:
  - `permanent`: Always active
  - `periods`: Active in specified periods
  - `weekly`: Recurring weekly schedule
* `weekly_periods` - (Optional, Set) Weekly recurring time schedules. See [`weekly_periods`](#config-timer-weekly_periods) below.
* `zone` - (Optional, Int) The time zone. If it is not specified, the default value is UTC +00:00.  Example: 8 means East Zone 8,-8 means West Zone 8  Range:-12 -+14

### `config-timer-periods`

The config-timer-periods supports the following:
* `end` - (Optional) End time, value is UTC time in RFC3339 format.
* `start` - (Optional) Start time, value is UTC time in RFC3339 format.

### `config-timer-weekly_periods`

The config-timer-weekly_periods supports the following:
* `daily_periods` - (Optional, Set) Daily effective time periods within a weekly schedule. See [`daily_periods`](#config-timer-weekly_periods-daily_periods) below.
* `days` - (Optional) Cycle, multiple use comma separated, 1-7 respectively represent Monday-Sunday.  Example: Monday, Wednesday value is "1,3"

### `config-timer-weekly_periods-daily_periods`

The config-timer-weekly_periods-daily_periods supports the following:
* `end` - (Optional) End time in HH:mm:ss format
* `start` - (Optional) Start time in HH:mm:ss format

### `config-rate_limit-characteristics`

The config-rate_limit-characteristics supports the following:
* `criteria` - (Optional, Set) The details of logical databases. See [`criteria`](#config-rate_limit-characteristics-criteria) below.
* `logic` - (Optional) Logical relationship, the current value of ogic is only and, format reference above.

### `config-rate_limit-threshold`

The config-rate_limit-threshold supports the following:
* `distinct_managed_rules` - (Optional, Int) The maximum number of distinct managed rules that can be triggered.
* `managed_rules_blocked` - (Optional, Int) The maximum number of times that managed rules can be triggered.
* `request` - (Optional, Int) The maximum number of allowed requests within a time interval.
* `response_status` - (Optional, List) Limits on the frequency of returning specific HTTP status codes. See [`response_status`](#config-rate_limit-threshold-response_status) below.
* `traffic` - (Optional) The maximum allowed traffic within a time interval (deprecated).

### `config-rate_limit-threshold-response_status`

The config-rate_limit-threshold-response_status supports the following:
* `code` - (Optional, Int) HTTP response code. Currently, the available value of code is only 404. Valid values: '404'
* `count` - (Optional, Int) The maximum number of times the specified status code can be returned.
* `ratio` - (Optional, Int) The upper limit of the percentage of occurrences of the specified status code among all responses.

### `config-rate_limit-characteristics-criteria`

The config-rate_limit-characteristics-criteria supports the following:
* `criteria` - (Optional, Set) Nested collection of matching criteria. See [`criteria`](#config-rate_limit-characteristics-criteria-criteria) below.
* `logic` - (Optional) Logical relationship between multiple matching conditions.
* `match_type` - (Optional) The request field targeted by the match condition.

### `config-rate_limit-characteristics-criteria-criteria`

The config-rate_limit-characteristics-criteria-criteria supports the following:
* `criteria` - (Optional, Set) Nested collection of matching criteria. See [`criteria`](#config-rate_limit-characteristics-criteria-criteria-criteria) below.
* `logic` - (Optional) Logical relationship between multiple matching conditions.
* `match_type` - (Optional) The request field targeted by the match condition.

### `config-rate_limit-characteristics-criteria-criteria-criteria`

The config-rate_limit-characteristics-criteria-criteria-criteria supports the following:
* `match_type` - (Optional) The request field targeted by the match condition.

### `config-managed_rulesets-managed_rules`

The config-managed_rulesets-managed_rules supports the following:
* `action` - (Optional) The action performed on requests that match the managed rule.
* `id` - (Optional, Int) The unique identifier of a managed rule.
* `status` - (Optional) The status of the managed rule: whether it is enabled or disabled.

### `config-app_sdk-custom_sign`

The config-app_sdk-custom_sign supports the following:
* `key` - (Optional) The name of the custom signature field used for validation.
* `value` - (Optional) The value of the custom signature field used for validation.

### `config-app_package-package_signs`

The config-app_package-package_signs supports the following:
* `name` - (Optional) The package name of an authorized application.
* `sign` - (Optional) The digital signature of a legitimate app package.

### `config-actions-bypass`

The config-actions-bypass supports the following:
* `custom_rules` - (Optional, Set) The IDs of custom rules to skip.
* `regular_rules` - (Optional, Set) The IDs of specific managed rules to skip.
* `regular_types` - (Optional, Set) The types of managed rules to skip.
* `skip` - (Optional) The scope that is skipped when requests match conditions defined in the whitelist rule.
* `tags` - (Optional, Set) The rule categories that are skipped when requests match conditions defined in the whitelist rule.

### `config-actions-response`

The config-actions-response supports the following:
* `code` - (Optional, Int) The HTTP response code returned to the client.
* `id` - (Optional, Int) The ID of the custom error page, which can be obtained by calling the ListPages operation.

### `shared`

The shared supports the following:
* `action` - (Optional) The default action executed under shared configuration.
* `actions` - (Optional, List) Extended action configurations under shared settings. See [`actions`](#shared-actions) below.
* `cross_site_id` - (Optional, Int) Specify the cross-domain site ID.
* `expression` - (Optional) The match expression used in shared configuration.
* `match` - (Optional, List) Configuration of the request matching logic engine. See [`match`](#shared-match) below.
* `mode` - (Optional) The integration mode of the Web SDK:
  - `automatic`: Automatically integrated
  - `manual`: Manually integrated
* `name` - (Optional) The display name of the ruleset.
* `target` - (Optional) The target type protected by this rule: web or app.

### `shared-actions`

The shared-actions supports the following:
* `response` - (Optional, List) Custom response configuration under shared settings. See [`response`](#shared-actions-response) below.

### `shared-match`

The shared-match supports the following:
* `criteria` - (Optional, List) A collection of matching criteria. See [`criteria`](#shared-match-criteria) below.
* `logic` - (Optional) Logical relationship between multiple matching conditions.
* `match_type` - (Optional) The request field targeted by the match condition.

### `shared-match-criteria`

The shared-match-criteria supports the following:
* `criteria` - (Optional, List) A collection of matching criteria. See [`criteria`](#shared-match-criteria-criteria) below.
* `logic` - (Optional) Logical relationship between multiple matching conditions.
* `match_type` - (Optional) The request field targeted by the match condition.

### `shared-match-criteria-criteria`

The shared-match-criteria-criteria supports the following:
* `criteria` - (Optional, List) A collection of matching criteria. See [`criteria`](#shared-match-criteria-criteria-criteria) below.
* `logic` - (Optional) Logical relationship between multiple matching conditions.
* `match_type` - (Optional) The request field targeted by the match condition.

### `shared-match-criteria-criteria-criteria`

The shared-match-criteria-criteria-criteria supports the following:
* `match_type` - (Optional) The request field targeted by the match condition.

### `shared-actions-response`

The shared-actions-response supports the following:
* `code` - (Optional, Int) The HTTP response code returned under shared configuration.
* `id` - (Optional, Int) The ID of the custom error page used in shared configuration.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<waf_rule_id>`.
* `config` - The specific configuration of the WAF rule.
  * `id` - The internal unique ID of the WAF rule.
  * `managed_group_id` - The ID of the managed rule group (deprecated).
  * `managed_rulesets` - The managed rulesets referenced by this rule and their configurations.
    * `number_enabled` - Number of rules currently enabled.
    * `number_total` - Total number of rules in this ruleset.
* `waf_rule_id` - The unique identifier of the WAF rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Waf Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Waf Rule.
* `update` - (Defaults to 5 mins) Used when update the Waf Rule.

## Import

ESA Waf Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_waf_rule.example <site_id>:<waf_rule_id>
```