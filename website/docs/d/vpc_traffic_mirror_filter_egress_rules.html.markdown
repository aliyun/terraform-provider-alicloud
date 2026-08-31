---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_traffic_mirror_filter_egress_rules"
sidebar_current: "docs-alicloud-datasource-vpc-traffic-mirror-filter-egress-rules"
description: |-
  Provides a list of Vpc Traffic Mirror Filter Egress Rules to the user.
---

# alicloud\_vpc\_traffic\_mirror\_filter\_egress\_rules

This data source provides the Vpc Traffic Mirror Filter Egress Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_traffic_mirror_filter_egress_rules" "ids" {
  traffic_mirror_filter_id = "example_traffic_mirror_filter_id"
  ids                      = ["example_id"]
}
output "vpc_traffic_mirror_filter_egress_rule_id_1" {
  value = data.alicloud_vpc_traffic_mirror_filter_egress_rules.ids.rules.0.id
}

data "alicloud_vpc_traffic_mirror_filter_egress_rules" "status" {
  traffic_mirror_filter_id = "example_traffic_mirror_filter_id"
  ids                      = ["example_id"]
  status                   = "Created"
}
output "vpc_traffic_mirror_filter_egress_rule_id_2" {
  value = data.alicloud_vpc_traffic_mirror_filter_egress_rules.status.rules.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Traffic Mirror Filter Egress Rule IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values:`Creating`, `Created`, `Modifying` and `Deleting`.
* `traffic_mirror_filter_id` - (Required, ForceNew) The ID of the Traffic Mirror Filter.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `rules` - A list of Vpc Traffic Mirror Filter Egress Rules. Each element contains the following attributes:
	* `destination_cidr_block` - The destination CIDR block of the outbound traffic.
	* `destination_port_range` - The destination port range of the outbound traffic.
	* `id` - The ID of the Traffic Mirror Filter Egress Rule.
	* `priority` - The priority of the outbound rule. A smaller value indicates a higher priority. The maximum value is `10`, which indicates that you can configure at most 10 inbound rules for a filter.
	* `rule_action` - The collection policy of the inbound rule. Valid values: `accept` or `drop`. `accept`: collects network traffic. `drop`: does not collect network traffic.
	* `protocol` - The transport protocol used by outbound traffic that needs to be mirrored. Valid values: `ALL`, `ICMP`, `TCP`, `UDP`.
	* `source_cidr_block` - The source CIDR block of the outbound traffic.
	* `source_port_range` - The source port range of the outbound traffic.
	* `status` - The status of the resource. Valid values:`Creating`, `Created`, `Modifying` and `Deleting`.
	* `traffic_mirror_filter_id` - The ID of the filter associated with the outbound rule.
	* `traffic_mirror_filter_rule_id` - The first ID of the resource.
