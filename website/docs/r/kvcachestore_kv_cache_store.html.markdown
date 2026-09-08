---
subcategory: "Kvcachestore"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvcachestore_kv_cache_store"
description: |-
  Provides a Alicloud Kvcachestore Kv Cache Store resource.
---

# alicloud_kvcachestore_kv_cache_store

Provides a Kvcachestore Kv Cache Store resource.

For information about Kvcachestore Kv Cache Store and how to use it, see [What is Kv Cache Store](https://next.api.alibabacloud.com/document/Kvcachestore/2026-06-17/CreateKVCacheStore).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_kvcachestore_kv_cache_store" "default" {
  name        = var.name
  description = "example kvcachestore instance"
  zone_id     = "cn-shanghai-cloudspe-b"
  hpn_zone    = "b1"
  capacity    = 307200
}
```

## Argument Reference

The following arguments are supported:
* `capacity` - (Required, Int) The storage capacity of the instance. Unit: GiB. The minimum value is `307200` (300 TiB), and the value must be a multiple of `307200`.
* `description` - (Optional) The description of the instance.
* `hpn_zone` - (Required, ForceNew) The cluster ID (HPN zone) of the instance.
* `name` - (Optional) The name of the instance.
* `payment_type` - (Optional, ForceNew, Computed) The payment type of the instance. Valid values: `PREPAY`, `POSTPAY`.
* `zone_id` - (Required, ForceNew) The zone ID of the instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is the KvcsId of the instance.
* `create_time` - The creation time of the instance (ISO 8601).
* `status` - The status of the instance. Valid values: `Creating`, `Available`, `InUse`, `Stopping`, `Stopped`, `Deleting`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Kv Cache Store.
* `delete` - (Defaults to 5 mins) Used when delete the Kv Cache Store.
* `update` - (Defaults to 5 mins) Used when update the Kv Cache Store.

## Import

Kvcachestore Kv Cache Store can be imported using the id, e.g.

```shell
$ terraform import alicloud_kvcachestore_kv_cache_store.example <kvcs_id>
```
