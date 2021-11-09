---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_endpoint_groups"
sidebar_current: "docs-alicloud-datasource-ga-endpoint-groups"
description: |-
  Provides a list of Global Accelerator (GA) Endpoint Groups to the user.
---

# alicloud\_ga\_endpoint\_groups

This data source provides the Global Accelerator (GA) Endpoint Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_endpoint_groups" "example" {
  accelerator_id = "example_value"
  ids            = ["example_value"]
  name_regex     = "the_resource_name"
}

output "first_ga_endpoint_group_id" {
  value = data.alicloud_ga_endpoint_groups.example.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator instance to which the endpoint group will be added.
* `endpoint_group_type` - (Optional, ForceNew) The endpoint group type. Valid values: `default`, `virtual`. Default value is `default`.
* `ids` - (Optional, ForceNew, Computed)  A list of Endpoint Group IDs.
* `listener_id` - (Optional, ForceNew) The ID of the listener that is associated with the endpoint group.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Endpoint Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the endpoint group. Valid values: `active`, `configuring`, `creating`, `init`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Endpoint Group names.
* `groups` - A list of Ga Endpoint Groups. Each element contains the following attributes:
	* `description` - The description of the endpoint group.
	* `endpoint_configurations` - The endpointConfigurations of the endpoint group.
		* `probe_port` - Probe Port.
		* `probe_protocol` - Probe Protocol.
		* `type` - The type of Endpoint N in the endpoint group.
		* `weight` - The weight of Endpoint N in the endpoint group.
		* `enable_clientip_preservation` - Indicates whether client IP addresses are reserved.
		* `endpoint` - The IP address or domain name of Endpoint N in the endpoint group.
	* `endpoint_group_id` - The endpoint_group_id of the Endpoint Group.
	* `endpoint_group_region` - The ID of the region where the endpoint group is deployed.
	* `health_check_interval_seconds` - The interval between two consecutive health checks. Unit: seconds.
	* `health_check_path` - The path specified as the destination of the targets for health checks.
	* `health_check_port` - The port that is used for health checks.
	* `health_check_protocol` - The protocol that is used to connect to the targets for health checks.
	* `id` - The ID of the Endpoint Group.
	* `listener_id` - The ID of the listener that is associated with the endpoint group.
	* `name` - The name of the endpoint group.
	* `port_overrides` - Mapping between listening port and forwarding port of boarding point.
		* `endpoint_port` - Forwarding port.
		* `listener_port` - Listener port.
	* `status` - The status of the endpoint group.
	* `threshold_count` - The number of consecutive failed heath checks that must occur before the endpoint is deemed unhealthy.
	* `traffic_percentage` - The weight of the endpoint group when the corresponding listener is associated with multiple endpoint groups.
