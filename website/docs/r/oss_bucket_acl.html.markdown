---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_acl"
description: |-
  Provides a Alicloud OSS Bucket Acl resource.
---

# alicloud_oss_bucket_acl

Provides a OSS Bucket Acl resource. The Access Control List (ACL) of a specific bucket.

For information about OSS Bucket Acl and how to use it, see [What is Bucket Acl](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketacl).

-> **NOTE:** Available since v1.220.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_acl&exampleId=c03819bd-4f62-ea6a-f9be-631dffb63da23ea61b1c&activeTab=example&spm=docs.r.oss_bucket_acl.0.c03819bd4f&intl_lang=EN_US" target="_blank">
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

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_oss_bucket_acl" "default" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
  acl    = "private"
}
```

### Deleting `alicloud_oss_bucket_acl` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_acl`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `acl` - (Required) Bucket-level Access Control List (ACL)，Valid values: `private`, `public-read`, `public-read-write`.
* `bucket` - (Required, ForceNew) The name of the bucket to which the current ACL configuration belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Acl.
* `update` - (Defaults to 5 mins) Used when update the Bucket Acl.

## Import

OSS Bucket Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_acl.example <id>
```