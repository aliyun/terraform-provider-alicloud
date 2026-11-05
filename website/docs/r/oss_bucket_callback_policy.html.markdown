---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_callback_policy"
description: |-
  Provides a Alicloud OSS Bucket Callback Policy resource.
---

# alicloud_oss_bucket_callback_policy

Provides a OSS Bucket Callback Policy resource. The bucket callback policy pre-defines the callback configuration that OSS applies to upload requests against the bucket, so that clients no longer have to send callback parameters with every upload.

For information about OSS Bucket Callback Policy and how to use it, see [PutBucketCallbackPolicy](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutBucketCallbackPolicy).

-> **NOTE:** Available since v2.3.0.

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

resource "alicloud_oss_bucket" "default" {
  bucket        = var.name
  storage_class = "Standard"
}

resource "alicloud_oss_bucket_callback_policy" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  policy_item {
    policy_name  = "name1"
    # callback and callback_var are Base64-encoded JSON documents.
    callback     = "eyJjYWxsYmFja1VybCI6InRlc3QuY29tL3Rlc3QucGhwIiwgImNhbGxiYWNrQm9keSI6ImJ1Y2tldD0ke2J1Y2tldH0mb2JqZWN0PSR7b2JqZWN0fSZhPSR7eDphfSZiPSR7eDpifSJ9"
    callback_var = "eyJ4OmEiOiJhIiwgIng6YiI6ImIifQ=="
  }
  policy_item {
    policy_name  = "name2"
    callback     = "eyJjYWxsYmFja1VybCI6InRlc3QyLmNvbS90ZXN0LnBocCIsICJjYWxsYmFja0JvZHkiOiJidWNrZXQ9JHtidWNrZXR9Jm9iamVjdD0ke29iamVjdH0mYT0ke3g6YX0mYj0ke3g6Yn0ifQ=="
    callback_var = "eyJ4OmEiOiJhIiwgIng6YiI6ImQifQ=="
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the bucket.
* `policy_item` - (Optional, List) The container that holds the callback policy configuration items. See [`policy_item`](#policy_item) below.

### `policy_item`

The policy_item supports the following:

* `callback` - (Optional) Callback address. A Base64-encoded JSON document containing the `callbackUrl` and `callbackBody` fields.
* `callback_var` - (Optional) Callback parameters. A Base64-encoded JSON document mapping custom callback variables (for example `{"x:a":"a","x:b":"b"}`).
* `policy_name` - (Optional) Name of the callback policy configuration item.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, which is the name of the owning bucket.

## Import

OSS Bucket Callback Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_callback_policy.example <bucket>
```
