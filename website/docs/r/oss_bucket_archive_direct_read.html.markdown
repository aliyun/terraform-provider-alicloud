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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_archive_direct_read&exampleId=2f53f111-5267-6bb6-3441-f494cea59b0f0c4b597a&activeTab=example&spm=docs.r.oss_bucket_archive_direct_read.0.2f53f11152&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oss_bucket_archive_direct_read&spm=docs.r.oss_bucket_archive_direct_read.example&intl_lang=EN_US)

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