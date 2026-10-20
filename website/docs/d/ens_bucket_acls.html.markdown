---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_bucket_acls"
description: |-
  Provides the ENS Bucket Acls of the current Alibaba Cloud user.
---

# alicloud\_ens\_bucket\_acls

This data source provides the ENS Bucket Acl of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ens_bucket_acls" "default" {
  bucket_name = "your-bucket-name"
}

output "ens_bucket_acl_id" {
  value = data.alicloud_ens_bucket_acls.default.bucket_acls.0.id
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required) The name of the ENS bucket.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of ENS Bucket Acl IDs.
* `bucket_acls` - A list of ENS Bucket Acls. Each element contains the following attributes:
  * `bucket_acl` - The access control list (ACL) of the ENS bucket.
  * `bucket_name` - The name of the ENS bucket.
  * `id` - The ID of the ENS Bucket Acl, formatted as `<bucket_name>`.
