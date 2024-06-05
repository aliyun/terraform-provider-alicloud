---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_data_redundancy_transition"
description: |-
  Provides a Alicloud OSS Bucket Data Redundancy Transition resource.
---

# alicloud_oss_bucket_data_redundancy_transition

Provides a OSS Bucket Data Redundancy Transition resource. Create a storage redundancy transition task to convert local redundant storage(LRS) to zone redundant storage(ZRS).

For information about OSS Bucket Data Redundancy Transition and how to use it, see [What is Bucket Data Redundancy Transition](https://www.alibabacloud.com/help/en/oss/developer-reference/createbucketdataredundancytransition).

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


resource "alicloud_oss_bucket_data_redundancy_transition" "default" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) Storage space name.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<bucket>:<task_id>`.
* `create_time` - Stores the creation time of the redundant transformation task.
* `status` - Stores the state of the redundant translation task. The values are as follows:  Queueing: in the queue.  Processing: In progress.  Finished: Finished.
* `task_id` - Unique identification of the storage redundancy conversion task.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Data Redundancy Transition.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Data Redundancy Transition.

## Import

OSS Bucket Data Redundancy Transition can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_data_redundancy_transition.example <bucket>:<task_id>
```