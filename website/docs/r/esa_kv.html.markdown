---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_kv"
description: |-
  Provides a Alicloud ESA Kv resource.
---

# alicloud_esa_kv

Provides a ESA Kv resource.



For information about ESA Kv and how to use it, see [What is Kv](https://next.api.alibabacloud.com/document/ESA/2024-09-10/PutKv).

-> **NOTE:** Available since v1.251.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_kv&exampleId=cb271c28-3f51-69f5-a8dd-586c35fb7afd83436251&activeTab=example&spm=docs.r.esa_kv.0.cb271c283f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_esa_kv_namespace" "default" {
  description  = "this is a example namespace."
  kv_namespace = "namespace1"
}

resource "alicloud_esa_kv" "default" {
  isbase         = "false"
  expiration_ttl = "360"
  value          = "example_value"
  expiration     = "1690"
  namespace      = alicloud_esa_kv_namespace.default.id
  key            = "example_key"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_kv&spm=docs.r.esa_kv.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `expiration` - (Optional, Int) The content of the key, which can be up to 2 MB (2 × 1000 × 1000). If the content is larger than 2 MB, call [PutKvWithHighCapacity](https://www.alibabacloud.com/help/en/doc-detail/2850486.html).
* `expiration_ttl` - (Optional, Int) The time when the key-value pair expires, which cannot be earlier than the current time. The value is a timestamp in seconds. If you specify both Expiration and ExpirationTtl, only ExpirationTtl takes effect.
* `isbase` - (Optional) The relative expiration time. Unit: seconds. If you specify both Expiration and ExpirationTtl, only ExpirationTtl takes effect.
* `key` - (Required, ForceNew) kv
* `namespace` - (Required, ForceNew) The name specified when calling [CreatevNamespace](https://help.aliyun.com/document_detail/2850317.html).
* `url` - (Optional) The key name. The name can be up to 512 characters in length and cannot contain spaces or backslashes (\\).
* `value` - (Optional) The content of the key. If the content has more than 256 characters in length, the system displays the first 100 and the last 100 characters, and omits the middle part.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<namespace>:<key>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Kv.
* `delete` - (Defaults to 5 mins) Used when delete the Kv.
* `update` - (Defaults to 5 mins) Used when update the Kv.

## Import

ESA Kv can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_kv.example <namespace>:<key>
```