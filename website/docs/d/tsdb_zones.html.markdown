---
subcategory: "Time Series Database (TSDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_tsdb_zones"
sidebar_current: "docs-alicloud-datasource-tsdb-zones"
description: |-
  Provides a list of Time Series Database (TSDB) instance available zones to the user.
---

# alicloud\_tsdb\_zones

This data source provides the available zones with the Time Series Database (TSDB) Instance of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.112.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_tsdb_zones" "example" {}

output "first_tsdb_zones_id" {
  value = data.alicloud_tsdb_zones.example.zones.0.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of TSDB instance zone IDs.
* `zones` - A list of TSDB Instance zones. Each element contains the following attributes:
	* `id` - The ID of zone.
	* `zone_id` - The zone ID.
	* `local_name` - The local name.
