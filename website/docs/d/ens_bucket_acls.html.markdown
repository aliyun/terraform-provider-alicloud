---
subcategory: "ENS (Edge Node Service)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_bucket_acls"
description: |-
  Provides a list of ENS Bucket Acls to the user.
---

# alicloud_ens_bucket_acls

This data source provides the ENS Bucket Acls of the current Alibaba Cloud user.

-> **NOTE:** Available since v 1.245.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ens_bucket_acls" "default" {
  bucket_name = "your-bucket-name"
}

output "bucket_acl" {
  value = data.alicloud_ens_bucket_acls.default.bucket_acls.0.bucket_acl
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required) The name of the ENS bucket to query the ACL for.

## Attributes Reference

The following attributes are exported in addition to the `bucket_name` argument above:

* `ids` - A list of Bucket Acl IDs. The ID is the bucket name.
* `bucket_acls` - A list of ENS Bucket Acls. Each element contains the following attributes:
  * `id` - The ID of the Bucket Acl. The value is the bucket name.
  * `bucket_name` - The name of the ENS bucket.
  * `bucket_acl` - The read and write permissions of the bucket. Valid values: `public-read-write`, `public-read`, `private`.
