---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_lifecycle"
description: |-
  Provides a Alicloud OSS Bucket Lifecycle resource.
---

# alicloud_oss_bucket_lifecycle

Provides a OSS Bucket Lifecycle resource. OSS lifecycle rules allow you to automatically transition objects between storage classes or delete them after a configurable period.

For information about OSS Bucket Lifecycle and how to use it, see [What is Bucket Lifecycle](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketlifecycle).

-> **NOTE:** Available since v1.295.0.

`alicloud_oss_bucket` keeps the embedded `lifecycle_rule` block for backward compatibility. `alicloud_oss_bucket_lifecycle` is the standalone resource that manages the same lifecycle configuration. When you adopt the standalone resource, add `ignore_changes = [lifecycle_rule]` on `alicloud_oss_bucket` so that the two resources do not fight over the same configuration, the same way `alicloud_oss_bucket_cors` / `alicloud_oss_bucket_logging` already coexist with the embedded blocks.

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

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
  lifecycle {
    # When you use `alicloud_oss_bucket_lifecycle`, you must add `ignore_changes`
    # for the `lifecycle_rule` attribute on `alicloud_oss_bucket` to avoid
    # unexpected diffs caused by both resources managing the same configuration.
    ignore_changes = [
      lifecycle_rule,
    ]
  }
}

resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    expiration {
      days = 365
    }
  }

  lifecycle_rule {
    id      = "rule2"
    prefix  = "path2/"
    enabled = true

    transitions {
      days          = 30
      storage_class = "IA"
    }

    abort_multipart_upload {
      days = 10
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the Bucket.
* `lifecycle_rule_allow_same_action_overlap` - (Optional) Specifies whether to allow lifecycle rules to overlap. Default value: `false`.
* `lifecycle_rule` - (Optional) The lifecycle rule configuration of the Bucket. See [`lifecycle_rule`](#lifecycle_rule) below.

### `lifecycle_rule`

The lifecycle_rule supports the following:

* `id` - (Optional, Computed) The ID of the lifecycle rule. The ID can be up to 255 characters in length.
* `prefix` - (Optional) The prefix that the lifecycle rule applies to. Only objects whose names start with the prefix are subject to the rule.
* `enabled` - (Required) Specifies whether to enable the lifecycle rule. Valid values: `true`, `false`.
* `tags` - (Optional) The tags of the lifecycle rule. Key-value pairs.
* `expiration` - (Optional) The expiration configuration of the lifecycle rule. See [`expiration`](#expiration) below.
* `transitions` - (Optional) The storage class transition configuration of the lifecycle rule. See [`transitions`](#transitions) below.
* `abort_multipart_upload` - (Optional) The AbortMultipartUpload configuration of the lifecycle rule. See [`abort_multipart_upload`](#abort_multipart_upload) below.
* `noncurrent_version_expiration` - (Optional) The expiration configuration of the noncurrent version objects. See [`noncurrent_version_expiration`](#noncurrent_version_expiration) below.
* `noncurrent_version_transition` - (Optional) The storage class transition configuration of the noncurrent version objects. See [`noncurrent_version_transition`](#noncurrent_version_transition) below.
* `filter` - (Optional) The filter configuration of the lifecycle rule. See [`filter`](#filter) below.

### `expiration`

The expiration supports the following:

* `date` - (Optional) The date when the rule takes effect. The date follows the ISO 8601 `YYYY-MM-DD` format. The date must be a future date. Only one of `date`, `created_before_date`, `days`, and `expired_object_delete_marker` can be specified.
* `created_before_date` - (Optional) The date before which the rule takes effect. The date follows the ISO 8601 `YYYY-MM-DD` format. Only one of `date`, `created_before_date`, `days`, and `expired_object_delete_marker` can be specified.
* `days` - (Optional) The number of days after which the rule takes effect after the object is last modified. Only one of `date`, `created_before_date`, `days`, and `expired_object_delete_marker` can be specified.
* `expired_object_delete_marker` - (Optional) Specifies whether to delete the delete marker of an object. Only one of `date`, `created_before_date`, `days`, and `expired_object_delete_marker` can be specified.

### `transitions`

The transitions supports the following:

* `created_before_date` - (Optional) The date before which the rule takes effect. The date follows the ISO 8601 `YYYY-MM-DD` format. Only one of `created_before_date` and `days` can be specified.
* `days` - (Optional) The number of days after which the rule takes effect after the object is last modified. Only one of `created_before_date` and `days` can be specified.
* `storage_class` - (Required) The storage class to which objects are transitioned. Valid values: `Standard`, `IA`, `Archive`, `ColdArchive`, `DeepColdArchive`.
* `is_access_time` - (Optional) Specifies whether to transition objects based on the access time.
* `return_to_std_when_visit` - (Optional) Specifies whether to transition objects back to Standard storage class when they are accessed.

### `abort_multipart_upload`

The abort_multipart_upload supports the following:

* `created_before_date` - (Optional) The date before which the rule takes effect. The date follows the ISO 8601 `YYYY-MM-DD` format. Only one of `created_before_date` and `days` can be specified.
* `days` - (Optional) The number of days after which the rule takes effect after the multipart upload is initiated. Only one of `created_before_date` and `days` can be specified.

### `noncurrent_version_expiration`

The noncurrent_version_expiration supports the following:

* `days` - (Required) The number of days after which noncurrent version objects expire.

### `noncurrent_version_transition`

The noncurrent_version_transition supports the following:

* `days` - (Required) The number of days after which noncurrent version objects are transitioned.
* `storage_class` - (Required) The storage class to which noncurrent version objects are transitioned. Valid values: `Standard`, `IA`, `Archive`, `ColdArchive`, `DeepColdArchive`.
* `is_access_time` - (Optional) Specifies whether to transition noncurrent version objects based on the access time.
* `return_to_std_when_visit` - (Optional) Specifies whether to transition noncurrent version objects back to Standard storage class when they are accessed.

### `filter`

The filter supports the following:

* `object_size_greater_than` - (Optional) The minimum object size (in bytes) that the lifecycle rule applies to. The value must be greater than 0.
* `object_size_less_than` - (Optional) The maximum object size (in bytes) that the lifecycle rule applies to. The value must be greater than 0.
* `not` - (Optional) The filter conditions that the lifecycle rule does not apply to. See [`not`](#not) below.

### `not`

The not supports the following:

* `prefix` - (Optional) The prefix that the lifecycle rule does not apply to.
* `tag` - (Optional) The tag that the lifecycle rule does not apply to. See [`tag`](#tag) below.

### `tag`

The tag supports the following:

* `key` - (Required) The key of the tag.
* `value` - (Required) The value of the tag.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Bucket Lifecycle resource. It is the same as the bucket name.

## Import

OSS Bucket Lifecycle can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_lifecycle.example <bucket-name>
```
