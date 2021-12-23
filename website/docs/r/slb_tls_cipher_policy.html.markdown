---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_tls_cipher_policy"
sidebar_current: "docs-alicloud-resource-slb-tls-cipher-policy"
description: |-
  Provides a Alicloud SLB Tls Cipher Policy resource.
---

# alicloud\_slb\_tls\_cipher\_policy

Provides a SLB Tls Cipher Policy resource.

For information about SLB Tls Cipher Policy and how to use it, see [What is Tls Cipher Policy](https://www.alibabacloud.com/help/doc-detail/196714.htm).

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_slb_tls_cipher_policy" "example" {
  tls_cipher_policy_name = "Test-example_value"
  tls_versions           = ["TLSv1.2"]
  ciphers                = ["AES256-SHA256", "AES128-GCM-SHA256"]
}
```

## Argument Reference

The following arguments are supported:

* `tls_cipher_policy_name` - (Required) TLS policy name. Length is from 2 to 128, or in both the English and Chinese characters must be with an uppercase/lowercase letter or a Chinese character and the beginning, may contain numbers, in dot `.`, underscore `_` or dash `-`.
* `tls_versions` - (Required) The version of TLS protocol. You can find the corresponding value description in the document center [What is Tls Cipher Policy](https://www.alibabacloud.com/help/doc-detail/196714.htm).
* `ciphers` - (Required) The encryption algorithms supported. It depends on the value of `tls_versions`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Tls Cipher Policy.
* `status` - TLS policy instance state.

## Import

SLB Tls Cipher Policy can be imported using the id, e.g.

```
$ terraform import alicloud_slb_tls_cipher_policy.example <id>
```
