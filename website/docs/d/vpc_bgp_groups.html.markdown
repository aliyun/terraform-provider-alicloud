---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_bgp_groups"
sidebar_current: "docs-alicloud-datasource-vpc-bgp-groups"
description: |-
  Provides a list of Vpc Bgp Groups to the user.
---

# alicloud\_vpc\_bgp\_groups

This data source provides the Vpc Bgp Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_bgp_groups" "ids" {
  ids = ["example_value"]
}
output "vpc_bgp_group_id_1" {
  value = data.alicloud_vpc_bgp_groups.ids.groups.0.id
}

data "alicloud_vpc_bgp_groups" "nameRegex" {
  name_regex = "^my-BgpGroup"
}
output "vpc_bgp_group_id_2" {
  value = data.alicloud_vpc_bgp_groups.nameRegex.groups.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Bgp Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Bgp Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `router_id` - (Optional, ForceNew) The ID of the virtual border router (VBR) that is associated with the BGP group.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Available`, `Deleting` and `Pending`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Bgp Group names.
* `groups` - A list of Vpc Bgp Groups. Each element contains the following attributes:
	* `auth_key` - The key used by the BGP group.
	* `bgp_group_name` - The name of the BGP group.
	* `description` - Description of the BGP group.
	* `hold` - The hold time to wait for the incoming BGP message. If no message has been passed in after the hold time, the BGP neighbor is considered disconnected.
	* `id` - The ID of the Bgp Group.
	* `ip_version` - IP version.
	* `is_fake_asn` - Whether the AS number is false.
	* `keepalive` - The keepalive time.
	* `local_asn` - The local AS number.
	* `peer_asn` - The autonomous system (AS) number of the BGP peer.
	* `route_limit` - Routing limits.
	* `router_id` - The ID of the VBR.
	* `status` - The status of the resource.