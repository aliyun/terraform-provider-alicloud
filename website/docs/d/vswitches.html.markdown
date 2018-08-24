---
layout: "alicloud"
page_title: "Alicloud: alicloud_vswitches"
sidebar_current: "docs-alicloud-datasource-vswitches"
description: |-
    Provides a list of VSwitch owned by an Alibaba Cloud account.
---

# alicloud\_vswitches

This data source provides a list of VSwitches owned by an Alibaba Cloud account.

## Example Usage

```
data "alicloud_vswitches" "vswitches_ds" {
  cidr_block = "172.16.0.0/12"
  name_regex = "^foo"
}

resource "alicloud_instance" "foo" {
  # ...
  instance_name = "in-the-vpc"
  vswitch_id = "${data.alicloud_vswitches.vswitches_ds.vswitches.0.id}"
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional) Filter results by a specific CIDR block. For example: "172.16.0.0/12".
* `zone_id` - (Optional) The availability zone of the VSwitch.
* `name_regex` - (Optional) A regex string to filter results by name.
* `is_default` - (Optional, type: bool) Indicate whether the VSwitch is created by the system.
* `vpc_id` - (Optional) ID of the VPC that owns the VSwitch.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `vswitches` - A list of VSwitches. Each element contains the following attributes:
  * `id` - ID of the VSwitch.
  * `zone_id` - ID of the availability zone where the VSwitch is located.
  * `vpc_id` - ID of the VPC that owns the VSwitch.
  * `name` - Name of the VSwitch.
  * `instance_ids` - List of ECS instance IDs in the specified VSwitch.
  * `cidr_block` - CIDR block of the VSwitch.
  * `description` - Description of the VSwitch.
  * `is_default` - Whether the VSwitch is the default one in the region.
  * `creation_time` - Time of creation.
