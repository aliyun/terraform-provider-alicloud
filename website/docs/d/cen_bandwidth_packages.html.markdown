---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidth_packages"
sidebar_current: "docs-alicloud-datasource-cen-bandwidth-packages"
description: |-
    Provides a list of CEN Bandwidth Packages owned by an Alibaba Cloud account.
---

# alicloud\_cen\_bandwidth\_packages

This data source provides CEN Bandwidth Packages available to the user.

## Example Usage

```
data "alicloud_cen_bandwidth_packages" "bwp" {
	instance_ids = ["cen-id1"]
	status = "Idle"
	name_regex="^foo"
}

output "first_cen_bandwidth_package_id" {
  value = "${data.alicloud_cen_bandwidth_packages.bwp.bandwidth_packages.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Optional) A list of CEN instances IDs.
* `status` - (Optional) Limit search to specific status - valid value is "Idle" or "InUse".
* `bandwidth_package_ids` - (Optional) Limit search to a list of specific CEN Bandwidth Package IDs.
* `name_regex` - (Optional) A regex string to filter CEN Bandwidth Package by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `bandwidth_packages` - A list of CEN bandwidth package. Each element contains the following attributes:
  * `id` - ID of the CEN Bandwidth Package.
  * `instance_id` - ID of CEN instance that owns the CEN Bandwidth Package.
  * `name` - Name of the CEN Bandwidth Package.
  * `description` - Description of the CEN Bandwidth Package.
  * `business_status` - Status of the CEN Bandwidth Package, including "Normal", "FinancialLocked" and "SecurityLocked".
  * `status` - Status of the CEN Bandwidth Package in CEN instance, including "Idle" and "InUse".
  * `bandwidth` - The bandwidth in Mbps of the CEN bandwidth package.
  * `creation_time` - Creation time of the CEN bandwidth package.
  * `bandwidth_package_charge_type` - The billing method, including "POSTPAY" and "PREPAY".
  * `geographic_region_a_id` - Region ID of the interconnected regions.
  * `geographic_region_b_id` - Region ID of the interconnected regions.