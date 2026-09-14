---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_escalation_policies"
description: |-
  Provides a list of Cms Escalation Policies to the user.
---

# alicloud_cms_escalation_policies

This data source provides the list of Cms Escalation Policies of current Alibaba Cloud user.

For information about Cms Escalation Policy and how to use it, see [What is Escalation Policy](https://next.api.alibabacloud.com/document/Cms/2024-03-30/ListEscalationPolicy).

-> **NOTE:** Available since v1.297.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_escalation_policies" "default" {
  workspace  = "example-workspace"
  name_regex = "tf-acc-escalation"
  ids        = ["example-uuid:example-workspace"]
}

output "first_escalation_policy_id" {
  value = data.alicloud_cms_escalation_policies.default.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Optional) The name of the workspace to which the escalation policies belong.
* `ids` - (Optional) A list of escalation policy IDs. The value is formatted as `<uuid>:<workspace>`.
* `name_regex` - (Optional) A regex string to filter results by the escalation policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `ids` - A list of escalation policy IDs. The value is formatted as `<uuid>:<workspace>`.
* `policies` - A list of Cms Escalation Policies. Each element contains the following attributes:
  * `id` - The ID of the escalation policy. It is formatted as `<uuid>:<workspace>`.
  * `uuid` - The UUID of the escalation policy.
  * `name` - The name of the escalation policy.
  * `workspace` - The name of the workspace to which the escalation policy belongs.
  * `description` - The description of the escalation policy.
  * `enable` - Whether the escalation policy is enabled.
  * `create_time` - The creation time of the resource.
  * `update_time` - The last modified time of the resource.
