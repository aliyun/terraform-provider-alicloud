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

## Argument Reference

The following arguments are supported:
* `access_type` - (Optional, ForceNew) The DNS setup. Valid values:
  - `NS`
  - `CNAME`
* `add_client_geolocation_header` - (Optional, Available since v1.244.0) Add visitor geolocation header. Value range:
  - `on`: Enable.
  - `off`: Disable.
* `add_real_client_ip_header` - (Optional, Available since v1.244.0) Add the "ali-real-client-ip" header containing the real client IP. Value range:
  - `on`: Enable.
  - `off`: Disable.
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
* `resource_group_id` - (Optional, ForceNew, Computed) The ID of the resource group
* `seo_bypass` - (Optional, Available since v1.251.0) Release the search engine crawler configuration. Value:
  - `on`
  - `off`
* `site_name` - (Required, ForceNew) The website name.
* `site_name_exclusive` - (Optional, Available since v1.251.0) Specifies whether to enable site hold.After you enable site hold, other accounts cannot add your website domain or its subdomains to ESA. Valid values:
  - `on`
  - `off`
* `site_version` - (Optional, Int, Available since v1.244.0) The version number of the site. For sites with version management enabled, you can use this parameter to specify the site version for which the configuration will take effect, defaulting to version 0.
* `tag_name` - (Optional, Available since v1.251.0) Custom CacheTag name.
* `tags` - (Optional, Map) Resource tags
* `version_management` - (Optional) Version management enabled. When true, version management is turned on for the table site.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the website was added. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Site.
* `delete` - (Defaults to 5 mins) Used when delete the Site.
* `update` - (Defaults to 5 mins) Used when update the Site.

## Import

ESA Site can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_site.example <id>
```