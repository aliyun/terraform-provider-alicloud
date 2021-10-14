---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_traffic_mirror_filters"
sidebar_current: "docs-alicloud-datasource-vpc-traffic-mirror-filters"
description: |-
  Provides a list of Vpc Traffic Mirror Filters to the user.
---

# alicloud\_vpc\_traffic\_mirror\_filters

This data source provides the Vpc Traffic Mirror Filters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_traffic_mirror_filters" "ids" {
  ids = ["example_id"]
}
output "vpc_traffic_mirror_filter_id_1" {
  value = data.alicloud_vpc_traffic_mirror_filters.ids.filters.0.id
}

data "alicloud_vpc_traffic_mirror_filters" "nameRegex" {
  name_regex = "^my-TrafficMirrorFilter"
}
output "vpc_traffic_mirror_filter_id_2" {
  value = data.alicloud_vpc_traffic_mirror_filters.nameRegex.filters.0.id
}

data "alicloud_vpc_traffic_mirror_filters" "filterName" {
  traffic_mirror_filter_name = "example_traffic_mirror_filter_name"
}
output "vpc_traffic_mirror_filter_id_3" {
  value = data.alicloud_vpc_traffic_mirror_filters.filterName.filters.0.id
}

data "alicloud_vpc_traffic_mirror_filters" "status" {
  status = "^my-TrafficMirrorFilter"
}
output "vpc_traffic_mirror_filter_id_4" {
  value = data.alicloud_vpc_traffic_mirror_filters.status.filters.0.id
}
            
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Traffic Mirror Filter IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Traffic Mirror Filter name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The state of the filter. Valid values:`Creating`, `Created`, `Modifying` and `Deleting`. `Creating`: The filter is being created. `Created`: The filter is created. `Modifying`: The filter is being modified. `Deleting`: The filter is being deleted.
* `traffic_mirror_filter_ids` - (Optional, ForceNew) The traffic mirror filter ids.
* `traffic_mirror_filter_name` - (Optional, ForceNew) The name of the filter. The name must be `2` to `128` characters in length, and can contain digits, periods (.), underscores (_), and hyphens (-). It must start with a letter and cannot start with `http://` or `https://`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Traffic Mirror Filter names.
* `filters` - A list of Vpc Traffic Mirror Filters. Each element contains the following attributes:
	* `egress_rules` - The list of details about outbound rules.
		* `destination_port_range` - The destination port range of the outbound traffic.
		* `protocol` - The transport protocol used by outbound traffic that needs to be mirrored. Valid values: `ALL`, `ICMP`, `TCP`, `UDP`.
		* `source_cidr_block` - The source CIDR block of the outbound traffic.
		* `source_port_range` - The source port range of the outbound traffic.
		* `traffic_direction` - The direction of the network traffic. Valid values: `egress` or `ingress`. `egress`: outbound `ingress`: inbound.
		* `action` - The collection policy of the outbound rule. Valid values: `accept` or `drop`. `accept`: collects network traffic. `drop`: does not collect network traffic.
		* `destination_cidr_block` - The destination CIDR block of the outbound traffic.
		* `priority` - The priority of the outbound rule. A smaller value indicates a higher priority.
		* `traffic_mirror_filter_id` - The ID of the filter associated with the outbound rule.
		* `traffic_mirror_filter_rule_id` - The ID of the outbound rule.
		* `traffic_mirror_filter_rule_status` - The state of the outbound rule. Valid values:`Creating`, `Created`, `Modifying` and `Deleting`.
	* `id` - The ID of the Traffic Mirror Filter.
	* `ingress_rules` - The list of details about inbound rules.
		* `destination_port_range` - The destination port range of the inbound traffic.
		* `protocol` - The transport protocol used by inbound traffic that needs to be mirrored. Valid values: `ALL`, `ICMP`, `TCP`, `UDP`.
		* `source_cidr_block` - The source CIDR block of the inbound traffic.
		* `source_port_range` - The source port range of the inbound traffic.
		* `traffic_direction` - The direction of the network traffic. Valid values: `egress` or `ingress`. `egress`: outbound `ingress`: inbound.
		* `action` - The collection policy of the inbound rule. Valid values: `accept` or `drop`. `accept`: collects network traffic. `drop`: does not collect network traffic.
		* `destination_cidr_block` - The destination CIDR block of the inbound traffic.
		* `priority` - The priority of the inbound rule. A smaller value indicates a higher priority.
		* `traffic_mirror_filter_id` - The ID of the filter associated with the inbound rule.
		* `traffic_mirror_filter_rule_id` - The ID of the inbound rule.
		* `traffic_mirror_filter_rule_status` - The state of the inbound rule. Valid values:`Creating`, `Created`, `Modifying` and `Deleting`.
	* `status` - The state of the filter. Valid values:`Creating`, `Created`, `Modifying` and `Deleting`. `Creating`: The filter is being created. `Created`: The filter is created. `Modifying`: The filter is being modified. `Deleting`: The filter is being deleted.
	* `traffic_mirror_filter_description` - The description of the filter.
	* `traffic_mirror_filter_id` - The ID of the filter.
	* `traffic_mirror_filter_name` - The name of the filter.