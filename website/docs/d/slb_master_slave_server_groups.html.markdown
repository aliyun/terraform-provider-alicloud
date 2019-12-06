---
subcategory: "Server Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_master_slave_server_groups"
sidebar_current: "docs-alicloud-datasource-slb-master-slave-server-groups"
description: |-
    Provides a list of master slave server groups related to a server load balancer to the user.
---

# alicloud\_slb\_master\_slave\_server\_groups

This data source provides the master slave server groups related to a server load balancer.

-> **NOTE:** Available in 1.54.0+

## Example Usage

```
data "alicloud_slb_master_slave_server_groups" "sample_ds" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
}

output "first_slb_server_group_id" {
  value = "${data.alicloud_slb_master_slave_server_groups.sample_ds.groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - ID of the SLB.
* `ids` - (Optional) A list of master slave server group IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by master slave server group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB master slave server groups IDs.
* `names` - A list of SLB master slave server groups names.
* `groups` - A list of SLB master slave server groups. Each element contains the following attributes:
  * `id` - master slave server group ID.
  * `name` - master slave server group name.
  * `servers` - ECS instances associated to the group. Each element contains the following attributes:
    * `instance_id` - ID of the attached ECS instance.
    * `weight` - Weight associated to the ECS instance.
    * `port` - The port used by the master slave server group.
    * `server_type` - The server type of the attached ECS instance.
    * `is_backup` - (Removed from v1.63.0) Determine if the server is executing.

