---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_site"
description: |-
  Provides a Alicloud ESA Site resource.
---

# alicloud_esa_site

Provides a ESA Site resource.



For information about ESA Site and how to use it, see [What is Site](https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/user-guide/site-management).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_site&exampleId=526a6ad2-d747-8099-80fe-4adba67e465727294d3a&activeTab=example&spm=docs.r.esa_site.0.526a6ad2d7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_esa_rate_plan_instance" "defaultIEoDfU" {
  type         = "NS"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "basic"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_esa_site" "default" {
  site_name         = "bcd${random_integer.default.result}.com"
  coverage          = "overseas"
  access_type       = "NS"
  instance_id       = alicloud_esa_rate_plan_instance.defaultIEoDfU.id
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_site&spm=docs.r.esa_site.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `access_type` - (Required, ForceNew) The DNS setup. Valid values:
  - `NS`
  - `CNAME`
* `add_client_geolocation_header` - (Optional, Available since v1.244.0) Add visitor geolocation header. Value range:
  - `on`: Enable.
  - `off`: Disable.
* `add_real_client_ip_header` - (Optional, Available since v1.244.0) Add the "ali-real-client-ip" header containing the real client IP. Value range:
  - `on`: Enable.
  - `off`: Disable.
* `ai_mode` - (Optional, Computed, Available since v1.272.1) HTTP DDoS Intelligent Protection Mode. Valid values:
  - `observe`: Observe.
  - `defense`: Block.
* `ai_template` - (Optional, Computed, Available since v1.272.1) HTTP DDoS Intelligent Protection Level. Values:
  - `level0`: Very Loose.
  - `level30`: Loose.
  - `level60`: Normal.
  - `level90`: Strict.
* `automatic_frequency_control_action_type` - (Optional, Available since v1.276.0) AutomaticFrequencyControl Disposal action.Valid values:
  - `observe`: Observe.
  - `deny`: Block.
  - `js`: Js verification.
* `automatic_frequency_control_enable` - (Optional, Available since v1.276.0) AutomaticFrequencyControl Switch.Valid values:
  - `on`: on.
  - `off`: off.
* `automatic_frequency_control_level` - (Optional, Available since v1.276.0) AutomaticFrequencyControl Protection Level.Valid values:
  - `loose`: Loose.
  - `normal`: Normal.
  - `strict`: Strict.
* `cache_architecture_mode` - (Optional, Computed, Available since v1.244.0) Multi-level cache architecture mode. Possible values:
  - `edge`: Edge cache layer.
  - `edge_smart`: Edge cache layer + intelligent cache layer.
  - `edge_regional`: Edge cache layer + regional cache layer.
  - `edge_regional_smart`: Edge cache layer + regional cache layer + intelligent cache layer.
* `cache_reserve_enable` - (Optional, Available since v1.251.0) Site cache hold switch. Value:
  - `on`
  - `off`
* `cache_reserve_instance_id` - (Optional, Available since v1.251.0) The ID of the cache instance.
* `case_insensitive` - (Optional, Available since v1.251.0) Whether the custom CacheTag name ignores case. Value range:
  - `on`:  Enabled, ignores case.
  - `off`: Disabled, does not ignore case.
* `coverage` - (Optional) The service location. Valid values:
  - `domestic`: the Chinese mainland
  - `global`: global
  - `overseas`: outside the Chinese mainland
* `cross_border_optimization` - (Optional, Available since v1.251.0) Whether to enable mainland China network access optimization, the default is off. Value range:
  - `on`
  - `off`
* `development_mode` - (Optional, Available since v1.251.0) Development mode configuration function switch. Value range:
  - `on`
  - `off`
* `flatten_mode` - (Optional, Available since v1.251.0) CNAME flattening mode. Possible values:
  - `flatten_all`: Flatten all.
  - `flatten_at_root`: Flatten only the root domain. The default is to flatten the root domain.
* `global_mode` - (Optional, Computed, Available since v1.272.1) HTTP DDoS Attack Protection Policy Modes. Valid values:
  - `very weak`: indicates a very permissive setting.
  - `weak`: indicates a permissive setting.
  - `default`: indicates a normal setting.
  - `hard`: indicates a strict setting.
* `instance_id` - (Required, ForceNew) The ID of the associated package instance.
* `ipv6_enable` - (Optional, Computed, Available since v1.244.0) Specifies whether to enable IPv6. Valid values:
  - `on`
  - `off`
* `ipv6_region` - (Optional, Computed, Available since v1.251.0) The region in which Ipv6 is enabled. The default value is x.x:
  - 'x.x': Global.
  - 'Cn.cn ': Mainland China.
* `paused` - (Optional, Available since v1.266.0) Specifies whether to temporarily pause ESA on the website. If you set this parameter to true, all requests to the domains in your DNS records go directly to your origin server. Valid values:
  - `true`
  - `false`
* `performance_data_collection_enable` - (Optional, Available since v1.276.0) Whether to enable the quality data collection switch.Valid values:
  - `on`: on.
  - `off`: off.
* `real_client_ip_header_name` - (Optional, Available since v1.276.0) Real client IP header name.
* `resource_group_id` - (Optional, ForceNew, Computed) The ID of the resource group
* `seo_bypass` - (Optional, Available since v1.251.0) Release the search engine crawler configuration. Value:
  - `on`
  - `off`
* `site_name` - (Required, ForceNew) The website name.
* `site_name_exclusive` - (Optional, Available since v1.251.0) Specifies whether to enable site hold.After you enable site hold, other accounts cannot add your website domain or its subdomains to ESA. Valid values:
  - `on`
  - `off`
* `site_version` - (Optional, Int, Available since v1.244.0) The version number of the site. For sites with version management enabled, you can use this parameter to specify the site version for which the configuration will take effect, defaulting to version 0.
* `site_waf_settings` - (Optional, Set, Available since v1.276.0) Site WAF Configuration Details. See [`site_waf_settings`](#site_waf_settings) below.
* `tag_name` - (Optional, Available since v1.251.0) Custom CacheTag name.
* `tags` - (Optional, Map) Resource tags
* `version_management` - (Optional) Version management enabled. When true, version management is turned on for the table site.

### `site_waf_settings`

The site_waf_settings supports the following:
* `add_bot_protection_headers` - (Optional, Set, Available since v1.276.0) Add BOT Protection Header. See [`add_bot_protection_headers`](#site_waf_settings-add_bot_protection_headers) below.
* `add_security_headers` - (Optional, Set, Available since v1.276.0) Add Security Header. See [`add_security_headers`](#site_waf_settings-add_security_headers) below.
* `bandwidth_abuse_protection` - (Optional, Set, Available since v1.276.0) Anti-theft brush. See [`bandwidth_abuse_protection`](#site_waf_settings-bandwidth_abuse_protection) below.
* `bot_management` - (Optional, Set, Available since v1.276.0) Bot Management. See [`bot_management`](#site_waf_settings-bot_management) below.
* `client_ip_identifier` - (Optional, Set, Available since v1.276.0) Client IP Identification. See [`client_ip_identifier`](#site_waf_settings-client_ip_identifier) below.
* `security_level` - (Optional, Set, Available since v1.276.0) Security Level. See [`security_level`](#site_waf_settings-security_level) below.

### `site_waf_settings-add_bot_protection_headers`

The site_waf_settings-add_bot_protection_headers supports the following:
* `enable` - (Optional, Available since v1.276.0) Switch.

### `site_waf_settings-add_security_headers`

The site_waf_settings-add_security_headers supports the following:
* `enable` - (Optional, Available since v1.276.0) Switch.

### `site_waf_settings-bandwidth_abuse_protection`

The site_waf_settings-bandwidth_abuse_protection supports the following:
* `action` - (Optional, Available since v1.276.0) Anti-theft brush Rule Action.Valid values:
  - `deny`: Deny.
  - `monitor`: Monitor.
  - `captcha`: Captcha.
* `status` - (Optional, Available since v1.276.0) Anti-theft brush Rule Status.

### `site_waf_settings-bot_management`

The site_waf_settings-bot_management supports the following:
* `definite_bots` - (Optional, Set, Available since v1.276.0) Definitely Bot. See [`definite_bots`](#site_waf_settings-bot_management-definite_bots) below.
* `effect_on_static` - (Optional, Set, Available since v1.276.0) VApply to Static Resource Requests. See [`effect_on_static`](#site_waf_settings-bot_management-effect_on_static) below.
* `js_detection` - (Optional, Set, Available since v1.276.0) JavaScript Challenge. See [`js_detection`](#site_waf_settings-bot_management-js_detection) below.
* `likely_bots` - (Optional, Set, Available since v1.276.0) Likely Bot. See [`likely_bots`](#site_waf_settings-bot_management-likely_bots) below.
* `verified_bots` - (Optional, Set, Available since v1.276.0) Verified Bot. See [`verified_bots`](#site_waf_settings-bot_management-verified_bots) below.

### `site_waf_settings-client_ip_identifier`

The site_waf_settings-client_ip_identifier supports the following:
* `headers` - (Optional, List, Available since v1.276.0) Specify Headers.
* `mode` - (Optional, Available since v1.276.0) Identification Mode.Valid values:
  - `headers`: Specify headers.
  - `connection_ip`: Connection IP.

### `site_waf_settings-security_level`

The site_waf_settings-security_level supports the following:
* `value` - (Optional, Available since v1.276.0) Security level value. Valid values:
  - `high`: High.
  - `low`: Low.
  - `under_attack`: I'm under attack.
  - `medium`: Medium.
  - `essentially_off`: Essentially off.
  - `off`: Completely off.

### `site_waf_settings-bot_management-definite_bots`

The site_waf_settings-bot_management-definite_bots supports the following:
* `action` - (Optional, Available since v1.276.0) Action.Valid values:
  - `allow`: Allow.
  - `deny`: Deny.
  - `monitor`: Monitor.
  - `captcha`: Captcha.

### `site_waf_settings-bot_management-effect_on_static`

The site_waf_settings-bot_management-effect_on_static supports the following:
* `enable` - (Optional, Available since v1.276.0) Switch.

### `site_waf_settings-bot_management-js_detection`

The site_waf_settings-bot_management-js_detection supports the following:
* `enable` - (Optional, Available since v1.276.0) Switch.

### `site_waf_settings-bot_management-likely_bots`

The site_waf_settings-bot_management-likely_bots supports the following:
* `action` - (Optional, Available since v1.276.0) Action.Valid values:
  - `allow`: Allow.
  - `deny`: Deny.
  - `monitor`: Monitor.
  - `captcha`: Captcha.

### `site_waf_settings-bot_management-verified_bots`

The site_waf_settings-bot_management-verified_bots supports the following:
* `action` - (Optional, Available since v1.276.0) Action.Valid values:
  - `allow`: Allow.
  - `deny`: Deny.
  - `monitor`: Monitor.
  - `captcha`: Captcha.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the website was added. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `site_waf_settings` - Site WAF Configuration Details.
  * `bandwidth_abuse_protection` - Anti-theft brush.
    * `id` - Anti-theft brush Rule ID.
  * `bot_management` - Bot Management.
    * `definite_bots` - Definitely Bot.
      * `id` - Rule ID.
    * `likely_bots` - Likely Bot.
      * `id` - Rule ID.
    * `verified_bots` - Verified Bot.
      * `id` - Rule ID.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 45 mins) Used when create the Site.
* `delete` - (Defaults to 15 mins) Used when delete the Site.
* `update` - (Defaults to 5 mins) Used when update the Site.

## Import

ESA Site can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_site.example <site_id>
```
