---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_hana_backup_plans"
sidebar_current: "docs-alicloud-datasource-hbr-hana-backup-plans"
description: |-
  Provides a list of Hbr Hana Backup Plans to the user.
---

# alicloud\_hbr\_hana\_backup\_plans

This data source provides the Hbr Hana Backup Plans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.179.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_hana_backup_plans" "ids" {
  cluster_id = "example_value"
  ids        = ["example_value-1", "example_value-2"]
}
output "hbr_hana_backup_plan_id_1" {
  value = data.alicloud_hbr_hana_backup_plans.ids.plans.0.id
}
```

## Argument Reference

The following arguments are supported:

* `database_name` - (Optional, ForceNew) The name of the database.
* `ids` - (Optional, ForceNew, Computed) A list of Hana Backup Plan IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vault_id` - (Optional, ForceNew) The id of the vault.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Hana Backup Plan name.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `plans` - A list of Hbr Hana Backup Plans. Each element contains the following attributes:
	* `backup_prefix` - The backup prefix.
	* `backup_type` - The backup type.
	* `cluster_id` - The ID of the SAP HANA instance.
	* `database_name` - The name of the database.
	* `status` - The status of the resource.
	* `id` - The ID of the resource.
	* `plan_id` - The ID of the backup plan.
	* `plan_name` - The name of the backup plan.
	* `schedule` - The backup policy.
	* `vault_id` - The ID of the backup vault.