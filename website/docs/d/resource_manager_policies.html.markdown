---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_policies"
sidebar_current: "docs-alicloud-datasource-resource-manager-policies"
description: |-
    Provides a list of Resource Manager Policies to the user.
---

# alicloud\_resource\_manager\_policies

This data source provides the Resource Manager Policies of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.86.0+.

## Example Usage

```terraform
data "alicloud_resource_manager_policies" "example" {
  name_regex        = "tftest"
  description_regex = "tftest_policy"
  policy_type       = "Custom"
}

output "first_policy_id" {
  value = "${data.alicloud_resource_manager_policies.example.policies.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Resource Manager Policy IDs.
* `name_regex` - (Optional) A regex string to filter results by policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `policy_type` - (Optional) The type of the policy. If you do not specify this parameter, the system lists all types of policies. Valid values: `Custom` and `System`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of policy IDs.
* `names` - A list of policy names.
* `policies` - A list of policies. Each element contains the following attributes:
    * `id` - The ID of the policy.
    * `policy_name`- The name of the policy.
    * `policy_type`- The type of the policy.
    * `description` - The description of the policy.
    * `create_date` - The time when the policy was created.
    * `update_date` - The time when the policy was updated.
    * `default_version` - The default version of the policy.
    * `attachment_count` - The number of times the policy is referenced.
