---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_https_basic_configuration"
description: |-
  Provides a Alicloud ESA Https Basic Configuration resource.
---

# alicloud_esa_https_basic_configuration

Provides a ESA Https Basic Configuration resource.



For information about ESA Https Basic Configuration and how to use it, see [What is Https Basic Configuration](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateHttpsBasicConfiguration).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_esa_rate_plan_instance" "example" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_esa_site" "resource_HttpBasicConfiguration_set_example" {
  site_name   = "gositecdn-${random_integer.default.result}.cn"
  instance_id = alicloud_esa_rate_plan_instance.example.id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_https_basic_configuration" "default" {
  https       = "on"
  rule        = "true"
  rule_name   = "example2"
  site_id     = alicloud_esa_site.resource_HttpBasicConfiguration_set_example.id
  rule_enable = "on"
}
```

## Argument Reference

The following arguments are supported:
* `ciphersuite` - (Optional) Custom cipher suite, indicating the specific encryption algorithm selected when CiphersuiteGroup is set to custom.
* `ciphersuite_group` - (Optional) Cipher suite group. Default is all cipher suites. Possible values:
  - all: All cipher suites.
  - strict: Strong cipher suites.
  - custom: Custom cipher suites.
* `http2` - (Optional) Indicates whether HTTP2 is enabled. Default is on. Possible values:
  - on: Enabled.
  - off: Disabled.
* `http3` - (Optional) Whether to enable HTTP3, which is enabled by default. The value can be:
  - on: Enabled. 
  - off: Disabled.
* `https` - (Optional) Whether to enable HTTPS. Default is enabled. Possible values:
  - on: Enable.
  - off: Disable.
* `ocsp_stapling` - (Optional) Indicates whether OCSP is enabled. Default is off. Possible values:
  - on: Enabled.
  - off: Disabled.
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  -  Match all incoming requests: value set to true
  -  Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - on: open.
  - off: close.
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Int, Available since v1.263.0) The rule execution order prioritizes lower numerical values. It is only applicable when setting or modifying the order of individual rule configurations.
* `site_id` - (Required, ForceNew, Int) Site ID, which can be obtained by calling the [ListSites](~~ListSites~~) interface.
* `tls10` - (Optional) Whether to enable TLS1.0. Default is disabled. Possible values:
  - on: Enable.
  - off: Disable.
* `tls11` - (Optional) Whether to enable TLS1.1. Default is enabled. Possible values:
  - on: Enable.
  - off: Disable.
* `tls12` - (Optional) Whether to enable TLS1.2. Default is enabled. Possible values:
  - on: Enable.
  - off: Disable.
* `tls13` - (Optional) Whether to enable TLS1.3. Default is enabled. Possible values:
  - on: Enable.
  - off: Disable.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - ConfigId of the configuration, which can be obtained by calling the [ListHttpsBasicConfigurations](https://www.alibabacloud.com/help/en/doc-detail/2867470.html) interface.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Https Basic Configuration.
* `delete` - (Defaults to 5 mins) Used when delete the Https Basic Configuration.
* `update` - (Defaults to 5 mins) Used when update the Https Basic Configuration.

## Import

ESA Https Basic Configuration can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_https_basic_configuration.example <site_id>:<config_id>
```