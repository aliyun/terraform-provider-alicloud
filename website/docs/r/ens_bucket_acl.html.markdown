---
subcategory: "ENS (Edge Node Service)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_bucket_acl"
description: |-
  Provides a ENS Bucket Acl resource.
---

# alicloud_ens_bucket_acl

Provides a ENS Bucket Acl resource.

For information about ENS Bucket Acl and how to use it, see [PutBucketAcl](https://next.api.alibabacloud.com/document/ENS/2017-11-10/PutBucketAcl).

-> **NOTE:** Available since v 1.245.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ens_bucket_acl" "default" {
  bucket_name = "your-bucket-name"
  bucket_acl  = "public-read"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required, ForceNew) The name of the ENS bucket. The bucket name is the identity of the BucketAcl resource and cannot be changed after creation.
* `bucket_acl` - (Required) The read and write permissions of the bucket. Valid values: `public-read-write`, `public-read`, `private`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bucket Acl. The value is the bucket name.
* `bucket_name` - The name of the ENS bucket.
* `bucket_acl` - The read and write permissions of the bucket.

## Import

ENS Bucket Acl can be imported using the id, e.g.

```shell
terraform import alicloud_ens_bucket_acl.default <bucket_name>
```
