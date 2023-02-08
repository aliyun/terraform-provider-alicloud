---
subcategory: "EBS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_enterprise_snapshot_policies"
sidebar_current: "docs-alicloud-datasource-ebs-enterprise-snapshot-policies"
description: |-
  Provides a list of Ebs Enterprise Snapshot Policy owned by an Alibaba Cloud account.
---

# alicloud_ebs_enterprise_snapshot_policies

This data source provides Ebs Enterprise Snapshot Policy available to the user.[What is Enterprise Snapshot Policy](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available in 1.200.0+

## Example Usage

```terraform
data "alicloud_ebs_enterprise_snapshot_policies" "default" {
  ids        = ["${alicloud_ebs_enterprise_snapshot_policy.default.id}"]
  name_regex = alicloud_ebs_enterprise_snapshot_policy.default.name
}

output "alicloud_ebs_enterprise_snapshot_policy_example_id" {
  value = data.alicloud_ebs_enterprise_snapshot_policies.default.enterprise_snapshot_policies.0.id
}
```

## Argument Reference

The following arguments are supported:
* `enterprise_snapshot_policy_ids` - (ForceNew,Optional) A list of Enterprise Snapshot Policy IDs.
* `resource_group_id` - (ForceNew,Optional) The ID of the resource group
* `ids` - (Optional, ForceNew, Computed) A list of Enterprise Snapshot Policy IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Enterprise Snapshot Policy IDs.
* `names` - A list of name of Enterprise Snapshot Policys.
* `enterprise_snapshot_policies` - A list of Enterprise Snapshot Policy Entries. Each element contains the following attributes:
  * `id` - The ID of the resource.
  * `create_time` - The creation time of the resource.
  * `cross_region_copy_info` - Snapshot replication information.
    * `enabled` - Enable Snapshot replication.
    * `regions` - Destination region for Snapshot replication.
      * `region_id` - Destination region ID.
      * `retain_days` - Number of days of snapshot retention for replication.
  * `desc` - Description information representing the resource
  * `enterprise_snapshot_policy_id` - The first ID of the resource
  * `enterprise_snapshot_policy_name` - The name of the resource
  * `resource_group_id` - The ID of the resource group
  * `retain_rule` - Snapshot retention policy representing resources
    * `number` - Retention based on counting method.
    * `time_interval` - Time unit.
    * `time_unit` - Time-based retention.
  * `schedule` - The scheduling plan that represents the resource.
    * `cron_expression` - CronTab expression.
  * `status` - The status of the resource
  * `storage_rule` - Snapshot storage policy
    * `enable_immediate_access` - Snapshot speed available.
  * `tags` - The tag of the resource
  * `target_type` - Represents the target type of resource binding
