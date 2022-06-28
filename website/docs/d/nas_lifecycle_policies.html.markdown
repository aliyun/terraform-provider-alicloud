---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_lifecycle_policies"
sidebar_current: "docs-alicloud-datasource-nas-lifecycle-policies"
description: |-
  Provides a list of Nas Lifecycle Policies to the user.
---

# alicloud\_nas\_lifecycle\_policies

This data source provides the Nas Lifecycle Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_nas_lifecycle_policies" "ids" {
  file_system_id = "example_value"
  ids            = ["my-LifecyclePolicy-1", "my-LifecyclePolicy-2"]
}
output "nas_lifecycle_policy_id_1" {
  value = data.alicloud_nas_lifecycle_policies.ids.policies.0.id
}

data "alicloud_nas_lifecycle_policies" "nameRegex" {
  file_system_id = "example_value"
  name_regex     = "^my-LifecyclePolicy"
}
output "nas_lifecycle_policy_id_2" {
  value = data.alicloud_nas_lifecycle_policies.nameRegex.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `ids` - (Optional, ForceNew, Computed)  A list of Lifecycle Policy IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Lifecycle Policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Lifecycle Policy names.
* `policies` - A list of Nas Lifecycle Policies. Each element contains the following attributes:
	* `create_time` - The time when the lifecycle management policy was created.
	* `file_system_id` - The ID of the file system.
	* `id` - The ID of the Lifecycle Policy. Its value is same as Queue Name.
	* `lifecycle_policy_name` - The name of the lifecycle management policy.
	* `lifecycle_rule_name` - The rules in the lifecycle management policy.
	* `paths` - The list of absolute paths for multiple directories. In this case, you can associate a lifecycle management policy with each directory.
	* `storage_type` - The storage type of the data that is dumped to the IA storage medium.