---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_port_mappings"
sidebar_current: "docs-alicloud-datasource-ga-custom-routing-port-mappings"
description: |-
  Provides a list of Global Accelerator (GA) Custom Routing Port Mappings to the user.
---

# alicloud_ga_custom_routing_port_mappings

This data source provides the Global Accelerator (GA) Custom Routing Port Mappings of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_custom_routing_port_mappings" "default" {
  accelerator_id = "your_accelerator_id"
}

output "ga_custom_routing_port_mappings_accelerator_id_1" {
  value = data.alicloud_ga_custom_routing_port_mappings.default.custom_routing_port_mappings.0.accelerator_id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the GA instance.
* `listener_id` - (Optional, ForceNew) The ID of the listener.
* `endpoint_group_id` - (Optional, ForceNew) The ID of the endpoint group.
* `status` - (Optional, ForceNew) The access policy of traffic for the backend instance. Valid Values: `allow`, `deny`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `custom_routing_port_mappings` - A list of Custom Routing Port Mappings. Each element contains the following attributes:
  * `accelerator_id` - The ID of the GA instance.
  * `listener_id` - The ID of the listener.
  * `endpoint_group_id` - The ID of the endpoint group.
  * `endpoint_id` - The ID of the endpoint.
  * `accelerator_port` - The acceleration port.
  * `vswitch` - The ID of the endpoint (vSwitch).
  * `endpoint_group_region` - The ID of the region in which the endpoint group resides.
  * `protocols` - The protocol of the backend service.
  * `destination_socket_address` - The service IP address and port of the backend instance.
    * `ip_address` - The service IP address of the backend instance.
    * `port` - The service port of the backend instance.
  * `status` - The access policy of traffic for the backend instance.
  