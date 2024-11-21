---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_cname_token"
description: |-
  Provides a Alicloud OSS Bucket Cname Token resource.
---

# alicloud_oss_bucket_cname_token

Provides a OSS Bucket Cname Token resource.

The token used to verify the ownership of the bucket custom domain name.

For information about OSS Bucket Cname Token and how to use it, see [What is Bucket Cname Token](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.233.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_oss_bucket_cname_token&exampleId=1e4e64f5-e7ff-bc8d-f299-91a45975d9aac19c9f09&activeTab=example&spm=docs.r.oss_bucket_cname_token.0.1e4e64f5e7&intl_lang=EN_US" target="_blank">
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

resource "alicloud_oss_bucket" "CreateBucket" {
  bucket        = var.name
  storage_class = "Standard"
}

resource "alicloud_oss_bucket_cname_token" "defaultZaWJfG" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
  domain = "tftestacc.com"
}
```

### Deleting `alicloud_oss_bucket_cname_token` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_cname_token`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket
* `domain` - (Required, ForceNew) The custom domain

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<bucket>:<domain>`.
* `token` - Token used to verify domain ownership

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Cname Token.

## Import

OSS Bucket Cname Token can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_cname_token.example <bucket>:<domain>
```