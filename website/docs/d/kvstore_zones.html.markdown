---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_zones"
sidebar_current: "docs-alicloud-datasource-kvstore-zones"
description: |-
    Provides a list of availability zones for KVStore that can be used by an Alibaba Cloud account.
---

# alicloud\_kvstore\_zones

This data source provides availability zones for KVStore that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in v1.73.0+.

## Example Usage

```
# Declare the data source
data "alicloud_kvstore_zones" "zones_ids" {}

# Create an KVStore instance with the first matched zone
resource "alicloud_kvstore_instance" "kvstore" {
    availability_zone = data.alicloud_kvstore_zones.zones_ids.zones.0.id

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch KVStore instances.
* `instance_charge_type` - (Optional) Filter the results by a specific instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `engine` - (Optional) Database type. Options are `Redis`, `Memcache`. Default to `Redis`.
* product_type - (Optional, Available in v1.130.0+) The type of the service. Valid values:
    * Local: an ApsaraDB for Redis instance with a local disk.
    * OnECS: an ApsaraDB for Redis instance with a standard disk. This type is available only on the Alibaba Cloud China site.
    
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.
