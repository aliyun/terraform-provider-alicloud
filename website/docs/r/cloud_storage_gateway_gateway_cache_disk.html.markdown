---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_cache_disk"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-cache-disk"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway Cache Disk resource.
---

# alicloud_cloud_storage_gateway_gateway_cache_disk

Provides a Cloud Storage Gateway Gateway Cache Disk resource.

For information about Cloud Storage Gateway Gateway Cache Disk and how to use it, see [What is Gateway Cache Disk](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/creategatewaycachedisk).

-> **NOTE:** Available since v1.144.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_storage_gateway_gateway_cache_disk&exampleId=309571b6-ce43-f39a-c5e0-da625a6af796bf728756&activeTab=example&spm=docs.r.cloud_storage_gateway_gateway_cache_disk.0.309571b6ce&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_uuid" "default" {
}
resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
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
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_storage_gateway_gateway_cache_disk&spm=docs.r.cloud_storage_gateway_gateway_cache_disk.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The ID of the gateway.
* `cache_disk_size_in_gb` - (Required, Int) The capacity of the cache disk.
* `cache_disk_category` - (Optional, ForceNew) The type of the cache disk. Valid values: `cloud_efficiency`, `cloud_ssd`, `cloud_essd`. **NOTE:** From version 1.227.0, `cache_disk_category` can be set to `cloud_essd`.
* `performance_level` - (Optional, ForceNew, Available since v1.227.0) The performance level (PL) of the Enterprise SSD (ESSD). Valid values: `PL1`, `PL2`, `PL3`. **NOTE:** If `cache_disk_category` is set to `cloud_essd`, `performance_level` is required.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gateway Cache Disk. It formats as `<gateway_id>:<cache_id>:<local_file_path>`.
* `cache_id` - The ID of the cache disk.
* `local_file_path` - The path of the cache disk.
* `status` - The status of the Gateway Cache Disk.

## Import

Cloud Storage Gateway Gateway Cache Disk can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway_cache_disk.example <gateway_id>:<cache_id>:<local_file_path>
```
