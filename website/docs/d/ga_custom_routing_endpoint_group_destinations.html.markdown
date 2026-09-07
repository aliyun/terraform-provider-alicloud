---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint_group_destinations"
sidebar_current: "docs-alicloud-datasource-ga-custom-routing-endpoint-group-destinations"
description: |-
  Provides a list of Global Accelerator (GA) Custom Routing Endpoint Group Destinations to the user.
---

# alicloud_ga_custom_routing_endpoint_group_destinations

This data source provides the Global Accelerator (GA) Custom Routing Endpoint Group Destinations of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_custom_routing_endpoint_group_destinations" "ids" {
  ids            = ["example_id"]
  accelerator_id = "your_accelerator_id"
}
output "ga_custom_routing_endpoint_group_destinations_id_1" {
  value = data.alicloud_ga_custom_routing_endpoint_group_destinations.ids.custom_routing_endpoint_group_destinations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Custom Routing Endpoint Group Destination IDs.
* `accelerator_id` - (Required, ForceNew) The ID of the GA instance.
* `listener_id` - (Optional, ForceNew) The ID of the listener.
* `endpoint_group_id` - (Optional, ForceNew) The ID of the endpoint group.
* `protocols` - (Optional, ForceNew) The backend service protocol of the endpoint group. Valid values: `TCP`, `UDP`, `TCP, UDP`.
* `from_port` - (Optional, ForceNew) The start port of the backend service port range of the endpoint group. The `from_port` value must be smaller than or equal to the `to_port` value. Valid values: `1` to `65499`.
* `to_port` - (Optional, ForceNew) The end port of the backend service port range of the endpoint group. The `from_port` value must be smaller than or equal to the `to_port` value. Valid values: `1` to `65499`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `custom_routing_endpoint_group_destinations` - A list of Custom Routing Endpoint Group Destinations. Each element contains the following attributes:
  * `id` - The id of the Global Accelerator Custom Routing Endpoint Group Destination. It formats as `<endpoint_group_id>:<custom_routing_endpoint_group_destination_id>`.  
  * `endpoint_group_id` - The ID of the Custom Routing Endpoint Group.
  * `custom_routing_endpoint_group_destination_id` - The ID of the Custom Routing Endpoint Group Destination.
  * `accelerator_id` - The ID of the GA instance.
  * `listener_id` - The ID of the listener.
  * `protocols` - The backend service protocol of the endpoint group.
  * `from_port` - The start port of the backend service port range of the endpoint group.
  * `to_port` - The end port of the backend service port range of the endpoint group.
