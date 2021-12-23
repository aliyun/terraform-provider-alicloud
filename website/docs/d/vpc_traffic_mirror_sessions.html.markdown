---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_traffic_mirror_sessions"
sidebar_current: "docs-alicloud-datasource-vpc-traffic-mirror-sessions"
description: |-
  Provides a list of Vpc Traffic Mirror Sessions to the user.
---

# alicloud\_vpc\_traffic\_mirror\_sessions

This data source provides the Vpc Traffic Mirror Sessions of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_traffic_mirror_sessions" "ids" {
  ids = ["example_id"]
}
output "vpc_traffic_mirror_session_id_1" {
  value = data.alicloud_vpc_traffic_mirror_sessions.ids.sessions.0.id
}

data "alicloud_vpc_traffic_mirror_sessions" "nameRegex" {
  name_regex = "^my-TrafficMirrorSession"
}
output "vpc_traffic_mirror_session_id_2" {
  value = data.alicloud_vpc_traffic_mirror_sessions.nameRegex.sessions.0.id
}

data "alicloud_vpc_traffic_mirror_sessions" "enabled" {
  ids     = ["example_id"]
  enabled = "false"
}
output "vpc_traffic_mirror_session_id_3" {
  value = data.alicloud_vpc_traffic_mirror_sessions.enabled.sessions.0.id
}

data "alicloud_vpc_traffic_mirror_sessions" "priority" {
  ids      = ["example_id"]
  priority = "1"
}
output "vpc_traffic_mirror_session_id_4" {
  value = data.alicloud_vpc_traffic_mirror_sessions.priority.sessions.0.id
}

data "alicloud_vpc_traffic_mirror_sessions" "filterId" {
  ids                      = ["example_id"]
  traffic_mirror_filter_id = "example_value"
}
output "vpc_traffic_mirror_session_id_5" {
  value = data.alicloud_vpc_traffic_mirror_sessions.filterId.sessions.0.id
}

data "alicloud_vpc_traffic_mirror_sessions" "sessionName" {
  ids                         = ["example_id"]
  traffic_mirror_session_name = "example_value"
}
output "vpc_traffic_mirror_session_id_6" {
  value = data.alicloud_vpc_traffic_mirror_sessions.sessionName.sessions.0.id
}

data "alicloud_vpc_traffic_mirror_sessions" "sourceId" {
  ids                      = ["example_id"]
  traffic_mirror_source_id = "example_value"
}
output "vpc_traffic_mirror_session_id_7" {
  value = data.alicloud_vpc_traffic_mirror_sessions.sourceId.sessions.0.id
}

data "alicloud_vpc_traffic_mirror_sessions" "targetId" {
  ids                      = ["example_id"]
  traffic_mirror_target_id = "example_value"
}
output "vpc_traffic_mirror_session_id_8" {
  value = data.alicloud_vpc_traffic_mirror_sessions.targetId.sessions.0.id
}

data "alicloud_vpc_traffic_mirror_sessions" "status" {
  ids    = ["example_id"]
  status = "Created"
}
output "vpc_traffic_mirror_session_id_9" {
  value = data.alicloud_vpc_traffic_mirror_sessions.status.sessions.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enabled` - (Optional, ForceNew) Specifies whether to enable traffic mirror sessions. default to `false`.
* `ids` - (Optional, ForceNew, Computed)  A list of Traffic Mirror Session IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Traffic Mirror Session name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `priority` - (Optional, ForceNew) The priority of the traffic mirror session. Valid values: `1` to `32766`. A smaller value indicates a higher priority. You cannot specify the same priority for traffic mirror sessions that are created in the same region with the same Alibaba Cloud account.
* `status` - (Optional, ForceNew) The state of the traffic mirror session. Valid values: `Creating`, `Created`, `Modifying` and `Deleting`.
* `traffic_mirror_filter_id` - (Optional, ForceNew) The ID of the filter.
* `traffic_mirror_session_name` - (Optional, ForceNew) The name of the traffic mirror session. The name must be `2` to `128` characters in length and can contain digits, underscores (_), and hyphens (-). It must start with a letter.
* `traffic_mirror_source_id` - (Optional, ForceNew) The ID of the mirror source. You can specify only an elastic network interface (ENI) as the mirror source.
* `traffic_mirror_target_id` - (Optional, ForceNew) The ID of the mirror destination. You can specify only an ENI or a Server Load Balancer (SLB) instance as a mirror destination.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Traffic Mirror Session names.
* `sessions` - A list of Vpc Traffic Mirror Sessions. Each element contains the following attributes:
	* `enabled` - Indicates whether traffic mirror sessions are enabled. default to `false`.
	* `id` - The ID of the Traffic Mirror Session.
	* `packet_length` - The maximum transmission unit (MTU).
	* `priority` - The priority of the traffic mirror session. A smaller value indicates a higher priority.
	* `status` - The state of the traffic mirror session. Valid values: `Creating`, `Created`, `Modifying` and `Deleting`.
	* `traffic_mirror_filter_id` - The ID of the filter.
	* `traffic_mirror_session_business_status` - The state of the traffic mirror session. Valid values: `Normal` or `FinancialLocked`. `Normal`: working as expected. `FinancialLocked`: locked due to overdue payments.
	* `traffic_mirror_session_description` - The description of the traffic mirror session.
	* `traffic_mirror_session_id` - The first ID of the resource.
	* `traffic_mirror_session_name` - The name of the traffic mirror session.
	* `traffic_mirror_source_ids` - The ID of the mirror source.
	* `traffic_mirror_target_id` - The ID of the mirror destination. You can specify only an ENI or a Server Load Balancer (SLB) instance as a mirror destination.
	* `traffic_mirror_target_type` - The type of the mirror destination. Valid values: `NetworkInterface` or `SLB`. `NetworkInterface`: an ENI. `SLB`: an internal-facing SLB instance
	* `virtual_network_id` - You can specify VNIs to distinguish different mirrored traffic.
