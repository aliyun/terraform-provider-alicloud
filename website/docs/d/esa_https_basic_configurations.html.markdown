---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_https_basic_configurations"
description: |-
  Provides a list of Esa Https Basic Configurations to the user.
---

# alicloud_esa_https_basic_configurations

This data source provides the Esa Https Basic Configurations of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.279.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_https_basic_configuration" "default" {
  site_id           = data.alicloud_esa_sites.default.sites.0.id
  rule_enable       = "on"
  rule_name         = var.name
  rule              = "true"
  ciphersuite       = "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256"
  ciphersuite_group = "all"
  tls10             = "on"
  tls11             = "on"
  tls12             = "on"
  tls13             = "on"
  ocsp_stapling     = "on"
  http2             = "on"
  http3             = "on"
  https             = "on"
  sequence          = 1
}

data "alicloud_esa_https_basic_configurations" "ids" {
  ids     = [alicloud_esa_https_basic_configuration.default.id]
  site_id = alicloud_esa_https_basic_configuration.default.site_id
}

output "esa_waf_https_basic_configurations_id_0" {
  value = data.alicloud_esa_https_basic_configurations.ids.configurations.0.id
}
```

## Argument Reference

The following attributes are exported:

* `ids` - (Optional, List) A list of Https Basic Configuration IDs.
* `name_regex` - (Optional) A regex string to filter results by Https Basic Configuration name.
* `site_id` - (Required) The ID of the Site.
* `config_id` - (Optional) The ID of the Configuration.
* `config_type` - (Optional) The type of the Configuration. Valid values: `global`, `rule`.
* `rule_name` - (Optional) The name of the rule.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Https Basic Configuration names.
* `configurations` - A list of Https Basic Configurations. Each element contains the following attributes:
  * `id` - The ID of the Https Basic Configuration.
  * `config_id` - The ID of the Configuration.
  * `config_type` - The type of the Configuration.
  * `rule_enable` - Rule switch.
  * `rule_name` - Rule name.
  * `rule` - Rule content.
  * `sequence` - The rule execution order prioritizes lower numerical values.
  * `ciphersuite` - Custom cipher suite, indicating the specific encryption algorithm selected when CiphersuiteGroup is set to custom.
  * `ciphersuite_group` - Cipher suite group.
  * `ocsp_stapling` - Whether to enable OCSP.
  * `http2` - Whether to enable HTTP2.
  * `http3` - Whether to enable HTTP3.
  * `https` - Whether to enable HTTPS.
  * `tls10` - Whether to enable TLS1.0.
  * `tls11` - Whether to enable TLS1.1.
  * `tls12` - Whether to enable TLS1.2.
  * `tls13` - Whether to enable TLS1.3.
