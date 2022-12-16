---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_cross_regions"
sidebar_current: "docs-alicloud-datasource-rds-cross-regions"
description: |-
  Provide a list of available RDS remote disaster recovery regions for Alibaba Cloud accounts.
---

# alicloud\_rds\_cross\_regions

This data source provides an available area for remote disaster recovery for RDS.

-> **NOTE:** Available in v1.193.0+.

## Example Usage

```
# Declare the data source
data "alicloud_rds_cross_regions" "cross_regions" {}

output "first_rds_cross_regions" {
  value = data.alicloud_rds_cross_regions.regions.ids.0
}
```

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of region IDs.
* `regions` - The list of destination regions that support cross-region backup. Each element contains the following attributes:
  * `id` - ID of the region.
