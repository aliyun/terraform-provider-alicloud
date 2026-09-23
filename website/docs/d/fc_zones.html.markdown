---
subcategory: "Function Compute Service (FC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_zones"
sidebar_current: "docs-alicloud-datasource-fc-zones"
description: |-
    Provides a list of availability zones for FunctionCompute that can be used by an Alibaba Cloud account.
---

# alicloud\_fc\_zones

This data source provides availability zones for FunctionCompute that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in v1.74.0+.

## Example Usage

```
# Declare the data source
data "alicloud_fc_zones" "zones_ids" {}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
