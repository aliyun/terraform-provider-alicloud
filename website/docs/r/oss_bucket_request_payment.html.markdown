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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_request_payment&exampleId=9f16db7f-6cbf-083d-612e-45a23b3bfac3fc468f87&activeTab=example&spm=docs.r.oss_bucket_request_payment.0.9f16db7f6c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Request Payment.
* `update` - (Defaults to 5 mins) Used when update the Bucket Request Payment.

## Import

OSS Bucket Request Payment can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_request_payment.example <id>
```