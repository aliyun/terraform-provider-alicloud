---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_policies"
sidebar_current: "docs-alicloud-datasource-apig-policies"
description: |-
  Provides a list of Apig Policy owned by an Alibaba Cloud account.
---

# alicloud_apig_policies

This data source provides Apig Policy available to the user.[What is Policy](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateAndAttachPolicy)

-> **NOTE:** Available since v1.292.0.

## Example Usage

```terraform
variable "gateway_id" {
  description = "The ID of an existing APIG gateway"
  type        = string
}

data "alicloud_apig_policies" "default" {
  gateway_id = var.gateway_id
}

output "first_policy_id" {
  value = data.alicloud_apig_policies.default.policies[0].policy_id
}
```

## Argument Reference

The following arguments are supported:
* `attach_resource_ids` - (Optional) The ID of the attach point resource to filter policies by.
* `attach_resource_type` - (Optional) Policies support mount point types. Valid values: `HttpApi`, `Operation`, `GatewayRoute`, `GatewayService`, `GatewayServicePort`, `Domain`, `Gateway`.
* `environment_id` - (Optional) Environment id.
* `gateway_id` - (Optional) Gateway id.
* `ids` - (Optional, Computed) A list of Policy IDs.
* `name_regex` - (Optional) A regex string to filter results by policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional) Default to `false`. Set it to `true` to output more details about resource attributes that belong to this policy.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `names` - A list of Policy names.
* `policies` - A list of Policy Entries. Each element contains the following attributes:
  * `id` - The ID of the Policy.
  * `policy_id` - The first ID of the resource.
  * `policy_name` - Policy name.
  * `policy_class_id` - Policy class id.
  * `policy_class_name` - Policy class name.
  * `policy_config` - Policy configuration.
  * `attach_resource_ids` - The Mount point id list.
  * `attach_resource_type` - The mount point type of the policy.
  * `environment_id` - Environment id.
  * `gateway_id` - Gateway id.
  * `policy_attachment_id` - Policy attachment id.
