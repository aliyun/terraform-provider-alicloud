---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidth_packages"
sidebar_current: "docs-alicloud-datasource-cen-bandwidth-packages"
description: |-
    Provides a list of CEN Bandwidth Packages which owned of the CEN.
---

# alicloud\_cen\_bandwidth\_packages

The CEN bandwidth packages data source lists a number of CEN bandwidth package resource information of the CEN.

## Example Usage

```
data "alicloud_cen_bandwidth_packages" "multi_cen_bwp"{
	cen_ids = ["cen-abc123"]
	status = "Idle"
	name_regex="^foo"
}
```

## Argument Reference

The following arguments are supported:

* `cen_ids` - (Optional) Limit search to a list of specific CEN IDs, like ["cen-id1","cen-id2"].
* `status` - (Optional) Limit search to specific status - valid value is "Idle" or "InUse".
* `cen_bandwidth_package_ids` - (Optional) Limit search to a list of specific CEN Bandwidth Package IDs, like ["cen_bwp_id1", "cen_bwp_id2"].
* `name_regex` - (Optional) A regex string of CEN Bandwidth Package name.
* `output_file` - (Optional) The name of file that can save CEN Bandwidth Package data source after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the CEN Bandwidth Package.
* `cen_ids` - List of CEN IDs in the specified CEN Bandwidth Package.
* `name` - Name of the CEN Bandwidth Package.
* `description` - Description of the CEN Bandwidth Package.
* `business_status` - Status of the CEN Bandwidth Package in CEN, including "Idle" and "InUse".
* `status` - Status of the CEN Bandwidth Package, including "Normal", "FinancialLocked" and "SecurityLocked".
* `bandwidth` - The bandwidth in Mbps of the bandwidth package.
* `creation_time` - Time of creation.
* `bandwidth_package_charge_type` - The billing method, including "POSTPAY" and "PREPAY".
* `geographic_region_a_id` - Region ID of the interconnected regions.
* `geographic_region_b_id` - Region ID of the interconnected regions.