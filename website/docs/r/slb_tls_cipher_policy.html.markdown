---
subcategory: "Classic Load Balancer (SLB)"
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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_tls_cipher_policy&exampleId=3e716de4-15e8-7606-48c7-1de0162fe9474420b3df&activeTab=example&spm=docs.r.slb_tls_cipher_policy.0.3e716de415&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

```shell
$ terraform import alicloud_slb_tls_cipher_policy.example <id>
```
