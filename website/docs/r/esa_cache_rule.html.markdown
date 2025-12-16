---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_cache_rule"
description: |-
  Provides a Alicloud ESA Cache Rule resource.
---

# alicloud_esa_cache_rule

Provides a ESA Cache Rule resource.



For information about ESA Cache Rule and how to use it, see [What is Cache Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateCacheRule).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_cache_rule&exampleId=96d87eeb-0a12-a847-75b6-341ecbdf73cdb8bc3531&activeTab=example&spm=docs.r.esa_cache_rule.0.96d87eeb0a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_cache_rule" "default" {
  user_device_type            = "off"
  browser_cache_mode          = "no_cache"
  user_language               = "off"
  check_presence_header       = "headername"
  include_cookie              = "cookie_exapmle"
  edge_cache_mode             = "follow_origin"
  additional_cacheable_ports  = "2053"
  rule_name                   = "rule_example"
  edge_status_code_cache_ttl  = "300"
  browser_cache_ttl           = "300"
  query_string                = "example"
  user_geo                    = "off"
  sort_query_string_for_cache = "off"
  check_presence_cookie       = "cookiename"
  cache_reserve_eligibility   = "bypass_cache_reserve"
  query_string_mode           = "ignore_all"
  rule                        = "http.host eq \"video.example.com\""
  cache_deception_armor       = "off"
  site_id                     = data.alicloud_esa_sites.default.sites.0.id
  bypass_cache                = "cache_all"
  edge_cache_ttl              = "300"
  rule_enable                 = "off"
  site_version                = "0"
  include_header              = "example"
  serve_stale                 = "off"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_cache_rule&spm=docs.r.esa_cache_rule.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `additional_cacheable_ports` - (Optional) Enable caching on specified ports. Value range: 8880, 2052, 2082, 2086, 2095, 2053, 2083, 2087, 2096
* `browser_cache_mode` - (Optional) Browser cache mode. Possible values:
  - `no_cache`no_cache: Do not cache.
  - `follow_origin`: Follow the origin server's cache policy.
  - `override_origin`: Override the origin server's cache policy.
* `browser_cache_ttl` - (Optional) Browser cache expiration time in seconds.
* `bypass_cache` - (Optional) Set the bypass cache mode. Possible values:
  - `cache_all`: Cache all requests.
  - `bypass_all`: Bypass cache for all requests.
* `cache_deception_armor` - (Optional) Cache deception protection. Used to defend against web cache deception attacks, only the cache content that passes the validation will be cached. Value range:
  - `on`: Enabled.
  - `off`: Disabled.
* `cache_reserve_eligibility` - (Optional) Cache retention eligibility. Used to control whether user requests bypass the cache retention node when returning to the origin. Possible values:
  - `bypass_cache_reserve`: Requests bypass cache retention.
  - `eligible_for_cache_reserve`: Eligible for cache retention.
* `check_presence_cookie` - (Optional) When generating the cache key, check if the cookie exists. If it does, add the cookie name (case-insensitive) to the cache key. Multiple cookie names are supported, separated by spaces.
* `check_presence_header` - (Optional) When generating the cache key, check if the header exists. If it does, add the header name (case-insensitive) to the cache key. Multiple header names are supported, separated by spaces.
* `edge_cache_mode` - (Optional) Edge cache mode. Possible values:
  - `follow_origin`: Follow the origin server's cache policy (if it exists), otherwise use the default cache policy.
  - `no_cache`: Do not cache.
  - `override_origin`: Override the origin server's cache policy.
  - `follow_origin_bypass`: Follow the origin server's cache policy (if it exists), otherwise do not cache.
  - `follow_origin_override `: Follow the origin server's cache policy (if it exists), otherwise use custom cache TTL.
* `edge_cache_ttl` - (Optional) Edge cache expiration time in seconds.
* `edge_status_code_cache_ttl` - (Optional) Status code cache expiration time in seconds.
* `include_cookie` - (Optional) When generating the cache key, add the specified cookie names and their values. Multiple values are supported, separated by spaces.
* `include_header` - (Optional) When generating the cache key, add the specified header names and their values. Multiple values are supported, separated by spaces.
* `query_string` - (Optional) Query strings to be reserved or excluded. Multiple values are supported, separated by spaces.
* `query_string_mode` - (Optional) The processing mode for query strings when generating the cache key. Possible values:
  - `ignore_all`: Ignore all.
  - `exclude_query_string`: Exclude specified query strings.
  - `reserve_all`: Default, reserve all.
  - `include_query_string`: Include specified query strings.
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true.
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\").
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - `on`: open.
  - `off`: close.
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Computed, Int, Available since v1.262.0) The rule execution order prioritizes lower numerical values. It is only applicable when setting or modifying the order of individual rule configurations.
* `serve_stale` - (Optional) Serve stale cache. When enabled, the node can still respond to user requests with expired cached files when the origin server is unavailable. Value range:
  - `on`: Enabled.
  - `off`: Disabled.
* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the [ListSites] API.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.
* `sort_query_string_for_cache` - (Optional) Query string sorting, disabled by default. Possible values:
  - `on`: Enable.
  - `off`: Disable.
* `user_device_type` - (Optional) When generating the cache key, add the client device type. Possible values:
  - `on`: Enable.
  - `off`: Disable.
* `user_geo` - (Optional) When generating the cache key, add the client's geographic location. Possible values:
  - `on`: Enable.
  - `off`: Disable.
* `user_language` - (Optional) When generating cache keys, include the client's language type. Possible values:
  - `on`: Enable.
  - `off`: Disable.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<cache_rule_id>`.
* `cache_rule_id` - Cache Rule Id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cache Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Cache Rule.
* `update` - (Defaults to 5 mins) Used when update the Cache Rule.

## Import

ESA Cache Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_cache_rule.example <site_id>:<cache_rule_id>
```