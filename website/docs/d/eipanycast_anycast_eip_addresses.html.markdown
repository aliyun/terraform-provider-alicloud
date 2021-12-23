---
subcategory: "Anycast Elastic IP Address (Eipanycast)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eipanycast_anycast_eip_addresses"
sidebar_current: "docs-alicloud-datasource-eipanycast-anycast-eip-addresses"
description: |-
  Provides a list of Anycast Eip Addresses to the user.
---

# alicloud\_eipanycast\_anycast\_eip\_addresses

This data source provides the Eipanycast Anycast Eip Addresses of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_eipanycast_anycast_eip_addresses" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_eipanycast_anycast_eip_address_id" {
  value = data.alicloud_eipanycast_anycast_eip_addresses.example.addresses.0.id
}
```

## Argument Reference

The following arguments are supported:

* `anycast_eip_address_name` - (Optional, ForceNew) Anycast EIP instance name.
* `bind_instance_ids` - (Optional, ForceNew) The bind instance ids.
* `business_status` - (Optional, ForceNew) The business status of the Anycast EIP instance. -`Normal`: Normal state. -`FinancialLocked`: The status of arrears locked.
* `ids` - (Optional, ForceNew, Computed)  A list of Anycast Eip Address IDs.
* `internet_charge_type` - (Optional, ForceNew) The billing method of Anycast EIP instance. `PayByBandwidth`: refers to the method of billing based on traffic.
* `ip_address` - (Optional, ForceNew)  Anycast EIP instance IP address.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Anycast Eip Address name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `payment_type` - (Optional, ForceNew) The payment model of Anycast EIP instance. `PayAsYouGo`: Refers to the post-paid mode. Default value is `PayAsYouGo`.
* `service_location` - (Optional, ForceNew) Anycast EIP instance access area. `international`: Refers to areas outside of Mainland China.
* `status` - (Optional, ForceNew) IP status。- `Associating`, `Unassociating`, `Allocated`, `Associated`, `Modifying`, `Releasing`, `Released`. Valid values: `Allocated`, `Associated`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Anycast Eip Address names.
* `addresses` - A list of Eipanycast Anycast Eip Addresses. Each element contains the following attributes:
	* `ali_uid` - Anycast EIP instance account ID.
	* `anycast_eip_address_name` - Anycast EIP instance name.
	* `anycast_eip_bind_info_list` -  AnycastEip binding information.
		* `bind_instance_id` - The bound cloud resource instance ID.
		* `bind_instance_region_id` -  The region ID of the bound cloud resource instance.
		* `bind_instance_type` - Bind the cloud resource instance type.
		* `bind_time` -  Binding time.
	* `anycast_id` -  Anycast EIP instance ID.
	* `bandwidth` -  The peak bandwidth of the Anycast EIP instance, in Mbps.
	* `bid` - Anycast EIP instance account BID.
	* `business_status` - The business status of the Anycast EIP instance. -`Normal`: Normal state. -`FinancialLocked`: The status of arrears locked.
	* `description` - Anycast EIP instance description.
	* `id` - The ID of the Anycast Eip Address.
	* `internet_charge_type` - The billing method of Anycast EIP instance. `PayByBandwidth`: refers to the method of billing based on traffic.
	* `ip_address` -  Anycast EIP instance IP address.
	* `payment_type` - The payment model of Anycast EIP instance. "PostPaid": Refers to the post-paid mode.
	* `service_location` - Anycast EIP instance access area. "international": Refers to areas outside of Mainland China.
	* `status` - IP status。- `Associating`, `Unassociating`, `Allocated`, `Associated`, `Modifying`, `Releasing`, `Released`.
