---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoints"
sidebar_current: "docs-alicloud-datasource-ga-custom-routing-endpoints"
description: |-
  Provides a list of Global Accelerator (GA) Custom Routing Endpoints to the user.
---

# alicloud_ga_custom_routing_endpoints

This data source provides the Global Accelerator (GA) Custom Routing Endpoints of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_custom_routing_endpoints" "ids" {
  ids            = ["example_id"]
  accelerator_id = "your_accelerator_id"
}

output "ga_custom_routing_endpoints_id_1" {
  value = data.alicloud_ga_custom_routing_endpoints.ids.custom_routing_endpoints.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Custom Routing Endpoint IDs.
* `accelerator_id` - (Required, ForceNew) The ID of the GA instance.
* `listener_id` - (Optional, ForceNew) The ID of the custom routing listener.
* `endpoint_group_id` - (Optional, ForceNew) The ID of the endpoint group.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `custom_routing_endpoints` - A list of Custom Routing Endpoints. Each element contains the following attributes:
  * `id` - The id of the Global Accelerator Custom Routing Endpoint. It formats as `<endpoint_group_id>:<custom_routing_endpoint_id>`.
  * `endpoint_group_id` - The ID of the Custom Routing Endpoint Group.
  * `custom_routing_endpoint_id` - The ID of the Custom Routing Endpoint.
  * `accelerator_id` - The ID of the GA instance with which the endpoint is associated.
  * `listener_id` - The ID of the listener with which the endpoint is associated.
  * `endpoint` - The ID of the endpoint (vSwitch).
  * `type` - The backend service type of the endpoint.
  * `traffic_to_endpoint_policy` - The access policy of traffic for the specified endpoint.
  