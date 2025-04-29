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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_access_monitor&exampleId=82af3813-21a3-964c-e231-a8a40cbfcdc3c64f23aa&activeTab=example&spm=docs.r.oss_bucket_access_monitor.0.82af381321&intl_lang=EN_US" target="_blank">
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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Access Monitor.
* `update` - (Defaults to 5 mins) Used when update the Bucket Access Monitor.

## Import

OSS Bucket Access Monitor can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_access_monitor.example <id>
```