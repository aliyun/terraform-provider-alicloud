---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_access_monitor"
description: |-
  Provides a Alicloud OSS Bucket Access Monitor resource.
---

# alicloud_oss_bucket_access_monitor

Provides a OSS Bucket Access Monitor resource. Enables or disables access tracking for a bucket.

For information about OSS Bucket Access Monitor and how to use it, see [What is Bucket Access Monitor](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketaccessmonitor).

-> **NOTE:** Available since v1.224.0.

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
}


resource "alicloud_oss_bucket_access_monitor" "default" {
  status = "Enabled"
  bucket = alicloud_oss_bucket.CreateBucket.bucket
}
```

### Deleting `alicloud_oss_bucket_access_monitor` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_access_monitor`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.
* `status` - (Required) Specifies whether to enable access tracking for the bucket. Valid values: Enabled: enables access tracking. Disabled: disables access tracking.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Access Monitor.
* `update` - (Defaults to 5 mins) Used when update the Bucket Access Monitor.

## Import

OSS Bucket Access Monitor can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_access_monitor.example <id>
```