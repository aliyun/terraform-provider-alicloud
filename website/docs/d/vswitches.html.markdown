---
layout: "alicloud"
page_title: "Alicloud: alicloud_vswitches"
sidebar_current: "docs-alicloud-datasource-vswitches"
description: |-
    Provides a list of VSwitch which owned by an Alicloud account.
---

# alicloud\_vswitches

The Virtual sunbet data source lists a list of vswitches resource information owned by an Alicloud account,
and each vswitch including its basic attribution, VPC ID and containing ECS instance IDs.

## Example Usage

```
data "alicloud_vswitches" "subnets"{
    cidr_block="172.16.0.0/12"
    name_regex="^foo"
}

resource "alicloud_instance" "foo" {
    ...
    instance_name =  "in-the-vpc"
    vswitch_id = "${data.alicloud_vswitches.subnets.vswitches.0.id}"
    ...
}

```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional) Limit search to specific cidr block,like "172.16.0.0/12".
* `zone_id` - (Optional) The availability zone for one vswitch.
* `name_regex` - (Optional) A regex string of VSwitch name.
* `is_default` - (Optional) Whether the Vswitch is created by system - valid value is true or false.
* `vpc_id` - (Optional) VPC ID in which vswitch belongs.
* `output_file` - (Optional) The name of file that can save vswitches data source after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `vswitches` A list of vswitches. It contains several attributes to `Block VSwitches`.

### Block VSwitches

Attributes for vswitches:

* `id` - ID of the VSwitch.
* `zone_id` - ID of the availability zone where VSwitch belongs.
* `vpc_id` - ID of the VPC where VSwitch belongs.
* `name` - Name of the VSwitch.
* `instance_ids` - List of ECS instance IDs in the specified VSwitch.
* `cidr_block` - CIDR block of the VSwitch.
* `description` - Description of the VSwitch
* `is_default` - Whether the VSwitch is the default VSwitch in the belonging region.
* `creation_time` - Time of creation.
