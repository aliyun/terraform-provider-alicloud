---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_express_sync"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-express-sync"
description: |-
  Provides a Alicloud Cloud Storage Gateway Express Sync resource.
---

# alicloud\_cloud\_storage\_gateway\_express\_sync

Provides a Cloud Storage Gateway Express Sync resource.

For information about Cloud Storage Gateway Express Sync and how to use it, see [What is Express Sync](https://www.alibabacloud.com/help/en/doc-detail/53972.htm).

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tftest"
}

variable "region" {
  default = "cn-shanghai"
}

data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = true
  public_network_bandwidth = 10
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
  gateway_name             = var.name
}

resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
  acl    = "public-read-write"
}

resource "alicloud_cloud_storage_gateway_gateway_file_share" "default" {
  gateway_file_share_name = var.name
  gateway_id              = alicloud_cloud_storage_gateway_gateway.default.id
  local_path              = alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path
  oss_bucket_name         = alicloud_oss_bucket.default.bucket
  oss_endpoint            = alicloud_oss_bucket.default.extranet_endpoint
  protocol                = "NFS"
  remote_sync             = true
  polling_interval        = 4500
  fe_limit                = 0
  backend_limit           = 0
  cache_mode              = "Cache"
  squash                  = "none"
  lag_period              = 5
}

resource "alicloud_cloud_storage_gateway_express_sync" "default" {
  bucket_name       = alicloud_cloud_storage_gateway_gateway_file_share.default.oss_bucket_name
  bucket_region     = var.region
  description       = var.name
  express_sync_name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required, ForceNew) The name of the OSS Bucket.
* `bucket_prefix` - (Optional, ForceNew) The prefix of the OSS Bucket.
* `bucket_region` - (Required, ForceNew) The region of the OSS Bucket.
* `description` - (Optional, ForceNew) The description of the Express Sync. The length of the name is limited to `1` to `255` characters.
* `express_sync_name` - (Required, ForceNew) The name of the ExpressSync. The length of the name is limited to `1` to `128` characters. It can contain uppercase and lowercase letters, Chinese characters, numbers, English periods (.), underscores (_), or hyphens (-), and must start with  letters.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of ExpressSync. The value is formate as <express_sync_id>.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `delete` - (Defaults to 5 mins) Used when delete the Express Sync.

## Import

Cloud Storage Gateway Express Sync can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_storage_gateway_express_sync.example <id>
```