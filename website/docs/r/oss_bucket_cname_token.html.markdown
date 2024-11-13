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

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_oss_bucket" "defaultWWM58I" {
  bucket        = var.name
  storage_class = "Standard"
}


resource "alicloud_oss_bucket_cname_token" "default" {
  bucket = alicloud_oss_bucket.defaultWWM58I.bucket
  domain = "dinary.top"
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