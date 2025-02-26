---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_https_application_configuration"
description: |-
  Provides a Alicloud ESA Https Application Configuration resource.
---

# alicloud_esa_https_application_configuration

Provides a ESA Https Application Configuration resource.



For information about ESA Https Application Configuration and how to use it, see [What is Https Application Configuration](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateHttpsApplicationConfiguration).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "httpsapplicationconfiguration.example.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "domestic"
  access_type = "NS"
}

resource "alicloud_esa_https_application_configuration" "default" {
  hsts_max_age            = "31536000"
  alt_svc_clear           = "off"
  rule                    = "http.host eq \"video.example.com\""
  https_force             = "off"
  alt_svc_ma              = "86400"
  hsts                    = "off"
  rule_name               = "rule_example"
  rule_enable             = "off"
  site_id                 = alicloud_esa_site.default.id
  alt_svc_persist         = "off"
  hsts_preload            = "off"
  hsts_include_subdomains = "off"
  alt_svc                 = "off"
  https_force_code        = "301"
  site_version            = "0"
}
```

## Argument Reference

The following arguments are supported:
* `alt_svc` - (Optional) Function switch, default off. Value range:
  - on: open.
  - off: off.
* `alt_svc_clear` - (Optional) Alt-Svc whether The header contains the clear parameter. This parameter is disabled by default. Value range:
  - on: open.
  - off: off.
* `alt_svc_ma` - (Optional) The effective time of the Alt-Svc, in seconds. The default value is 86400 seconds.
* `alt_svc_persist` - (Optional) Alt-Svc whether The header contains the persist parameter. This parameter is disabled by default. Value range:
  - on: open.
  - off: off.
* `hsts` - (Optional) Whether to enable HSTS. It is disabled by default. Value range:
  - on: open.
  - off: off.
* `hsts_include_subdomains` - (Optional) Whether to include subdomains in HSTS is disabled by default. Value range:
  - on: open.
  - off: off.
* `hsts_max_age` - (Optional) The expiration time of HSTS, in seconds.
* `hsts_preload` - (Optional) Whether to enable HSTS preloading. It is disabled by default. Value range:
  - on: open.
  - off: off.
* `https_force` - (Optional) Whether to enable forced HTTPS. It is disabled by default. Value range:
  - on: open.
  - off: off.
* `https_force_code` - (Optional) Forced HTTPS jump status code, value range:
  - 301
  - 302
  - 307
  - 308
* `rule` - (Optional) Rule Content.
* `rule_enable` - (Optional) Rule switch. Value range:
  - on: Open.
  - off: off.
* `rule_name` - (Optional) Rule name, you can find out the rule whose rule name is the passed field.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Optional, ForceNew, Int) The version of the website configurations.
        

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Https Application Configuration.
* `delete` - (Defaults to 5 mins) Used when delete the Https Application Configuration.
* `update` - (Defaults to 5 mins) Used when update the Https Application Configuration.

## Import

ESA Https Application Configuration can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_https_application_configuration.example <site_id>:<config_id>
```