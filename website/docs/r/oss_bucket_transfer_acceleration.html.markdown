---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_transfer_acceleration"
description: |-
  Provides a Alicloud OSS Bucket Transfer Acceleration resource.
---

# alicloud_oss_bucket_transfer_acceleration

Provides a OSS Bucket Transfer Acceleration resource. Transfer acceleration configuration of a bucket.

For information about OSS Bucket Transfer Acceleration and how to use it, see [What is Bucket Transfer Acceleration](https://www.alibabacloud.com/help/en/oss/developer-reference/putbuckettransferacceleration).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_transfer_acceleration&exampleId=902fa53d-3c73-cdea-c3c6-acf529710ae9a8bd1b49&activeTab=example&spm=docs.r.oss_bucket_transfer_acceleration.0.902fa53d3c&intl_lang=EN_US" target="_blank">
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
  lifecycle {
    ignore_changes = [
      transfer_acceleration,
    ]
  }
}


resource "alicloud_oss_bucket_transfer_acceleration" "default" {
  bucket  = alicloud_oss_bucket.CreateBucket.bucket
  enabled = true
}
```

### Deleting `alicloud_oss_bucket_transfer_acceleration` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_transfer_acceleration`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oss_bucket_transfer_acceleration&spm=docs.r.oss_bucket_transfer_acceleration.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the Bucket.
* `enabled` - (Optional) Specifies whether to enable transfer acceleration for the bucket. Valid values: true: transfer acceleration for the bucket is enabled. false: transfer acceleration for the bucket is disabled.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Transfer Acceleration.
* `update` - (Defaults to 5 mins) Used when update the Bucket Transfer Acceleration.

## Import

OSS Bucket Transfer Acceleration can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_transfer_acceleration.example <id>
```