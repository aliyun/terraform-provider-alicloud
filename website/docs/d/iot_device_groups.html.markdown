---
subcategory: "Internet of Things (Iot)"
layout: "alicloud"
page_title: "Alicloud: alicloud_iot_device_groups"
sidebar_current: "docs-alicloud-datasource-iot-device-groups"
description: |-
  Provides a list of Iot Device Groups to the user.
---

# alicloud\_iot\_device\_groups

This data source provides the Iot Device Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_iot_device_groups" "ids" {}
output "iot_device_group_id_1" {
  value = data.alicloud_iot_device_groups.ids.groups.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `group_name` - (Optional, ForceNew) The GroupName of the device group.
* `ids` - (Optional, ForceNew, Computed)  A list of device group IDs.
* `name_regex` - (Optional) A regex string to filter CEN instances by name.
* `iot_instance_id` - (Optional, ForceNew) The id of the Iot Instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `super_group_id` - (Optional, ForceNew) The id of the SuperGroup.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `groups` - A list of Iot Device Groups. Each element contains the following attributes:
	* `create_time` - The Group CreateTime.
	* `device_active` - The Group Number of activated devices.
	* `device_count` - The Group Total number of devices.
	* `device_online` - The Group Number of online devices.
	* `error_message` - The Error_Message of the device group.
	* `group_desc` - The GroupDesc of the device group.
	* `group_id` - The GroupId of the device group.
	* `group_name` - The GroupName of the device group.
	* `id` - The ID of the device group.
	* `success` - Whether the call is successful.
