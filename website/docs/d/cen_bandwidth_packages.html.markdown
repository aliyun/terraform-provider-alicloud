---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidth_packages"
sidebar_current: "docs-alicloud-datasource-cen-bandwidth-packages"
description: |-
    Provides a list of CEN Bandwidth Packages owned by an Alibaba Cloud account.
---

# alicloud\_cen\_bandwidth\_packages

This data source provides CEN Bandwidth Packages available to the user.

## Example Usage

```terraform
data "alicloud_cen_bandwidth_packages" "example" {
  instance_id = "cen-id1"
  name_regex  = "^foo"
}

output "first_cen_bandwidth_package_id" {
  value = data.alicloud_cen_bandwidth_packages.example.packages.0.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional) ID of a CEN instance.
* `ids` - (Optional) Limit search to a list of specific CEN Bandwidth Package IDs.
* `name_regex` - (Optional) A regex string to filter CEN Bandwidth Package by name.
* `include_reservation_data` - (Optional, ForceNew, Available in 1.98.0+) -Indicates whether to include renewal data. Valid values: `true`: Return renewal data in the response. `false`: Do not return renewal data in the response.
* `status` - (Optional, ForceNew, Available in 1.98.0+) Status of the CEN Bandwidth Package in CEN instance, Valid value: `Idle` and `InUse`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of specific CEN Bandwidth Package IDs.
* `names` (Available in 1.98.0+) - A list of CEN Bandwidth Package Names.
* `packages` - A list of CEN bandwidth package. Each element contains the following attributes:
  * `id` - ID of the CEN Bandwidth Package.
  * `name` - Name of the CEN Bandwidth Package.
  * `description` - Description of the CEN Bandwidth Package.
  * `business_status` - Status of the CEN Bandwidth Package, including `Normal`, `FinancialLocked` and `SecurityLocked`.
  * `status` - Status of the CEN Bandwidth Package in CEN instance, including `Idle` and `InUse`.
  * `instance_id` - The ID of the CEN instance that are associated with the bandwidth package.
  * `bandwidth` - The bandwidth in Mbps of the CEN bandwidth package.
  * `bandwidth_package_charge_type` - The billing method, including `POSTPAY` and `PREPAY`.
  * `geographic_region_a_id` - Region ID of the interconnected regions.
  * `geographic_region_b_id` - Region ID of the interconnected regions.
  * `cen_bandwidth_package_id` - The ID of the bandwidth package.
  * `cen_bandwidth_package_name` - The name of the bandwidth package.
  * `cen_ids` - The list of CEN instances that are associated with the bandwidth package.
  * `geographic_span_id` - The area ID of the cross-area connection.
  * `has_reservation_data` - Indicates whether renewal data is involved.
  * `is_cross_border` - Indicates whether the bandwidth package is a cross-border bandwidth package.
  * `payment_type` - The billing method of the bandwidth package.
  * `reservation_active_time` - The expiration time of the temporary upgrade.
  * `reservation_bandwidth` - The restored bandwidth after the temporary upgrade.
  * `reservation_internet_charge_type` - The billing method after the configuration change.
  * `reservation_order_type` - The type of the configuration change.
