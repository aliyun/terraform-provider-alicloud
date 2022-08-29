---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_peerings"
sidebar_current: "docs-alicloud-datasource-vpc-peerings"
description: |-
  Provides a list of Vpc Peerings to the user.
---

# alicloud\_vpc\_peerings

This data source provides the Vpc Peerings of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.184.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_peerings" "ids" {}
output "vpc_peering_id_1" {
  value = data.alicloud_vpc_peerings.ids.peerings.0.id
}

data "alicloud_vpc_peerings" "nameRegex" {
  name_regex = "^my-Peering"
}
output "vpc_peering_id_2" {
  value = data.alicloud_vpc_peerings.nameRegex.peerings.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Peering IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Peering name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `peering_name` - (Optional, ForceNew) The name of the resource.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Accepting`, `Activated`, `Creating`, `Deleted`, `Deleting`, `Expired`, `Rejected`, `Updating`.
* `vpc_id` - (Optional, ForceNew) The ID of the requester VPC.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Peering names.
* `peerings` - A list of Vpc Peerings. Each element contains the following attributes:
	* `accepting_ali_uid` - The ID of the Alibaba Cloud account (primary account) of the receiving end of the VPC peering connection to be created.
	* `accepting_region_id` - The region ID of the recipient of the VPC peering connection to be created.
	* `accepting_vpc_id` - The VPC ID of the receiving end of the VPC peer connection.
	* `bandwidth` - The bandwidth of the VPC peering connection to be modified. Unit: Mbps.
	* `create_time` - The creation time of the resource.
	* `description` - The description of the VPC peer connection to be created.
	* `id` - The ID of the Peering.
	* `peering_id` - The first ID of the resource.
	* `peering_name` - The name of the resource.
	* `status` - The status of the resource.
	* `vpc_id` - The ID of the requester VPC.