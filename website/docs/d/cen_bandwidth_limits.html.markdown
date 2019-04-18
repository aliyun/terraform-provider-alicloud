---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidth_limits"
sidebar_current: "docs-alicloud-datasource-cen-bandwidth-limits"
description: |-
    Provides a list of CEN Bandwidth Limits owned by an Alibaba Cloud account.
---

# alicloud\_cen\_bandwidth\_limits

This data source provides CEN Bandwidth Limits available to the user.

## Example Usage

```
data "alicloud_cen_bandwidth_limits" "bwl"{
	instance_ids = ["cen-id1"]
}

output "first_cen_bandwidth_limits_local_region_id" {
  value = "${data.alicloud_cen_bandwidth_packages.bwl.limits.0.local_region_id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Optional) A list of CEN instances IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `limits` - A list of CEN Bandwidth Limits. Each element contains the following attributes:
  * `instance_id` - ID of the CEN instance.
  * `local_region_id` - ID of local region.
  * `opposite_region_id` - ID of opposite region.
  * `status` - Status of the CEN Bandwidth Limit, including "Active" and "Modifying".
  * `bandwidth_limit` - The bandwidth limit configured for the interconnected regions communication.