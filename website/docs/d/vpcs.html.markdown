---
layout: "alicloud"
page_title: "Alicloud: alicloud_vpcs"
sidebar_current: "docs-alicloud-datasource-vpcs"
description: |-
    Provides a list of VPCs which owned by an Alicloud account.
---

# alicloud\_vpcs

The VPCs data source lists a number of VPCs resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_vpcs" "multi_vpc"{
	cidr_block="172.16.0.0/12"
	status="Available"
	name_regex="^foo"
}

```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional) Limit search to specific cidr block,like "172.16.0.0/12".
* `status` - (Optional) Limit search to specific status - valid value is "Pending" or "Available".
* `name_regex` - (Optional) A regex string of VPC name.
* `is_default` - (Optional) Whether the VPC is the default VPC in the specified region - valid value is true or false.
* `vswitch_id` - (Optional) Retrieving VPC according to the specified VSwitch.
* `output_file` - (Optional) The name of file that can save vpcs data source after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the VPC.
* `region_id` - ID of the region where VPC belongs.
* `status` - Status of the VPC.
* `vpc_name` - Name of the VPC.
* `vswitch_ids` - List of VSwitch IDs in the specified VPC
* `cidr_block` - CIDR block of the VPC.
* `vrouter_id` - ID of the VRouter
* `route_table_id` - Route table ID of the VRouter
* `description` - Description of the VPC
* `is_default` - Whether the VPC is the default VPC in the belonging region.
* `creation_time` - Time of creation.