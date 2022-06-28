---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_backup_policies"
sidebar_current: "docs-alicloud-datasource-click-house-backup-policies"
description: |-
  Provides a list of Click House Backup Policies to the user.
---

# alicloud\_click\_house\_backup\_policies

This data source provides the Click House Backup Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_click_house_backup_policies" "example" {
  db_cluster_id = "example_value"
}
output "click_house_backup_policy_id_1" {
  value = data.alicloud_click_house_backup_policies.example.policies.0.id
}

```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Request, ForceNew) The db cluster id.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `policies` - A list of Click House Backup Policies. Each element contains the following attributes:
    * `backup_retention_period` - Data backup days. Valid values: `7` to `730`.
    * `db_cluster_id` - The db cluster id.
    * `id` - The ID of the Backup Policy.
    * `preferred_backup_period` - DBCluster Backup period.
    * `preferred_backup_time` - Backup Time, UTC time.
    * `status` - The status of the resource.