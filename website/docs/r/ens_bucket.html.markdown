---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_bucket"
sidebar_current: "docs-alicloud-resource-ens-bucket"
description: |-
  Provides a Alicloud ENS Bucket resource.
---

# alicloud_ens_bucket

Provides a ENS Bucket resource.

ENS Bucket (Edge Object Storage) is an object storage service provided by Edge Node Service (ENS). It is used to aggregate and manage files at edge nodes.

For information about ENS Bucket and how to use it, see [What is Bucket](https://www.alibabacloud.com/help/en/ens/latest/putbucket).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ens_bucket" "example" {
  bucket_name         = "terraform-example"
  comment             = "terraform example bucket"
  bucket_acl          = "private"
  logical_bucket_type = "sink"
  ens_region_id       = "cn-fuyou-1"
  dispatch_scope      = "domestic"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required, ForceNew) The name of the bucket. The logical bucket name is used to aggregate and manage files.
* `comment` - (Optional) The description of the bucket.
* `bucket_acl` - (Optional, ForceNew) The read and write permission type of the bucket. Valid values: `private`, `public-read`, `public-read-write`. Default value: `private`.
* `logical_bucket_type` - (Optional, ForceNew) The logical bucket type. Valid values: `sink` (single-node storage), `standard` (standard storage).
* `ens_region_id` - (Optional, ForceNew) The ENS node ID. Required when `logical_bucket_type` is `sink`.
* `dispatch_scope` - (Optional, ForceNew) The scheduling scope. This parameter is valid only for global scheduling buckets. Valid values: `domestic` (Mainland China), `oversea` (outside Mainland China).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Bucket. The value is the bucket name.
* `create_time` - The creation time of the bucket. The date format is in accordance with ISO8601 notation and uses UTC time. The format is yyyy-MM-ddTHH:mm:ssZ.
* `modify_time` - The modification time of the bucket. The date format is in accordance with ISO8601 notation and uses UTC time. The format is yyyy-MM-ddTHH:mm:ssZ.

## Import

ENS Bucket can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_bucket.example <bucket_name>
```
