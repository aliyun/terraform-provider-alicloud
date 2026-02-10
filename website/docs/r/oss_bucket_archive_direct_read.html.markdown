---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_archive_direct_read"
description: |-
  Provides a Alicloud OSS Bucket Archive Direct Read resource.
---

# alicloud_oss_bucket_archive_direct_read

Provides a OSS Bucket Archive Direct Read resource.

Real-time access Archive objects in the bucket without the need to restore the Archive objects.

For information about OSS Bucket Archive Direct Read and how to use it, see [What is Bucket Archive Direct Read](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutBucketArchiveDirectRead).

-> **NOTE:** Available since v1.271.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
}


resource "alicloud_oss_bucket_archive_direct_read" "default" {
  bucket  = alicloud_oss_bucket.CreateBucket.id
  enabled = true
}
```

### Deleting `alicloud_oss_bucket_archive_direct_read` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_archive_direct_read`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket
* `enabled` - (Required) Specifies whether to enable real-time access of Archive objects for a bucket. Valid values: true and false.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Archive Direct Read.
* `update` - (Defaults to 5 mins) Used when update the Bucket Archive Direct Read.

## Import

OSS Bucket Archive Direct Read can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_archive_direct_read.example <bucket>
```