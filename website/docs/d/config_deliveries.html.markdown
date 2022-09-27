---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_deliveries"
sidebar_current: "docs-alicloud-datasource-config-deliveries"
description: |-
  Provides a list of Config Deliveries to the user.
---

# alicloud\_config\_deliveries

This data source provides the Config Deliveries of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.171.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_config_deliveries" "ids" {
  ids = ["example_id"]
}
output "config_delivery_id_1" {
  value = data.alicloud_config_deliveries.ids.deliveries.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Delivery IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by delivery channel name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the delivery method. Valid values: `0`: The delivery method is disabled. `1`: The delivery destination is enabled.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `deliveries` - A list of Config Deliveries. Each element contains the following attributes:
	* `account_id` - The Aliyun User Id.
	* `configuration_item_change_notification` - Open or close delivery configuration change history.
	* `configuration_snapshot` - Open or close timed snapshot of shipping resources.
	* `delivery_channel_assume_role_arn` - The Alibaba Cloud Resource Name (ARN) of the role to be assumed by the delivery method.
	* `delivery_channel_condition` - The rule attached to the delivery method.
	* `delivery_channel_id` - The ID of the delivery method.
	* `delivery_channel_name` - The name of the delivery method.
	* `delivery_channel_target_arn` - The ARN of the delivery destination.
	* `delivery_channel_type` - The type of the delivery method.
	* `description` - The description of the delivery method.
	* `id` - The ID of the Delivery.
	* `non_compliant_notification` - Open or close non-compliance events of delivery resources.
	* `oversized_data_oss_target_arn` - The oss ARN of the delivery channel when the value data oversized limit.
	* `status` - The status of the delivery method. Valid values: `0`: The delivery method is disabled. `1`: The delivery destination is enabled.