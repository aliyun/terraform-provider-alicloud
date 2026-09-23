---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_traffic_marking_policies"
sidebar_current: "docs-alicloud-datasource-cen-traffic-marking-policies"
description: |-
  Provides a list of Cen Traffic Marking Policies to the user.
---

# alicloud\_cen\_traffic\_marking\_policies

This data source provides the Cen Traffic Marking Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.173.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_traffic_marking_policies" "ids" {
  transit_router_id = "example_value"
  ids               = ["example_value-1", "example_value-2"]
}
output "cen_traffic_marking_policy_id_1" {
  value = data.alicloud_cen_traffic_marking_policies.ids.policies.0.id
}

data "alicloud_cen_traffic_marking_policies" "nameRegex" {
  transit_router_id = "example_value"
  name_regex        = "^my-TrafficMarkingPolicy"
}
output "cen_traffic_marking_policy_id_2" {
  value = data.alicloud_cen_traffic_marking_policies.nameRegex.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Traffic Marking Policy IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Traffic Marking Policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource.  Valid values: `Active`, `Creating`, `Deleting`, `Updating`.
* `traffic_marking_policy_description` - (Optional, ForceNew) The traffic marking policy description.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Traffic Marking Policy names.
* `policies` - A list of Cen Traffic Marking Policies. Each element contains the following attributes:
	* `description` - The description of the Traffic Marking Policy.
	* `id` - The ID of the resource. The value is formatted `<transit_router_id>:<traffic_marking_policy_id>`.
	* `marking_dscp` - The DSCP(Differentiated Services Code Point) of the Traffic Marking Policy.
	* `priority` - The Priority of the Traffic Marking Policy.
	* `transit_router_id` - The ID of the transit router.
	* `status` - The status of the resource.
	* `traffic_marking_policy_id` - The ID of the Traffic Marking Policy.
	* `traffic_marking_policy_name` - The name of the Traffic Marking Policy.