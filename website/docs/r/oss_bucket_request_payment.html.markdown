---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_request_payment"
description: |-
  Provides a Alicloud OSS Bucket Request Payment resource.
---

# alicloud_oss_bucket_request_payment

Provides a OSS Bucket Request Payment resource. Whether to enable pay-by-requester for a bucket.

For information about OSS Bucket Request Payment and how to use it, see [What is Bucket Request Payment](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketrequestpayment).

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
}


resource "alicloud_oss_bucket_request_payment" "default" {
  payer  = "Requester"
  bucket = alicloud_oss_bucket.CreateBucket.bucket
}
```

### Deleting `alicloud_oss_bucket_request_payment` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_request_payment`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.
* `payer` - (Optional, Computed) The payer of the request and traffic fees.Valid values: BucketOwner: request and traffic fees are paid by the bucket owner. Requester: request and traffic fees are paid by the requester.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Request Payment.
* `update` - (Defaults to 5 mins) Used when update the Bucket Request Payment.

## Import

OSS Bucket Request Payment can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_request_payment.example <id>
```