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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Acl.
* `update` - (Defaults to 5 mins) Used when update the Bucket Acl.

## Import

OSS Bucket Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_acl.example <id>
```