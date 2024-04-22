---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_versioning"
description: |-
  Provides a Alicloud OSS Bucket Versioning resource.
---

# alicloud_oss_bucket_versioning

Provides a OSS Bucket Versioning resource. Configures the versioning state for a bucket.

For information about OSS Bucket Versioning and how to use it, see [What is Bucket Versioning](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketversioning).

-> **NOTE:** Available since v1.222.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-${random_integer.default.result}"
  lifecycle {
    ignore_changes = [
      versioning,
    ]
  }
}


resource "alicloud_oss_bucket_versioning" "default" {
  status = "Enabled"
  bucket = alicloud_oss_bucket.CreateBucket.bucket
}
```

### Deleting `alicloud_oss_bucket_versioning` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_versioning`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.
* `status` - (Optional, Computed) A bucket can be in one of the following versioning states: disabled, enabled, or suspended. By default, versioning is disabled for a bucket. Updating the value from Enabled or Suspended to Disabled will result in errors, because OSS does not support returning buckets to an unversioned state. .

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Versioning.
* `update` - (Defaults to 5 mins) Used when update the Bucket Versioning.

## Import

OSS Bucket Versioning can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_versioning.example <id>
```