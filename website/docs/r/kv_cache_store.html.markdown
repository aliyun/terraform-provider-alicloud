---
subcategory: "Kvcachestore"
layout: "alicloud"
page_title: "Alicloud: alicloud_kv_cache_store"
description: |-
  Provides a Alicloud Kvcachestore KVCacheStore resource.
---

# alicloud_kv_cache_store

Provides a Kvcachestore KVCacheStore resource. KVCacheStore is a high-performance key-value storage service on Alibaba Cloud.

For information about Kvcachestore KVCacheStore and how to use it, see [What is KVCacheStore](https://next.api.alibabacloud.com/document/Kvcachestore/2026-06-17/CreateKVCacheStore).

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_kv_cache_store" "default" {
  capacity     = 307200
  zone_id      = "cn-hangzhou-b"
  hpn_zone     = "default"
  name         = var.name
  description  = "terraform-example-description"
  payment_type = "POSTPAY"
  tags = {
    Created = "Terraform"
    Env     = "dev"
  }
}
```

## Argument Reference

The following arguments are supported:

* `capacity` - (Required, Int) Storage capacity in GiB. Minimum 307200 GiB (300 TiB), expanding in steps of 300 TiB (307200 GiB).
* `zone_id` - (Required, ForceNew) Zone ID. Query available zones for a region via the DescribeZones API. Changing it will create a new resource.
* `hpn_zone` - (Required, ForceNew) HPN cluster code used for affinity scheduling between the KVCacheStore and a specified HPN cluster. This is a cluster identifier (e.g. `default`, `B6`), not a zone ID. Changing it will create a new resource.
* `name` - (Optional) The name of the KVCacheStore instance.
* `description` - (Optional) The description of the KVCacheStore instance.
* `payment_type` - (Optional, ForceNew, Computed) The payment type of the resource. Valid values: `POSTPAY` (pay-as-you-go), `PREPAY` (subscription). Create currently only supports `POSTPAY`; defaults to `POSTPAY` when omitted.
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the KVCacheStore belongs. Updating this will trigger a ChangeResourceGroup operation.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the KVCacheStore. It is the same as `kvcs_id`.
* `kvcs_id` - The primary identifier of the KVCacheStore resource.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.
* `status` - The status of the resource. Valid values: `Creating`, `Available`, `InUse`, `Stopping`, `Stopped`, `Deleting`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the KVCacheStore.
* `update` - (Defaults to 20 mins) Used when update the KVCacheStore.
* `delete` - (Defaults to 10 mins) Used when delete the KVCacheStore.

## Import

KVCacheStore can be imported using the id, e.g.

```shell
$ terraform import alicloud_kv_cache_store.example <kvcs_id>
```
