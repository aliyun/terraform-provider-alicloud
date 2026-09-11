---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_custom_check_standard_policies"
sidebar_current: "docs-alicloud-datasource-threat-detection-custom-check-standard-policies"
description: |-
  This data source provides the Threat Detection Custom Check Standard Policies of the current Alibaba Cloud user.
---

# alicloud_threat_detection_custom_check_standard_policies

This data source provides the Threat Detection Custom Check Standard Policies of the current Alibaba Cloud user.

For information about Threat Detection Custom Check Standard Policies, see [ListCheckPolicies](https://www.alibabacloud.com/help/en/security-center/latest/api-sas-2018-12-03-listcheckpolicies).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_custom_check_standard_policies" "default" {
  policy_type = "STANDARD"
  type        = "AISPM"
}

output "first_policy_id" {
  value = data.alicloud_threat_detection_custom_check_standard_policies.default.policies.0.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `dependent_policy_id` - (Optional) The ID of the associated parent policy. The dependency order from low to high is: Section -> Requirement -> Standard.
* `ids` - (Optional) A list of Custom Check Standard Policy IDs. Each ID is formatted as `<policy_id>:<policy_type>`.
* `name_regex` - (Optional) A regex string to filter the results by the `policy_show_name` of the policy.
* `output_file` - (Optional) File name where to save data source results after running `terraform plan`.
* `policy_type` - (Required) The policy type of the custom check item rule. Valid values: `STANDARD`, `REQUIREMENT`, `SECTION`.
* `type` - (Optional) The name of the major policy category. Valid values: `AISPM`, `IDENTITY_PERMISSION`, `RISK`, `COMPLIANCE`.

## Attributes Reference

The following attributes are exported in addition to the `arguments` above:

* `policies` - A list of Custom Check Standard Policies. Each element contains the following attributes:
  * `check_type` - The check type of the policy.
  * `dependent_policy_id` - The ID of the associated parent policy.
  * `id` - The resource ID in Terraform, formatted as `<policy_id>:<policy_type>`.
  * `policy_id` - The ID of the custom policy.
  * `policy_show_name` - The name of the custom policy.
  * `policy_type` - The policy category type.
  * `type` - The name of the major policy category.
* `names` - A list of policy names.
