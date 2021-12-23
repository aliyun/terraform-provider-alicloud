---
subcategory: "Cassandra"
layout: "alicloud"
page_title: "Alicloud: alicloud_cassandra_backup_plans"
sidebar_current: "docs-alicloud-datasource-cassandra-backup-plans"
description: |-
  Provides a list of Cassandra Backup Plans to the user.
---

# alicloud\_cassandra\_backup\_plans

This data source provides the Cassandra Backup Plans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.128.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cassandra_backup_plans" "example" {
  cluster_id = "example_value"
}

output "first_cassandra_backup_plan_id" {
  value = data.alicloud_cassandra_backup_plans.example.plans.0.id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the cluster for the backup.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Backup Plan IDs.
* `plans` - A list of Cassandra Backup Plans. Each element contains the following attributes:
	* `active` - Specifies whether to activate the backup plan.
	* `backup_period` - The backup cycle. Valid values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, and Sunday.
	* `backup_time` - The start time of the backup task each day. The time is displayed in UTC and denoted by Z.
	* `cluster_id` - The ID of the cluster for the backup.
	* `create_time` - The time when the backup plan was created.
	* `data_center_id` - The ID of the data center for the backup in the cluster.
	* `id` - The ID of the Backup Plan.
	* `retention_period` - The duration for which you want to retain the backup. Valid values: 1 to 30. Unit: days.
