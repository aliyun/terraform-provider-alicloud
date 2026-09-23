---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_zones"
sidebar_current: "docs-alicloud-datasource-polardb-zones"
description: |-
    Provides a list of availability zones for PolarDB that can be used by an Alibaba Cloud account.
---

# alicloud\_polardb\_zones

This data source provides availability zones for PolarDB that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in v1.74.0+.

## Example Usage

```terraform
# Declare the data source
data "alicloud_polardb_zones" "zones_ids" {}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch PolarDB instances.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.
