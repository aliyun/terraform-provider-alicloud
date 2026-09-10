---
subcategory: "Tair (Redis OSS-Compatible) And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_zones"
sidebar_current: "docs-alicloud-datasource-kvstore-zones"
description: |-
    Provides a list of availability zones for Tair (Redis OSS-Compatible) And Memcache (KVStore) that can be used by an Alibaba Cloud account.
---

# alicloud_kvstore_zones

This data source provides availability zones for Tair (Redis OSS-Compatible) And Memcache (KVStore) that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available since v1.73.0.

## Example Usage

```terraform
# Declare the data source
data "alicloud_kvstore_zones" "zones_ids" {
  instance_charge_type = "PostPaid"
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch Tair (Redis OSS-Compatible) And Memcache (KVStore) instances.
* `instance_charge_type` - (Optional) Filter the results by a specific instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `engine` - (Optional) Database type. Options are `Redis`, `Memcache`. Default to `Redis`.
* product_type - (Optional, Available since v1.130.0+) The type of the service. Valid values: `Local`, `Tair_rdb`, `Tair_scm`, `Tair_essd`, `OnECS`.
    
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.
