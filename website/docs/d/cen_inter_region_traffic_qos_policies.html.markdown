---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_inter_region_traffic_qos_policies"
sidebar_current: "docs-alicloud-datasource-cen-inter-region-traffic-qos-policies"
description: |-
  Provides a list of Cen Inter Region Traffic Qos Policies to the user.
---

# alicloud\_cen\_inter\_region\_traffic\_qos\_policies

This data source provides the Cen Inter Region Traffic Qos Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_inter_region_traffic_qos_policies" "ids" {
  ids                          = ["example_id"]
  transit_router_id            = "your_transit_router_id"
  transit_router_attachment_id = "your_transit_router_attachment_id"
}

output "cen_inter_region_traffic_qos_policy_id_0" {
  value = data.alicloud_cen_inter_region_traffic_qos_policies.ids.policies.0.id
}

data "alicloud_cen_inter_region_traffic_qos_policies" "nameRegex" {
  name_regex                   = "^my-name"
  transit_router_id            = "your_transit_router_id"
  transit_router_attachment_id = "your_transit_router_attachment_id"
}

output "cen_inter_region_traffic_qos_policy_id_1" {
  value = data.alicloud_cen_inter_region_traffic_qos_policies.nameRegex.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Inter Region Traffic Qos Policy IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Inter Region Traffic Qos Policy name.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router.
* `transit_router_attachment_id` - (Required, ForceNew) The ID of the inter-region connection.
* `traffic_qos_policy_id` - (Optional, ForceNew) The ID of the QoS policy.
* `traffic_qos_policy_name` - (Optional, ForceNew) The name of the QoS policy.
* `traffic_qos_policy_description` - (Optional, ForceNew) The description of the QoS policy.
* `status` - (Optional, ForceNew) The status of the traffic scheduling policy. Valid Value: `Creating`, `Active`, `Modifying`, `Deleting`, `Deleted`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Inter Region Traffic Qos Policy names.
* `policies` - A list of Cen Inter Region Traffic Qos Policies. Each element contains the following attributes:
	* `id` - The ID of the Inter Region Traffic Qos Policy.
	* `transit_router_id` - The ID of the transit router.
	* `transit_router_attachment_id` - The ID of the inter-region connection.
	* `inter_region_traffic_qos_policy_id` - The ID of the Inter Region Traffic Qos Policy.
	* `inter_region_traffic_qos_policy_name` - The name of the Inter Region Traffic Qos Policy.
	* `inter_region_traffic_qos_policy_description` - The description of the Inter Region Traffic Qos Policy.
	* `status` - The status of the Inter Region Traffic Qos Policy.
	