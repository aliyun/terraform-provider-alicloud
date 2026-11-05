---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_bucket"
sidebar_current: "docs-alicloud-resource-ens-bucket"
description: |-
  Provides a ENS Bucket resource.
---

# alicloud\_ens\_bucket

Provides a ENS Bucket resource.

For information about ENS Bucket and how to use it, see [What is Bucket](https://www.alibabacloud.com/help/en/ens/developer-reference/).

-> **NOTE:** Available since v2.3.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_bucket" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

resource "alicloud_ens_bucket" "default" {
  bucket_name = "tf-example-bucket"
  comment     = "managed by terraform"
  bucket_acl  = "private"
}
```

Single-node (sink) bucket

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

resource "alicloud_ens_bucket" "sink" {
  bucket_name         = "tf-example-sink-bucket"
  logical_bucket_type = "sink"
  ens_region_id       = "cn-beijing-1"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required, ForceNew) The name of the bucket. 3 to 50 characters in length. Valid values can only contain lowercase letters, digits and hyphens (-), and must not start or end with a hyphen.
* `logical_bucket_type` - (Optional, ForceNew) The logical bucket type. Valid values: `sink`, `standard`. Defaults to `standard`. `sink` is a single-node bucket bound to one edge node and requires `ens_region_id`; `standard` is a global bucket dispatched across nodes.
* `comment` - (Optional, ForceNew) The remark of the bucket. 0 to 128 characters in length. Changing this field recreates the resource, because the ENS API exposes no in-place update for the bucket comment.
* `bucket_acl` - (Optional) The access control list of the bucket. Valid values: `private`, `public-read`, `public-read-write`. Defaults to `private`. Updatable in place via PutBucketAcl.
* `ens_region_id` - (Optional, ForceNew) The node region id. Required when `logical_bucket_type` is `sink`. Write-only: not returned by GetBucketInfo, so the configured value is preserved in state and cannot be read back or imported.
* `dispatch_scope` - (Optional, ForceNew) The dispatch scope. Valid values: `domestic`, `oversea`. Optional for a global (`standard`) bucket; specifies the dispatch scope when explicit control is needed. Write-only: not returned by GetBucketInfo, so the configured value is preserved in state and cannot be read back or imported.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. It is the same as `bucket_name`.
* `create_time` - The time when the bucket was created.
* `modify_time` - The time when the bucket was last modified. This value is updated whenever the ACL is changed.

## Import

ENS bucket can be imported using the id, e.g.

```shell
terraform import alicloud_ens_bucket.default tf-example-bucket
```
