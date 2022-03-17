---
subcategory: "ApsaraDB for MyBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_dedicated_host_groups"
sidebar_current: "docs-alicloud-datasource-cddc-dedicated-host-groups"
description: |-
  Provides a list of Cddc Dedicated Host Groups to the user.
---

# alicloud\_cddc\_dedicated\_host\_groups

This data source provides the Cddc Dedicated Host Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cddc_dedicated_host_groups" "default" {
  engine = "MongoDB"
}
output "cddc_dedicated_host_group_id" {
  value = data.alicloud_cddc_dedicated_host_groups.default.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew, Available in 1.147.0+) A regex string to filter results by Dedicated Host Group name.
* `engine` - (Optional, ForceNew) Database Engine Type. Valid values:`Redis`, `SQLServer`, `MySQL`, `PostgreSQL`, `MongoDB`
* `ids` - (Optional, ForceNew, Computed)  A list of Dedicated Host Group IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Dedicated Host Group names.
* `groups` - A list of Cddc Dedicated Host Groups. Each element contains the following attributes:
	* `allocation_policy` -  The policy that is used to allocate resources in the dedicated cluster. Valid values:`Evenly`,`Intensively`
	* `bastion_instance_id` - The Bastion Instance id of the Dedicated Host Group.
	* `cpu_allocate_ration` - The CPU overcommitment ratio of the dedicated cluster. If you set this parameter to 200, the CPU resources that can be allocated are twice as many as the CPU resources that are provided. This maximizes the CPU utilization. Valid values: 100 to 300. Default value: 200.
	* `cpu_allocated_amount` - The CPU Allocated Amount of the Dedicated Host Group.
	* `cpu_allocation_ratio` - The CPU overcommitment ratio of the dedicated cluster.Valid values: 100 to 300. Default value: 200.
	* `create_time` - The Created Time of the Dedicated Host Group.
	* `dedicated_host_count_group_by_host_type` - The Dedicated Host Count Group by Host Type of the Dedicated Host Group.
	* `dedicated_host_group_desc` -The name of the dedicated cluster. The name must be 1 to 64 characters in length and can contain letters, digits, underscores (_), and hyphens (-). It must start with a letter.
	* `dedicated_host_group_id` - Dedicated Host Group ID.
	* `deploy_type` - The Deployment Type of the Dedicated Host Group.
	* `disk_allocate_ration` - The storage overcommitment ratio of the dedicated cluster.Valid values: 100 to 300. Default value: 200.
	* `disk_allocated_amount` - The Disk Allocated Amount of the Dedicated Host Group.
	* `disk_allocation_ratio` - The Disk Allocation Ratio of the Dedicated Host Group.
	* `disk_used_amount` - The DiskUsedAmount of the Dedicated Host Group.
	* `disk_utility` - The DiskUtility of the Dedicated Host Group.
	* `engine` - Database Engine Type.The database engine of the dedicated cluster. Valid values:`Redis`, `SQLServer`, `MySQL`, `PostgreSQL`, `MongoDB`
	* `host_number` - The Total Host Number  of the Dedicated Host Group.
	* `host_replace_policy` - The policy based on which the system handles host failures. Valid values:`Auto`,`Manual`
	* `id` - The ID of the Dedicated Host Group.
	* `instance_number` - The Total Instance Number of the Dedicated Host Group.
	* `mem_allocate_ration` - The maximum memory usage of each host in the dedicated cluster.Valid values: 0 to 90. Default value: 90.
	* `mem_allocated_amount` - The MemAllocatedAmount of the Dedicated Host Group.
	* `mem_allocation_ratio` - The Memory Allocation Ratio of the Dedicated Host Group.
	* `mem_used_amount` - The MemUsedAmount of the Dedicated Host Group.
	* `mem_utility` - The Mem Utility of the Dedicated Host Group.
	* `text` - The Text of the Dedicated Host Group.
	* `vpc_id` - The virtual private cloud (VPC) ID of the dedicated cluster.
	* `zone_id_list` - The ZoneIDList of the Dedicated Host Group.
