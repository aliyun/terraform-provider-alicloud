---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_metric_rule_black_lists"
sidebar_current: "docs-alicloud-datasource-cms-metric-rule-black-lists"
description: |-
  Provides a list of Cloud Monitor Service Metric Rule Black List owned by an Alibaba Cloud account.
---

# alicloud_cms_metric_rule_black_lists

This data source provides Cloud Monitor Service Metric Rule Black List available to the user.[What is Metric Rule Black List](https://www.alibabacloud.com/help/en/cloudmonitor/latest/describemetricruleblacklist)

-> **NOTE:** Available in 1.194.0+

## Example Usage

```
data "alicloud_cms_metric_rule_black_lists" "default" {
  ids        = ["${alicloud_cms_metric_rule_black_lists.default.id}"]
  category   = "ecs"
  namespace  = "acs_ecs_dashboard"
}

output "alicloud_cms_rule_black_list_example_id" {
  value = data.alicloud_cms_metric_rule_black_lists.lists.0.id
}
```

## Argument Reference

The following arguments are supported:
* `category` - (ForceNew,Optional) Cloud service classification. For example, Redis includes kvstore_standard, kvstore_sharding, and kvstore_splitrw.
* `metric_rule_black_list_id` - (ForceNew,Optional) The first ID of the resource
* `namespace` - (ForceNew,Optional) The data namespace of the cloud service.
* `ids` - (Optional, ForceNew, Computed) A list of Metric Rule Black List IDs.
* `metric_rule_black_list_names` - (Optional, ForceNew) The name of the Metric Rule Black List. You can specify at most 10 names.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Metric Rule Black List IDs.
* `names` - A list of name of Metric Rule Black Lists.
* `lists` - A list of Metric Rule Black List Entries. Each element contains the following attributes:
  * `category` - Cloud service classification. For example, Redis includes kvstore_standard, kvstore_sharding, and kvstore_splitrw.
  * `create_time` - The timestamp for creating an alert blacklist policy.Unit: milliseconds.
  * `effective_time` - The effective time range of the alert blacklist policy.
  * `enable_end_time` - The start timestamp of the alert blacklist policy.Unit: milliseconds.
  * `enable_start_time` - The end timestamp of the alert blacklist policy.Unit: milliseconds.
  * `instances` - The list of instances of cloud services specified in the alert blacklist policy.
  * `is_enable` - The status of the alert blacklist policy. Value:-true: enabled.-false: disabled.
  * `metric_rule_black_list_id` - The first ID of the resource
  * `metric_rule_black_list_name` - The name of the alert blacklist policy.
  * `metrics` - Monitoring metrics in the instance.
    * `metric_name` - The name of the monitoring indicator.
    * `resource` - The extended dimension information of the instance. For example, '{"device":"C:"}' indicates that the blacklist policy is applied to all C disks under the ECS instance.
  * `namespace` - The data namespace of the cloud service.
  * `scope_type` - The effective range of the alert blacklist policy. Value:-USER: The alert blacklist policy only takes effect in the current Alibaba cloud account.-GROUP: The alert blacklist policy takes effect in the specified application GROUP.
  * `scope_value` - Application Group ID list. The format is JSON Array.> This parameter is displayed only when 'ScopeType' is 'GROUP.
  * `update_time` - Modify the timestamp of the alert blacklist policy.Unit: milliseconds.
