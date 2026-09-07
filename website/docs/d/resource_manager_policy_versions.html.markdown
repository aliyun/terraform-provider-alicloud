---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_policy_versions"
sidebar_current: "docs-alicloud-datasource-resource-manager-policy-versions"
description: |-
    Provides a list of Resource Manager Policy Versions to the user.
---

# alicloud\_resource\_manager\_policy\_versions

This data source provides the Resource Manager Policy Versions of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.85.0+.

## Example Usage

```terraform
data "alicloud_resource_manager_policy_versions" "default" {
  policy_name = "tftest"
  policy_type = "Custom"
}

output "first_policy_version_id" {
  value = "${data.alicloud_resource_manager_policy_versions.default.versions.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of policy version IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `policy_name` - (Required) The name of the policy.
* `policy_type` - (Required) The type of the policy. Valid values:`Custom` and `System`.
* `enable_details` -(Optional, Available in v1.114.0+) Default to `false`. Set it to true can output more details.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of policy version IDs.
* `versions` - A list of policy versions. Each element contains the following attributes:
    * `id` - The ID of the resource, the value is `<policy_name>`:`<version_id>`.
    * `version_id`- The ID of the policy version.
    * `create_date`- (Removed form v1.114.0)The time when the policy version was created.
    * `is_default_version`- Indicates whether the policy version is the default version.
    * `policy_document`- (Available in v1.114.0+) The policy document of the policy version.
    
    
