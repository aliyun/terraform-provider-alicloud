---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_domain"
description: |-
  Provides a Alicloud APIG Domain resource.
---

# alicloud_apig_domain

Provides a APIG Domain resource.



For information about APIG Domain and how to use it, see [What is Domain](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateDomain).

-> **NOTE:** Available since v1.284.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_apig_domain" "default" {
  domain_name  = "example-domain-cspec-v6.example.com"
  gateway_type = "API"
  protocol     = "HTTP"
}
```

## Argument Reference

The following arguments are supported:
* `ca_cert_identifier` - (Optional) CA certificate identifier
* `cert_identifier` - (Optional) tls cert identifier
* `client_ca_cert` - (Optional) client CA certificate
* `domain_name` - (Required, ForceNew) domain name
* `domain_scope` - (Optional, Computed) domain scope
* `force_https` - (Optional) Set the HTTPS protocol type and whether to enable forced HTTPS redirection.
* `gateway_type` - (Optional, ForceNew) Gateway type. Valid values: `API`, `AI`.

* `http2_option` - (Optional) Whether to enable http2 settings.
* `m_tls_enabled` - (Optional) Whether to enable mTLS mutual authentication
* `protocol` - (Required) Protocol, HTTP/HTTPS
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `tls_cipher_suites_config` - (Optional, Computed, Set) TlsCipherSuitesConfig See [`tls_cipher_suites_config`](#tls_cipher_suites_config) below.
* `tls_max` - (Optional) The maximum version of the TLS protocol. The maximum version of the TLS protocol is 1.3.
* `tls_min` - (Optional) The minimum version of the TLS protocol. The minimum version of the TLS protocol is 1.0.

### `tls_cipher_suites_config`

The tls_cipher_suites_config supports the following:
* `config_type` - (Optional) config type, Default or Custom
* `tls_cipher_suite` - (Optional, List) tls Cipher Suite See [`tls_cipher_suite`](#tls_cipher_suites_config-tls_cipher_suite) below.

### `tls_cipher_suites_config-tls_cipher_suite`

The tls_cipher_suites_config-tls_cipher_suite supports the following:
* `name` - (Optional) cipher suite name
* `support_versions` - (Optional, List) support versions

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Domain.
* `update` - (Defaults to 5 mins) Used when update the Domain.

## Import

APIG Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_domain.example <domain_id>
```