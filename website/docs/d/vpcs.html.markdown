---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpcs"
sidebar_current: "docs-alicloud-datasource-vpcs"
description: |-
    Provides a list of VPCs owned by an Alibaba Cloud account.
---

# alicloud_vpcs

This data source provides VPCs available to the user.

## Example Usage

```terraform
data "alicloud_vpcs" "vpcs_ds" {
  cidr_block = "172.16.0.0/12"
  status     = "Available"
  name_regex = "^foo"
}

output "first_vpc_id" {
  value = "${data.alicloud_vpcs.vpcs_ds.vpcs.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional) Filter results by a specific CIDR block. For example: "172.16.0.0/12".
* `status` - (Optional) Filter results by a specific status. Valid value are `Pending` and `Available`.
* `name_regex` - (Optional) A regex string to filter VPCs by name.
* `is_default` - (Optional, type: bool) Indicate whether the VPC is the default one in the specified region.
* `vswitch_id` - (Optional) Filter results by the specified vSwitch.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `ids` - (Optional, Available in 1.52.0+) A list of VPC IDs.
* `resource_group_id` - (Optional, ForceNew, Available in 1.60.0+) The Id of resource group which VPC belongs.
* `dhcp_options_set_id` - (Optional, ForceNew, Available in v1.119.0+) The ID of dhcp options set.
* `dry_run` - (Optional, ForceNew, Available in v1.119.0+) Indicates whether to check this request only. Valid values: `true` and `false`.
* `vpc_name` - (Optional, ForceNew, Available in v1.119.0+) The name of the VPC.
* `vpc_owner_id` - (Optional, ForceNew, Available in v1.119.0+) The owner ID of VPC.
* `enable_details` -(Optional, Available in v1.119.0+) Default to `true`. Set it to true can output the `route_table_id`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of VPC IDs.
* `names` - A list of VPC names.
* `vpcs` - A list of VPCs. Each element contains the following attributes:
  * `id` - ID of the VPC.
  * `region_id` - ID of the region where the VPC is located.
  * `status` - Status of the VPC.
  * `vpc_name` - Name of the VPC.
  * `vswitch_ids` - List of vSwitch IDs in the specified VPC
  * `cidr_block` - CIDR block of the VPC.
  * `vrouter_id` - ID of the VRouter.
  * `route_table_id` - Route table ID of the VRouter.
  * `description` - Description of the VPC
  * `dhcp_options_set_id` - The ID of the DHCP options set.
  * `dhcp_options_set_status` - The status of the DHCP options set. Valid values: `Available`, `InUse`, `Deleted`, `Pending`.
  * `dns_hostname_status` - Indicates whether the Domain Name System (DNS) feature is enabled.
  * `is_default` - Whether the VPC is the default VPC in the region.
  * `creation_time` - Time of creation.
  * `tags` - A map of tags assigned to the VPC.
  * `ipv6_cidr_block` - The IPv6 CIDR block of the VPC.
  * `router_id` - The ID of the VRouter.
  * `secondary_cidr_blocks` - A list of secondary IPv4 CIDR blocks of the VPC.
  * `user_cidrs` - A list of user CIDRs.
  * `vpc_id` - ID of the VPC.
  
