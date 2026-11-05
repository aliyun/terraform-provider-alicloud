---
subcategory: "Kvcachestore"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvcachestore_kv_cache_stores"
sidebar_current: "docs-alicloud-datasource-kvcachestore-kv-cache-stores"
description: |-
  Provides a list of Kvcachestore Kv Cache Store owned by an Alibaba Cloud account.
---

# alicloud_kvcachestore_kv_cache_stores

This data source provides Kvcachestore Kv Cache Store available to the user. [What is Kv Cache Store](https://next.api.alibabacloud.com/document/Kvcachestore/2026-06-17/CreateKVCacheStore)

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_kvcachestore_kv_cache_stores" "default" {
  name = "terraform-example"
}

output "first_store_id" {
  value = data.alicloud_kvcachestore_kv_cache_stores.default.stores.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, Computed) A list of Kv Cache Store IDs.
* `kvcs_ids` - (Optional) The ID of the Kv Cache Store used to filter the results.
* `name` - (Optional) The instance name used to filter the results.
* `status` - (Optional) The instance status used to filter the results. Valid values: `Creating`, `Available`, `InUse`, `Stopping`, `Stopped`, `Deleting`.
* `zone_id` - (Optional) The zone ID used to filter the results.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `stores` - A list of Kv Cache Store entries. Each element contains the following attributes:
  * `id` - The ID of the instance. The value is the same as `kvcs_id`.
  * `kvcs_id` - The ID of the instance.
  * `capacity` - The storage capacity of the instance. Unit: GiB.
  * `create_time` - The creation time of the instance (ISO 8601).
  * `description` - The description of the instance.
  * `hpn_zone` - The cluster ID (HPN zone) of the instance.
  * `name` - The name of the instance.
  * `payment_type` - The payment type of the instance.
  * `region_id` - The region ID of the instance.
  * `status` - The status of the instance. **NOTE:** This field is only available when `enable_details` is `true`.
  * `zone_id` - The zone ID of the instance.
