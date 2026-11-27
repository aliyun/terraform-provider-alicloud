---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_zones"
sidebar_current: "docs-alicloud-datasource-gpdb-zones"
description: |-
    Provides a list of availability zones for Gpdb that can be used by an Alibaba Cloud account.
---

# alicloud_gpdb_zones

This data source provides availability zones for Gpdb that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available since v1.73.0+.

## Example Usage

```terraform
# Declare the data source
data "alicloud_gpdb_zones" "zones_ids" {}

# Create an Gpdb instance with the first matched zone
resource "alicloud_hbase_instance" "hbase" {
  availability_zone = data.alicloud_gpdb_zones.zones_ids.zones[0].id

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch Gpdb instances.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones.

---

### `zones`

The zones supports the following attributes:

* `id` - ID of the zone.
* `multi_zone_ids` - A list of zone ids in which the multi zone.
