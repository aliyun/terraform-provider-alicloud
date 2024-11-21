---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_express_sync_share_attachment"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-express-sync-share-attachment"
description: |-
  Provides a Alicloud Cloud Storage Gateway Express Sync Share Attachment resource.
---

# alicloud_cloud_storage_gateway_express_sync_share_attachment

Provides a Cloud Storage Gateway Express Sync Share Attachment resource.

For information about Cloud Storage Gateway Express Sync Share Attachment and how to use it, see [What is Express Sync Share Attachment](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/addsharestoexpresssync).

-> **NOTE:** Available since v1.144.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_storage_gateway_express_sync_share_attachment&exampleId=79463bda-8c0a-69f2-1308-f17e96246676e6cc7252&activeTab=example&spm=docs.r.cloud_storage_gateway_express_sync_share_attachment.0.79463bda8c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_regions" "default" {
  current = true
}
resource "random_uuid" "default" {
}
resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}

resource "alicloud_oss_bucket" "default" {
  bucket = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}

resource "alicloud_oss_bucket_acl" "default" {
  bucket = alicloud_oss_bucket.default.bucket
  acl    = "public-read-write"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}
data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  gateway_name             = var.name
  description              = var.name
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = true
  public_network_bandwidth = 40
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
}

resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
}

resource "alicloud_cloud_storage_gateway_gateway_file_share" "default" {
  gateway_file_share_name = var.name
  gateway_id              = alicloud_cloud_storage_gateway_gateway.default.id
  local_path              = alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path
  oss_bucket_name         = alicloud_oss_bucket.default.bucket
  oss_endpoint            = alicloud_oss_bucket.default.extranet_endpoint
  protocol                = "NFS"
  remote_sync             = false
  fe_limit                = 0
  backend_limit           = 0
  cache_mode              = "Cache"
  squash                  = "none"
  lag_period              = 5
}

resource "alicloud_cloud_storage_gateway_express_sync" "default" {
  bucket_name       = alicloud_cloud_storage_gateway_gateway_file_share.default.oss_bucket_name
  bucket_region     = data.alicloud_regions.default.regions.0.id
  description       = var.name
  express_sync_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cloud_storage_gateway_express_sync_share_attachment" "default" {
  express_sync_id = alicloud_cloud_storage_gateway_express_sync.default.id
  gateway_id      = alicloud_cloud_storage_gateway_gateway.default.id
  share_name      = alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_file_share_name
}
```

## Argument Reference

The following arguments are supported:

* `express_sync_id` - (Required, ForceNew) The ID of the ExpressSync.
* `gateway_id` - (Required, ForceNew) The ID of the Gateway.
* `share_name` - (Required, ForceNew) The name of the GatewayFileShare. **NOTE:** When GatewayFileShare is associated with a speed sync group, its reverse synchronization function will be turned off by default.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Express Sync Share Attachment. The value is formate as <express_sync_id>:<gateway_id>:<share_name>.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Express Sync Share Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Express Sync Share Attachment.

## Import

Cloud Storage Gateway Express Sync Share Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_storage_gateway_express_sync_share_attachment.example <express_sync_id>:<gateway_id>:<share_name>
```