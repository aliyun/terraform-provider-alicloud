---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_ciphertext"
sidebar_current: "docs-alicloud-resource-kms-ciphertext"
description: |-
  Encrypt data with KMS.
---

# alicloud_kms_ciphertext

Encrypt a given plaintext with KMS. The produced ciphertext stays stable across applies. If the plaintext should be re-encrypted on each apply use the [`alicloud_kms_ciphertext`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/kms_ciphertext) data source.

-> **NOTE:** Available since v1.63.0.

-> **NOTE**: Using this data provider will allow you to conceal secret data within your resource definitions but does not take care of protecting that data in all Terraform logging and state output. Please take care to secure your secret data beyond just the Terraform configuration.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kms_ciphertext&exampleId=089d6344-261f-2331-91ce-8f1e925d785989bd2ff0&activeTab=example&spm=docs.r.kms_ciphertext.0.089d634426&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_kms_key" "key" {
  description            = "example key"
  status                 = "Enabled"
  pending_window_in_days = 7
}

resource "alicloud_kms_ciphertext" "encrypted" {
  key_id    = alicloud_kms_key.key.id
  plaintext = "example"
}
```

## Argument Reference

The following arguments are supported:

* `plaintext` - (Required, ForceNew) The plaintext to be encrypted which must be encoded in Base64.
* `key_id` - (Required, ForceNew) The globally unique ID of the CMK.
* `encryption_context` - (Optional, ForceNew) The Encryption context. If you specify this parameter here, it is also required when you call the Decrypt API operation. For more information, see [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ciphertext_blob` - The ciphertext of the data key encrypted with the primary CMK version.
