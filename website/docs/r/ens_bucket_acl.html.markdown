---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_bucket_acl"
description: |-
  Provides a Alicloud ENS Bucket Acl resource.
---

# alicloud\_ens\_bucket\_acl

Provides a ENS Bucket Acl resource.

For information about ENS Bucket Acl and how to use it, see [What is BucketAcl](https://www.alibabacloud.com/help/en/ens/developer-reference/).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_bucket_acl" "default" {
  bucket_name = "your-bucket-name"
  bucket_acl  = "private"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_acl` - (Required) The access control list (ACL) of the ENS bucket. Valid values: `public-read-write`, `public-read`, `private`.
* `bucket_name` - (Required, ForceNew) The name of the ENS bucket.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. It is formatted as `<bucket_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Bucket Acl.
* `update` - (Defaults to 5 mins) Used when update the Bucket Acl.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Acl.

## Import

ENS Bucket Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_bucket_acl.example <bucket_name>
```
