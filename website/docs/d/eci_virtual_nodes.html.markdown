---
subcategory: "Elastic Container Instance (ECI)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_virtual_nodes"
sidebar_current: "docs-alicloud-datasource-eci-virtual-nodes"
description: |-
  Provides a list of Eci Virtual Nodes to the user.
---

# alicloud\_eci\_virtual\_nodes

This data source provides the Eci Virtual Nodes of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.145.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_eci_virtual_nodes" "ids" {
  ids = ["example_value-1", "example_value-2"]
}
output "eci_virtual_node_id_1" {
  value = data.alicloud_eci_virtual_nodes.ids.nodes.0.id
}

data "alicloud_eci_virtual_nodes" "nameRegex" {
  name_regex = "^my-VirtualNode"
}
output "eci_virtual_node_id_2" {
  value = data.alicloud_eci_virtual_nodes.nameRegex.nodes.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Virtual Node IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Virtual Node name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The resource group ID. If when you create a GPU does not specify a resource group instance will automatically add the account's default resource group.
* `security_group_id` - (Optional, ForceNew) VNode itself and by VNode created (ECI) the security group used by.
* `status` - (Optional, ForceNew) The Status of the virtual node. Valid values: `Cleaned`, `Failed`, `Pending`, `Ready`.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `virtual_node_name` - (Optional, ForceNew) The name of the virtual node.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Virtual Node names.
* `nodes` - A list of Eci Virtual Nodes. Each element contains the following attributes:
	* `cpu` - The Number of CPU.
	* `create_time` - The creation time of the virtual node.
	* `eni_instance_id` - The ENI instance ID.
	* `events` - The event list.
		* `type` - The Event type.
		* `count` - The number of occurrences.
		* `first_timestamp` - The first presentation time stamp.
		* `last_timestamp` - The most recent time stamp.
		* `message` - The event of the message body.
		* `name` - The name of the event.
		* `reason` - The causes of the incident.
	* `id` - The ID of the Virtual Node.
	* `internet_ip` - The IP address of a public network.
	* `intranet_ip` - The private IP address of the RDS instance.
	* `memory` - The memory size.
	* `ram_role_name` - The ram role.
	* `resource_group_id` - The resource group ID. 
	* `security_group_id` - The security group ID.
	* `status` - The Status of the virtual node.
	* `tags` - A mapping of tags to assign to the resource.
	* `virtual_node_id` - Of the virtual node number.
	* `virtual_node_name` - The name of the virtual node.
	* `vswitch_id` - The vswitch id.
	* `zone_id` - The Zone.