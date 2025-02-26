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

## Argument Reference

The following arguments are supported:
* `additional_cacheable_ports` - (Optional) Enable caching on the specified port. value: 8880, 2052, 2082, 2086, 2095, 2053, 2083, 2087, 2096.
* `browser_cache_mode` - (Optional) Browser cache mode. value:
           - no_cache: Not cached.
           - follow_origin: follows the Origin cache policy.
           - override_origin: replaces the Origin cache policy.
* `browser_cache_ttl` - (Optional) The browser cache expiration time, in seconds.
* `bypass_cache` - (Optional) Set the cache bypass mode. value:
           - cache_all: All requests are cached.
           - bypass_all: All requests bypass the cache.
* `cache_deception_armor` - (Optional) Cache spoofing defense. Used to defend against Web cache spoofing attacks, the cached content that passes the check is cached. value:
           - on: open.
           - off: off.
* `cache_reserve_eligibility` - (Optional) Cache retention eligibility. Used to control whether the user request bypasses the cache retention node when returning to the source. value:
           - bypass_cache_reserve: Request to bypass cache hold.
           - eligible_for_cache_reserve: Eligible for cache retention.
* `check_presence_cookie` - (Optional) When generating the cache key, check whether the cookie exists, and if so, add the cookie name to the cache key (the cookie name is not case sensitive). Multiple cookie names are supported, with multiple values separated by spaces.
* `check_presence_header` - (Optional) When the cache key is generated, check whether the header exists. If the header exists, add the header name to the cache key (the header name is not case sensitive). You can enter multiple header names, with multiple values separated by spaces.
* `edge_cache_mode` - (Optional) Node cache mode. value:
           - follow_origin: follows the Origin cache policy (if it exists), otherwise the default cache policy is used.
           - no_cache: Not cached.
           - override_origin: replaces the Origin cache policy.
           - follow_origin_bypass: follows the Origin cache policy (if it exists), otherwise it is not cached.
* `edge_cache_ttl` - (Optional) The node cache expiration time, in seconds.
* `edge_status_code_cache_ttl` - (Optional) Status code cache expiration time, in seconds.
* `include_cookie` - (Optional) When generating a cache key, it includes the specified cookie name and its value. You can enter multiple values separated by spaces.
* `include_header` - (Optional) When generating a cache key, it includes the specified header name and its value. You can enter multiple values separated by spaces.
* `query_string` - (Optional) The query string to be retained or deleted. You can enter multiple values separated by spaces.
* `query_string_mode` - (Optional) The processing mode for the query string when the cache key is generated. value:
           - ignore_all: Ignore all.
           - exclude_query_string: deletes the specified query string.
           - reserve_all: default value, all reserved.
           - include_query_string: Retains the specified query string.
* `rule` - (Optional) The rule content.
* `rule_enable` - (Optional) Rule switch. value:
           - on: open.
           - off: off.
* `rule_name` - (Optional) Rule name, you can find out the rule whose rule name is the passed field.
* `serve_stale` - (Optional) Response expiration cache. After enabling, nodes can still use cached expired files to respond to user requests even if the source server is unavailable. value:
           - on: open.
           - off: off.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the [ListSites](~~ ListSites ~~) API.
* `site_version` - (Optional, ForceNew, Int) Version number of the site configuration. For a site with configuration version management enabled, you can use this parameter to specify the site version in which the configuration takes effect. The default version is 0.
* `sort_query_string_for_cache` - (Optional) Query string sorting, which is disabled by default. value:
           - on: open.
           - off: off.
* `user_device_type` - (Optional) When generating the cache key, add the client device type. value:
           - on: open.
           - off: off.
* `user_geo` - (Optional) When generating the cache key, add the client geographic location. value:
           - on: open.
           - off: off.
* `user_language` - (Optional) When generating the cache key, add the client language type. value:
           - on: open.
           - off: off.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<cache_rule_id>`.
* `cache_rule_id` - The configured ConfigId. You can call the [ListCacheRules](~~ ListCacheRules ~~) operation to obtain the ConfigId.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cache Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Cache Rule.
* `update` - (Defaults to 5 mins) Used when update the Cache Rule.

## Import

ESA Cache Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_cache_rule.example <site_id>:<cache_rule_id>
```