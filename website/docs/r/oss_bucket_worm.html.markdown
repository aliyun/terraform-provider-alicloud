---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_worm"
description: |-
  Provides a Alicloud OSS Bucket Worm resource.
---

# alicloud_oss_bucket_worm

Provides a OSS Bucket Worm resource.

Bucket Retention Policy.

For information about OSS Bucket Worm and how to use it, see [What is Bucket Worm](https://www.alibabacloud.com/help/en/oss/developer-reference/initiatebucketworm).

-> **NOTE:** Available since v1.240.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_worm&exampleId=cdc3e025-aec0-bea3-a4ea-4b8b456dfa98f1a43289&activeTab=example&spm=docs.r.oss_bucket_worm.0.cdc3e025ae&intl_lang=EN_US" target="_blank">
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

resource "alicloud_oss_bucket" "defaulthNMfIF" {
  storage_class = "Standard"
}


resource "alicloud_oss_bucket_worm" "default" {
  bucket                   = alicloud_oss_bucket.defaulthNMfIF.bucket
  retention_period_in_days = "1"
  status                   = "InProgress"
}
```

### Deleting `alicloud_oss_bucket_worm` or removing it from your configuration

The `alicloud_oss_bucket_worm` resource allows you to manage  `status = "Locked"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket
* `retention_period_in_days` - (Optional, Int) The specified number of days to retain the Object.
* `status` - (Optional) The status of the compliance retention policy. Optional values:
  - `InProgress`: After a compliance retention policy is created, the policy is in the InProgress status by default, and the validity period of this status is 24 hours.
  - `Locked`: The compliance retention policy is Locked.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<bucket>:<worm_id>`.
* `create_time` - The creation time of the resource
* `worm_id` - The ID of the retention policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Worm.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Worm.
* `update` - (Defaults to 5 mins) Used when update the Bucket Worm.

## Import

OSS Bucket Worm can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_worm.example <bucket>:<worm_id>
```