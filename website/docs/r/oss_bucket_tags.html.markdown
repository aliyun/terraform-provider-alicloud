---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_tags"
description: |-
  Provides a Alicloud OSS Bucket Tags resource.
---

# alicloud_oss_bucket_tags

Provides a OSS Bucket Tags resource. The complete set of tags (key/value pairs) attached to an OSS bucket.

PutBucketTags is an idempotent overwrite: this resource always reconciles the bucket's tag set to the configured `tags` map. Setting `tags` to an empty map `{}` clears all tags on the bucket via DeleteBucketTags.

For information about OSS Bucket Tags and how to use it, see [PutBucketTags](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutBucketTags).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** This resource manages the full tag set of the bucket: creating or updating it overwrites the bucket tags (PutBucketTags), and destroying it (or removing it from the configuration) clears all tags (DeleteBucketTags). The inline `tags` attribute of `alicloud_oss_bucket` manages the same tag set. If the same bucket is also managed by `alicloud_oss_bucket`, add `lifecycle { ignore_changes = [tags] }` to the bucket resource; otherwise `alicloud_oss_bucket` will keep overwriting the tags set by this resource on every apply. Do not manage the same bucket's tags with both resources at the same time.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

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
  lifecycle {
    # When you use `alicloud_oss_bucket_tags`, you must add `ignore_changes` for the `tags` attribute
    # on `alicloud_oss_bucket` to avoid unexpected diffs caused by both resources managing the same tag set.
    ignore_changes = [
      tags,
    ]
  }
}

resource "alicloud_oss_bucket_tags" "default" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
  tags = {
    Created = "TF"
    Env     = "test"
  }
}
```

### Deleting `alicloud_oss_bucket_tags` or removing it from your configuration

Destroying `alicloud_oss_bucket_tags` clears all tags from the bucket via DeleteBucketTags. The bucket itself is not deleted.

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew, String) The name of the bucket whose tag set is managed.
* `tags` - (Required, Map) The complete set of tags applied to the bucket. Up to 20 tag pairs are supported. Setting this to an empty map `{}` clears all tags on the bucket.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The ID of the resource, which is the bucket name.

## Import

OSS Bucket Tags can be imported using the bucket name as the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_tags.example <bucket-name>
```
