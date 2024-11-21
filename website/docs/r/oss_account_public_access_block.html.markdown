---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_account_public_access_block"
description: |-
  Provides a Alicloud OSS Account Public Access Block resource.
---

# alicloud_oss_account_public_access_block

Provides a OSS Account Public Access Block resource. Blocking public access at the account level.

For information about OSS Account Public Access Block and how to use it, see [What is Account Public Access Block](https://www.alibabacloud.com/help/en/oss/developer-reference/putpublicaccessblock).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_account_public_access_block&exampleId=4fd29403-b452-6737-a26a-3f322f57e97bd53854d4&activeTab=example&spm=docs.r.oss_account_public_access_block.0.4fd29403b4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_oss_account_public_access_block" "default" {
  block_public_access = true
}
```

## Argument Reference

The following arguments are supported:
* `block_public_access` - (Required) Whether or not AlibabaCloud OSS should block public bucket policies for buckets in this account is enabled.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account Public Access Block.
* `delete` - (Defaults to 5 mins) Used when delete the Account Public Access Block.
* `update` - (Defaults to 5 mins) Used when update the Account Public Access Block.

## Import

OSS Account Public Access Block can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_account_public_access_block.example 
```