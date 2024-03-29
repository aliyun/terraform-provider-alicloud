---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_https_config"
description: |-
  Provides a Alicloud OSS Bucket Https Config resource.
---

# alicloud_oss_bucket_https_config

Provides a OSS Bucket Https Config resource. Whether the bucket can only be accessed with specific TLS versions.

For information about OSS Bucket Https Config and how to use it, see [What is Bucket Https Config](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.220.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
}


resource "alicloud_oss_bucket_https_config" "default" {
  tls_version = ["TLSv1.2"]
  bucket      = alicloud_oss_bucket.CreateBucket.bucket
  enable      = true
}
```

### Deleting `alicloud_oss_bucket_https_config` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_https_config`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.
* `enable` - (Required) Specifies whether to enable TLS version management for the bucket. Valid values: true, false.
* `tls_version` - (Optional) Specifies the TLS versions allowed to access this buckets.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Https Config.
* `update` - (Defaults to 5 mins) Used when update the Bucket Https Config.

## Import

OSS Bucket Https Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_https_config.example <id>
```