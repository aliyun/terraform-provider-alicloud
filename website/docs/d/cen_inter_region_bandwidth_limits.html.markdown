---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_inter_region_bandwidth_limits"
sidebar_current: "docs-alicloud-datasource-cen-inter-region-bandwidth-limits"
description: |-
    Provides a list of CEN Inter Region Bandwidth Limits of the CEN.
---

# alicloud\_cen\_inter\_region\_bandwidth\_limits

The CEN inter region bandwidth limits data source lists a number of CEN inter region bandwidth limit resource information of the CEN.

## Example Usage

```
data "alicloud_cen_inter_region_bandwidth_limits" "multi_cen_bwl"{
	cen_id = "cen-id1"
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Optional) Limit search to a CEN ID, like "cen-id1".
* `output_file` - (Optional) The name of file that can save CEN Bandwidth Package data source after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `cen_id` - ID of the CEN.
* `local_region_id` - ID of local region.
* `opposite_region_id` - ID of opposite region.
* `status` - Status of the CEN Inter Region Bandwidth, including "Active" and "Modifying".
* `bandwidth_limit` - The bandwidth configured for the interconnected regions communication.