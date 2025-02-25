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

resource "alicloud_esa_site" "resource_HttpBasicConfiguration_set_example" {
  site_name   = "gositecdn.cn"
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
* `ciphersuite` - (Optional) Ciphersuite
* `ciphersuite_group` - (Optional) CiphersuiteGroup
* `http2` - (Optional) https enable
* `http3` - (Optional) https enable
* `https` - (Optional) https enable
* `ocsp_stapling` - (Optional) OCSP enable
* `rule` - (Optional) rule
* `rule_enable` - (Optional) rule enable
* `rule_name` - (Optional) rule name
* `site_id` - (Required, ForceNew, Int) Site ID
* `tls10` - (Optional) Tls10 enable
* `tls11` - (Optional) Tls11 enable
* `tls12` - (Optional) Tls12 enable
* `tls13` - (Optional) Tls13 enable

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config ID

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Https Basic Configuration.
* `delete` - (Defaults to 5 mins) Used when delete the Https Basic Configuration.
* `update` - (Defaults to 5 mins) Used when update the Https Basic Configuration.

## Import

ESA Https Basic Configuration can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_https_basic_configuration.example <site_id>:<config_id>
```