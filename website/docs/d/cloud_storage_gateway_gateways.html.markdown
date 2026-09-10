---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateways"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-gateways"
description: |-
  Provides a list of Cloud Storage Gateway Gateways to the user.
---

# alicloud\_cloud\_storage\_gateway\_gateways

This data source provides the Cloud Storage Gateway Gateways of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_storage_gateway_storage_bundle" "example" {
  storage_bundle_name = "example_value"
}
data "alicloud_cloud_storage_gateway_gateways" "nameRegex" {
  storage_bundle_id = alicloud_cloud_storage_gateway_storage_bundle.example.id
  name_regex        = "^my-Gateway"
}
output "cloud_storage_gateway_gateway_id" {
  value = data.alicloud_cloud_storage_gateway_gateways.nameRegex.gateways.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Gateway IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Gateway name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) gateway status.
* `storage_bundle_id` - (Required, ForceNew) storage bundle id.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Gateway names.
* `gateways` - A list of Cloud Storage Gateway Gateways. Each element contains the following attributes:
	* `activated_time` - gateway .
	* `category` - gateway category.
	* `create_time` - gateway created timestamp in second format.
	* `description` - gateway description.
	* `ecs_instance_id` - gateway ecs instance id.
	* `expire_status` - gateway expiration status.
	* `expired_time` - gateway expiration timestamp in second format.
	* `gateway_class` - gateway class.
	* `gateway_id` - gateway id.
	* `gateway_name` - gateway name.
	* `gateway_version` - gateway version.
	* `id` - The ID of the Gateway.
	* `inner_ip` - gateway service ip.
	* `ip` - gateway public ip.
	* `is_release_after_expiration` - whether subscription gateway is released after expiration or not.
	* `location` - gateway location.
	* `payment_type` - gateway payment type. The Payment type of gateway. The valid value: `PayAsYouGo`, `Subscription`.
	* `public_network_bandwidth` - gateway public network bandwidth.
	* `status` - gateway status.
	* `storage_bundle_id` - storage bundle id.
	* `task_id` - gateway task id.
	* `type` - gateway type.
	* `vpc_id` - gateway vpc id.
	* `vswitch_id` - The vswitch id.
