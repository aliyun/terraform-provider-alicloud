---
subcategory: "Kvcachestore"
layout: "alicloud"
page_title: "Alicloud: alicloud_kv_cache_stores"
description: |-
  Provides a list of Kvcachestore KVCacheStores to the user.
---

# alicloud_kv_cache_stores

This data source provides the KVCacheStores of the current AlicaCloud user.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_kv_cache_stores" "default" {
  ids = ["<kvcs_id>"]
}

output "first_store_kvcs_id" {
  value = data.alicloud_kv_cache_stores.default.stores.0.kvcs_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of KVCacheStore IDs.
* `kvcs_ids` - (Optional, ForceNew) The IDs of the KVCacheStores. Multiple IDs are separated by commas (,).
* `name` - (Optional, ForceNew) The name of the KVCacheStore.
* `status` - (Optional, ForceNew) The status of the KVCacheStore.
* `zone_id` - (Optional, ForceNew) The zone ID of the KVCacheStore.
* `output_file` - (Optional) File path where data source results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `stores` - A list of KVCacheStores. Each element contains the following attributes:
  * `capacity` - Storage capacity in GiB. Minimum 307200 GiB (300 TiB), expanding in steps of 300 TiB (307200 GiB).
  * `create_time` - The creation time of the resource.
  * `description` - The description of the KVCacheStore instance.
  * `hpn_zone` - HPN cluster code used for affinity scheduling. This is a cluster identifier (e.g. `default`, `B6`), not a zone ID.
  * `kvcs_id` - The primary identifier of the KVCacheStore resource.
  * `name` - The name of the KVCacheStore instance.
  * `payment_type` - The payment type of the resource. Valid values: `POSTPAY` (pay-as-you-go), `PREPAY` (subscription).
  * `region_id` - The region ID of the resource.
  * `resource_group_id` - The ID of the resource group.
  * `status` - The status of the resource.
  * `tags` - A mapping of tags assigned to the resource.
  * `zone_id` - The zone ID of the KVCacheStore instance.
