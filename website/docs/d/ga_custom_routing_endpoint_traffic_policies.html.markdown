---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint_traffic_policies"
sidebar_current: "docs-alicloud-datasource-ga-custom-routing-endpoint-traffic-policies"
description: |-
  Provides a list of Global Accelerator (GA) Custom Routing Endpoint Traffic Policies to the user.
---

# alicloud_ga_custom_routing_endpoint_traffic_policies

This data source provides the Global Accelerator (GA) Custom Routing Endpoint Traffic Policies of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_custom_routing_endpoint_traffic_policies" "ids" {
  ids            = ["example_id"]
  accelerator_id = "your_accelerator_id"
}

output "ga_custom_routing_endpoint_traffic_policies_id_1" {
  value = data.alicloud_ga_custom_routing_endpoint_traffic_policies.ids.custom_routing_endpoint_traffic_policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Custom Routing Endpoint Traffic Policy IDs.
* `accelerator_id` - (Required, ForceNew) The ID of the GA instance to which the traffic policies belong.
* `listener_id` - (Optional, ForceNew) The ID of the listener to which the traffic policies belong.
* `endpoint_group_id` - (Optional, ForceNew) The ID of the endpoint group to which the traffic policies belong.
* `endpoint_id` - (Optional, ForceNew) The ID of the endpoint to which the traffic policies belong.
* `address` - (Optional, ForceNew) The IP addresses of the traffic policies.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `custom_routing_endpoint_traffic_policies` - A list of Custom Routing Endpoint Traffic Policies. Each element contains the following attributes:
  * `id` - The id of the Global Accelerator Custom Routing Endpoint Traffic Policy. It formats as `<endpoint_id>:<custom_routing_endpoint_traffic_policy_id>`.
  * `endpoint_id` - The ID of the Custom Routing Endpoint.
  * `custom_routing_endpoint_traffic_policy_id` - The ID of the Custom Routing Endpoint Traffic Policy.
  * `accelerator_id` - The ID of the GA instance to which the endpoint belongs.
  * `listener_id` - The ID of the custom routing listener to which the endpoint belongs.
  * `endpoint_group_id` - The ID of the Custom Routing Endpoint Group.
  * `address` - The IP address of the traffic policy.
  * `port_ranges` - The port range of the traffic policy.
    * `from_port` - The first port of the port range.
    * `to_port` - The last port of the port range.
  