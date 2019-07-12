---
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_groups"
sidebar_current: "docs-alicloud_ess_scaling_groups"
description: |-
    Provides a list of scaling groups available to the user.
---

# alicloud_ess_scaling_groups

This data source provides available scaling group resources. 

## Example Usage

```
data "alicloud_ess_scaling_groups" "scalinggroups_ds" {
  ids        = ["scaling_group_id1", "scaling_group_id2"]
  name_regex = "scaling_group_name"
}

output "first_scaling_group" {
  value = "${data.alicloud_ess_scaling_groups.scalinggroups_ds.groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter resulting scaling groups by name.
* `ids` - (Optional) A list of scaling group IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of scaling group ids.
* `names` - A list of scaling group names.
* `scaling_groups` - A list of scaling groups. Each element contains the following attributes:
  * `id` - ID of the scaling group.
  * `name` - Name of the scaling group.
  * `active_scaling_configuration` -Active scaling configuration for scaling group.
  * `launch_template_id` - Active launch template ID for scaling group.
  * `launch_template_version` - Version of active launch template.
  * `region_id` - Region ID the scaling group belongs to.
  * `min_size` - The minimum number of ECS instances.
  * `max_size` - The maximum number of ECS instances.
  * `cooldown_time` - Default cooldown time of scaling group.
  * `removal_policies` - Removal policy used to select the ECS instance to remove from the scaling group.
  * `load_balancer_ids` - Slb instances id which the ECS instance attached to.
  * `db_instance_ids` - Db instances id which the ECS instance attached to.
  * `vswitch_ids` - Vswitches id in which the ECS instance launched.
  * `lifecycle_state` - Lifecycle state of scaling group.
  * `total_capacity` - Number of instances in scaling group.
  * `active_capacity` - Number of active instances in scaling group.
  * `pending_capacity` - Number of pending instances in scaling group.
  * `removing_capacity` - Number of removing instances in scaling group.
  * `creation_time` - Creation time of scaling group.
  
  