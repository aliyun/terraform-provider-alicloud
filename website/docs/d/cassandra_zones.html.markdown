---
subcategory: "Cassandra"
layout: "alicloud"
page_title: "Alicloud: alicloud_cassandra_zones"
sidebar_current: "docs-alicloud-datasource-cassandra-zones"
description: |-
    Provides a list of availability zones for Cassandra that can be used by an Alibaba Cloud account.
---

# alicloud\_cassandra\_zones

This data source provides availability zones for Cassandra that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in v1.88.0+.

## Example Usage

```
# Declare the data source
data "alicloud_cassandra_zones" "zones_ids" {}

# Create an Cassandra cluster with the first matched zone
resource "alicloud_cassandra_cluster" "cassandra" {
    zone_id = data.alicloud_cassandra_zones.zones_ids.zones[0].id

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch Cassandra clusters.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.
