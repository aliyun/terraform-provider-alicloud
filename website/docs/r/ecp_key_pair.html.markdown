---
subcategory: "Elastic Cloud Phone (ECP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecp_key_pair"
sidebar_current: "docs-alicloud-resource-ecp-key-pair"
description: |-
  Provides a Alicloud Elastic Cloud Phone (ECP) Key Pair resource.
---

# alicloud\_ecp\_key\_pair

Provides a Elastic Cloud Phone (ECP) Key Pair resource.

For information about Elastic Cloud Phone (ECP) Key Pair and how to use it, see [What is Key Pair](https://next.api.aliyun.com/document/cloudphone/2020-12-30/ImportImage).

-> **NOTE:** Available since v1.130.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecp_key_pair&exampleId=c9d75f6d-9ddb-d0fb-a621-3bc08f3487041e3f24ac&activeTab=example&spm=docs.r.ecp_key_pair.0.c9d75f6d9d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ecp_key_pair" "example" {
  key_pair_name   = "my-KeyPair"
  public_key_body = "ssh-rsa AAAAxxxxxxxxxxtyuudsfsg"
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecp_key_pair&spm=docs.r.ecp_key_pair.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `key_pair_name` - (Required, ForceNew) The Key Name.
* `public_key_body` - (Required) The public key body.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Key Pair. Its value is same as `key_pair_name`.

## Import

Elastic Cloud Phone (ECP) Key Pair can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecp_key_pair.example <key_pair_name>
```
