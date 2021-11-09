---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eip_addresses"
sidebar_current: "docs-alicloud-datasource-eip-addresses"
description: |-
  Provides a list of Eip Addresses to the user.
---

# alicloud\_eip\_addresses

This data source provides the Eip Addresses of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.126.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_eip_addresses" "example" {
  ids        = ["eip-bp1jvx5ki6c********"]
  name_regex = "the_resource_name"
}

output "first_eip_address_id" {
  value = data.alicloud_eip_addresses.example.addresses.0.id
}
```

## Argument Reference

The following arguments are supported:

* `associated_instance_id` - (Optional, ForceNew) The associated instance id.
* `associated_instance_type` - (Optional, ForceNew) The associated instance type.
* `dry_run` - (Optional, ForceNew) The dry run.
* `ip_address` - (Optional, ForceNew) The eip address.
* `address_name` - (Optional, ForceNew) The eip name.
* `enable_details` - (Optional) Default to `tue`. Set it to `false` can hidden the `tags` to output.
* `ids` - (Optional, ForceNew, Computed)  A list of Address IDs.
* `include_reservation_data` - (Optional, ForceNew) The include reservation data. Valid values: `BGP` and `BGP_PRO`. 
* `isp` - (Optional, ForceNew) The Internet service provider (ISP). Valid values `BGP` and `BGP_PRO`.
* `lock_reason` - (Optional, ForceNew) The lock reason.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Address name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `payment_type` - (Optional, ForceNew) The billing method of the EIP. Valid values: `Subscription` and `PayAsYouGo`. 
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `segment_instance_id` - (Optional, ForceNew) The IDs of the contiguous EIPs.  This value is returned only when contiguous EIPs are specified.
* `status` - (Optional, ForceNew) The status of the EIP. Valid values:  `Associating`: The EIP is being associated. `Unassociating`: The EIP is being disassociated. `InUse`: The EIP is allocated. `Available`:The EIP is available.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Address names.
* `addresses` - A list of Eip Addresses. Each element contains the following attributes:
	* `address_name` - The name of the EIP.
	* `allocation_id` - The ID of the EIP.
	* `available_regions` - The ID of the region to which the EIP belongs.
	* `bandwidth` - The maximum bandwidth of the EIP. Unit: Mbit/s.
	* `bandwidth_package_bandwidth` - The bandwidth value of the EIP bandwidth plan with which the EIP is associated.
	* `bandwidth_package_id` - The ID of the EIP bandwidth plan.
	* `bandwidth_package_type` - The type of the bandwidth. Only CommonBandwidthPackage (an EIP bandwidth plan) is returned.
	* `create_time` - The time when the EIP was created.
	* `deletion_protection` - Indicates whether deletion protection is enabled.
	* `description` - The description of the EIP.
	* `expired_time` - The expiration date. The time follows the ISO 8601 standard and is displayed in UTC. Format: YYYY-MM-DDThh:mmZ.
	* `has_reservation_data` - Indicates whether renewal data is included. This parameter returns true only when the parameter IncludeReservationData is set to true, and some orders have not taken effect.
	* `hd_monitor_status` - Indicates whether fine-grained monitoring is enabled for the EIP.
	* `id` - The ID of the Address.
	* `instance_id` - The ID of the instance with which the EIP is associated.
	* `instance_region_id` - The region ID of the associated resource.
	* `instance_type` - The type of the instance with which the EIP is associated.
	* `internet_charge_type` - The metering method of the EIP.
	* `ip_address` - The IP address of the EIP.
	* `isp` - The Internet service provider (ISP).
	* `operation_locks` - The details about the locked EIP.
	* `payment_type` - The billing method of the EIP.
	* `reservation_active_time` - The time when the renewal takes effect.
	* `reservation_bandwidth` - The bandwidth after the renewal takes effect.
	* `reservation_internet_charge_type` - The metering method of the renewal. 
	* `reservation_order_type` - The type of the renewal order. 
	* `resource_group_id` - The ID of the resource group.
	* `second_limited` - Indicates whether level-2 throttling is configured.
	* `segment_instance_id` - The IDs of the contiguous EIPs.  
	* `status` - The status of the EIP. 
	* `tags` - A mapping of tags to assign to the resource.
