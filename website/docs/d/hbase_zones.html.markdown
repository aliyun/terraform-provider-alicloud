---
subcategory: "HBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbase_zones"
sidebar_current: "docs-alicloud-datasource-hbase-zones"
description: |-
    Provides a list of availability zones for HBase that can be used by an Alibaba Cloud account.
---

# alicloud\_hbase\_zones

This data source provides availability zones for HBase that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in v1.73.0+.

## Example Usage

```terraform
data "alicloud_hbase_zones" "zones_ids" {}

resource "alicloud_hbase_instance" "hbase" {
  zone_id = data.alicloud_hbase_zones.zones_ids.zones[0].id

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Removed from v1.99.0) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch HBase instances. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone. Removed from v1.99.0.
