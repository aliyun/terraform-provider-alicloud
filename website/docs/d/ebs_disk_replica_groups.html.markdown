---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_disk_replica_groups"
sidebar_current: "docs-alicloud-datasource-ebs-disk-replica-groups"
description: |-
  Provides a list of Ebs Disk Replica Groups to the user.
---

# alicloud\_ebs\_disk\_replica\_groups

This data source provides the Ebs Disk Replica Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.187.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ebs_disk_replica_groups" "ids" {
  ids = ["example_id"]
}

output "ebs_disk_replica_group_id_1" {
  value = data.alicloud_ebs_disk_replica_groups.ids.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Disk Replica Group IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `groups` - A list of Ebs Disk Replica Groups. Each element contains the following attributes:
	* `description` - The description of the consistent replication group.
	* `destination_region_id` - The ID of the region to which the disaster recovery site belongs.
	* `destination_zone_id` - The ID of the zone to which the disaster recovery site belongs.
	* `group_name` - Consistent replication group name.
	* `id` - The ID of the Disk Replica Group.
	* `primary_region` - The initial source region of the replication group.
	* `primary_zone` - The initial source available area of the replication group.
	* `replica_group_id` - The ID of the consistent replication group.
	* `rpo` - The recovery point objective (RPO) of the replication pair-consistent group.
	* `site` - Site information sources for replication pairs and consistent replication groups. 
	* `source_region_id` - The ID of the region to which the production site belongs.
	* `source_zone_id` - The ID of the zone to which the production site belongs.
	* `standby_region` - The initial destination region of the replication group.
	* `standby_zone` - The initial destination zone of the replication group.
	* `status` - The status of the consistent replication group. Possible values:
