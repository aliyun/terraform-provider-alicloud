---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_buckets"
sidebar_current: "docs-alicloud-datasource-ens-buckets"
description: |-
  Provides a list of ENS Buckets to the user.
---

# alicloud_ens_buckets

This data source provides the ENS Buckets of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ens_buckets" "example" {
  name_regex = "^terraform-example"
}
output "ens_bucket_id_1" {
  value = data.alicloud_ens_buckets.example.buckets.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by bucket name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `prefix` - (Optional) The prefix that the returned bucket name must start with. If this parameter is not set, the prefix information is not filtered.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Bucket IDs.
* `names` - A list of Bucket names.
* `buckets` - A list of ENS Buckets. Each element contains the following attributes:
  * `id` - The ID of the Bucket. The value is the bucket name.
  * `bucket_name` - The name of the bucket.
  * `comment` - The description of the bucket.
  * `bucket_acl` - The read and write permission type of the bucket.
  * `logical_bucket_type` - The logical bucket type.
  * `ens_region_id` - The ENS node ID.
  * `create_time` - The creation time of the bucket.
  * `modify_time` - The modification time of the bucket.
