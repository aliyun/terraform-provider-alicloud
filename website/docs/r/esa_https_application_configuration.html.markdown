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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_https_application_configuration&exampleId=75a33a9f-97bf-e4e4-fbba-10cb9f86ee147a44b6fc&activeTab=example&spm=docs.r.esa_https_application_configuration.0.75a33a9f97&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
  - `on`: on.
  - `off`: off.
* `alt_svc_clear` - (Optional) Alt-Svc whether The header contains the clear parameter. This parameter is disabled by default. Value range:
  - `on`: on.
  - `off`: off.
* `alt_svc_ma` - (Optional) The effective time of the Alt-Svc, in seconds. The default value is 86400 seconds.
* `alt_svc_persist` - (Optional) Alt-Svc whether The header contains the persist parameter. This parameter is disabled by default. Value range:
  - `on`: on.
  - `off`: off.
* `hsts` - (Optional) Whether to enable HSTS. It is disabled by default. Value range:
  - `on`: on.
  - `off`: off.
* `hsts_include_subdomains` - (Optional) Whether to include subdomains in HSTS is disabled by default. Value range:
  - `on`: on.
  - `off`: off.
* `hsts_max_age` - (Optional) The expiration time of HSTS, in seconds.
* `hsts_preload` - (Optional) Whether to enable HSTS preloading. It is disabled by default. Value range:
  - `on`: on.
  - `off`: off.
* `https_force` - (Optional) Whether to enable forced HTTPS. It is disabled by default. Value range:
  - `on`: on.
  - `off`: off.
* `https_force_code` - (Optional) Forced HTTPS jump status code, value range:
  - `301`
  - `302`
  - `307`
  - `308`
* `https_no_sni_deny` - (Optional, Available since v1.262.1) Whether to enable to reject TLS handshake requests without SNI. This parameter is disabled by default. Value range:
  - `on`: open.
  - `off`: off.
* `https_sni_verify` - (Optional, Available since v1.262.1) Whether to enable SNI verification. It is disabled by default. Value range:
  - `on`: open.
  - `off`: off.
* `https_sni_whitelist` - (Optional, Available since v1.262.1) Specifies the list of allowed SNI whitelists, separated by spaces.
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - `on`: open.
  - `off`: close.
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Int, Available since v1.262.1) The rule execution order prioritizes lower numerical values. It is only applicable when setting or modifying the order of individual rule configurations.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Https Application Configuration.
* `delete` - (Defaults to 5 mins) Used when delete the Https Application Configuration.
* `update` - (Defaults to 5 mins) Used when update the Https Application Configuration.

## Import

ESA Https Application Configuration can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_https_application_configuration.example <site_id>:<config_id>
```