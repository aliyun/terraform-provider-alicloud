---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vswitches"
sidebar_current: "docs-alicloud-datasource-vswitches"
description: |-
    Provides a list of vSwitch owned by an Alibaba Cloud account.
---

# alicloud_vswitches

This data source provides a list of VSwitches owned by an Alibaba Cloud account.

## Example Usage

```terraform
variable "name" {
  default = "vswitchDatasourceName"
}
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
  vpc_name   = "${var.name}"
}

resource "alicloud_vswitch" "vswitch" {
  vswitch_name      = "${var.name}"
  cidr_block        = "172.16.0.0/24"
  vpc_id            = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

data "alicloud_vswitches" "default" {
  name_regex = "${alicloud_vswitch.vswitch.vswitch_name}"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional) Filter results by a specific CIDR block. For example: "172.16.0.0/12".
* `zone_id` - (Optional) The availability zone of the vSwitch.
* `name_regex` - (Optional) A regex string to filter results by name.
* `is_default` - (Optional, type: bool) Indicate whether the vSwitch is created by the system.
* `vpc_id` - (Optional) ID of the VPC that owns the vSwitch.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `ids` - (Optional, Available in 1.52.0+) A list of vSwitch IDs.
* `resource_group_id` - (Optional, ForceNew, Available in 1.60.0+) The Id of resource group which VSWitch belongs.
* `dry_run` - (Optional, ForceNew, Available in 1.119.0+) Specifies whether to precheck this request only. Valid values: `true` and `false`.
* `route_table_id` - (Optional, ForceNew, Available in 1.119.0+) The route table ID of the vSwitch.
* `status` - (Optional, ForceNew, Available in 1.119.0+) The status of the vSwitch. Valid values: `Available` and `Pending`.
* `vswitch_name` - (Optional, ForceNew, Available in 1.119.0+) The name of the vSwitch.
* `vswitch_owner_id` - (Optional, ForceNew, Available in 1.119.0+) The vSwitch owner id.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of vSwitch IDs.
* `names` - A list of vSwitch names.
* `vswitches` - A list of VSwitches. Each element contains the following attributes:
  * `id` - ID of the vSwitch.
  * `zone_id` - ID of the availability zone where the vSwitch is located.
  * `vpc_id` - ID of the VPC that owns the vSwitch.
  * `name` - Name of the vSwitch.
  * `cidr_block` - CIDR block of the vSwitch.
  * `instance_ids` - (Deprecated in v1.119.0+) List of ECS instance IDs in the specified vSwitch.
  * `description` - Description of the vSwitch.
  * `is_default` - Whether the vSwitch is the default one in the region.
  * `creation_time` - Time of creation.
  * `available_ip_address_count` - The available ip address count of the vSwitch.
  * `resource_group_id` - The resource group ID of the vSwitch.
  * `route_table_id` - The route table ID of the vSwitch.
  * `status` - The status of the vSwitch.
  * `tags` - The Tags of the vSwitch.
  * `vswitch_id` - ID of the vSwitch.
  * `vswitch_name` - Name of the vSwitch.
  * `ipv6_cidr_block` - The IPv6 CIDR block of the switch.
